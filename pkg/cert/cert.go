package cert

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	client "github.com/linuxuser586/apis/grpc/pki/client/v1"
	"github.com/linuxuser586/cassandra/pkg/singleton/filesystem"
	"github.com/linuxuser586/common/pkg/net"
	"github.com/linuxuser586/common/pkg/os"
	"github.com/linuxuser586/common/pkg/os/exec"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

const (
	// KeystoreFile holds the cassandra client key
	KeystoreFile = "/tmp/keystore.p12"
	// TruststoreFile hods the CA certificate
	TruststoreFile = "/tmp/truststore.p12"
	keyFile        = "/tmp/client.key"
	crtFile        = "/tmp/client.crt"
	caFile         = "/conf/ca.crt"
	hostEnv        = "PKI_CLIENT_GRPC_HOST"
	portEnv        = "PKI_CLIENT_GRPC_PORT"
	passEnv        = "CASSANDRA_KEYSTORE_PASS"
	keytool        = "/usr/lib/jvm/java-1.8.0-openjdk-amd64/jre/bin/keytool"
)

var (
	// Pass is the keystore/trustore password
	Pass  = "cassandra"
	fs    = filesystem.Singleton()
	podIP = net.PodIP
	ip    string
)

func init() {
	p := os.Getenv(passEnv)
	if p != "" {
		Pass = p
	}
}

// Setup the certificates for Cassandra
func Setup() error {
	var err error
	ip, err = podIP()
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(address(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := load(client.NewClientServiceClient(conn)); err != nil {
		return err
	}
	if err := keystore(); err != nil {
		return err
	}
	if err := truststore(); err != nil {
		return err
	}
	return nil
}

func load(c client.ClientServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h, err := os.Hostname()
	if err != nil {
		return err
	}
	dns, err := podDNS(h)
	if err != nil {
		return err
	}
	r, err := c.NewCert(ctx, &client.CertRequest{Subjects: []string{dns, ip}})
	if err != nil {
		return err
	}
	if err := afero.WriteFile(fs, keyFile, []byte(r.GetKey()), 0644); err != nil {
		return err
	}
	if err := afero.WriteFile(fs, crtFile, []byte(r.GetCert()), 0644); err != nil {
		return err
	}
	return nil
}

func address() string {
	h := os.Getenv(hostEnv)
	if h == "" {
		h = "127.0.0.1"
	}
	p := os.Getenv(portEnv)
	if p == "" {
		p = "10041"
	}
	return h + ":" + p
}

func keystore() error {
	p := "pass:" + Pass
	args := []string{"pkcs12", "-export", "-inkey", keyFile, "-in", crtFile,
		"-out", KeystoreFile, "-name", "cassandra", "-password", p}
	cmd := exec.Command("openssl", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := fs.Chmod(KeystoreFile, 0644); err != nil {
		return err
	}
	return nil
}

func truststore() error {
	args := []string{"-import", "-trustcacerts", "-noprompt", "-alias", "cassandra",
		"-file", caFile, "-keystore", TruststoreFile, "-storepass", Pass}
	cmd := exec.Command(keytool, args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := fs.Chmod(TruststoreFile, 0644); err != nil {
		return err
	}
	return nil
}

func podDNS(h string) (string, error) {
	hn := os.Getenv("DOCKER_HOSTNAME")
	if hn != "" {
		return hn, nil
	}
	r, err := regexp.Compile("-\\d")
	if err != nil {
		return "", err
	}
	s := r.FindString(h)
	if s == "" {
		return "", fmt.Errorf("could not find valid host name for ip %v", ip)
	}
	parts := r.FindStringIndex(h)
	if len(parts) < 1 {
		return "", fmt.Errorf("host should be at least 1, got %v", len(parts))
	}
	i := parts[0]
	runes := []rune(h)
	svc := string(runes[0:i])
	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		return "", errors.New("POD_NAMESPACE must not be empty")
	}
	return fmt.Sprintf("%s.%s.%s.svc.cluster.local", h, svc, ns), nil
}
