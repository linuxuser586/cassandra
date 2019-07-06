package cassandra

import (
	"fmt"
	"io/ioutil"
	"log"
	goos "os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/linuxuser586/cassandra/pkg/cert"
	"github.com/linuxuser586/cassandra/pkg/commitlog"
	"github.com/linuxuser586/cassandra/pkg/config"
	"github.com/linuxuser586/cassandra/pkg/jvm"
	"github.com/linuxuser586/cassandra/pkg/user"
	"github.com/linuxuser586/cassandra/pkg/virtual"
	"github.com/linuxuser586/common/pkg/net"
	"golang.org/x/sys/unix"
	yaml "gopkg.in/yaml.v2"
)

const (
	dataDir     = "/var/lib/cassandra"
	mapCnt      = "/proc/sys/vm/max_map_count"
	customYAML  = "/conf/cassandra.yaml"
	casYAML     = "/etc/cassandra/cassandra.yaml"
	maxSeedHost = 3
	uk          = "unknown"
)

var (
	ip string
	ns string
)

var exec = virtual.NewExec()
var podIP = net.PodIP
var stdout = goos.Stdout
var stderr = goos.Stderr

// Start cassandra
func Start() {
	printHostAndIP()
	setNamespace()
	setRLimits()
	setVM()
	setupDataDirectory()
	setFDPermission()
	updateYAML()
	commitlog.Update()
	run()
}

func printHostAndIP() {
	host, err := goos.Hostname()
	if err != nil {
		host = uk
	}
	fmt.Printf("Host: %v\n", host)
	ip, err = podIP()
	if err != nil {
		fmt.Print(err)
		goos.Exit(1)
	}
	fmt.Printf("IP: %v\n", ip)
}

func setNamespace() {
	ns = goos.Getenv("POD_NAMESPACE")
	if ns == "" {
		ns = uk
	}
	fmt.Printf("Namespace: %v\n", ns)
}

func setRLimits() {
	var lim = syscall.Rlimit{}
	lim.Cur = unix.RLIM_INFINITY
	lim.Max = unix.RLIM_INFINITY
	err := syscall.Setrlimit(unix.RLIMIT_MEMLOCK, &lim)
	if err != nil {
		fmt.Printf("could not set RLIMIT_MEMLOCK: %v\n", err)
		fmt.Println("for Docker, use --cap-add SYS_RESOURCE or --privileged with run")
		fmt.Println("for Kubernetes, add SYS_RESOURCE to the container securityContext capabilities section")
		goos.Exit(1)
	}
	fmt.Printf("RLIMIT_MEMLOCK: %v\n", lim)

	lim.Cur = unix.RLIM_INFINITY
	lim.Max = unix.RLIM_INFINITY
	err = syscall.Setrlimit(syscall.RLIMIT_AS, &lim)
	if err != nil {
		fmt.Fprintf(stderr, "could not set RLIMIT_AS: %v\n", err)
		goos.Exit(1)
	}
	fmt.Printf("RLIMIT_AS: %v\n", lim)

	lim.Cur = 100000
	lim.Max = 100000
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
	if err != nil {
		fmt.Fprintf(stderr, "could not set RLIMIT_NOFILE: %v\n", err)
		goos.Exit(1)
	}
	fmt.Printf("RLIMIT_NOFILE: %v\n", lim)

	lim.Cur = 32768
	lim.Max = 32768
	err = syscall.Setrlimit(unix.RLIMIT_NPROC, &lim)
	if err != nil {
		fmt.Fprintf(stderr, "could not set RLIMIT_NPROC: %v\n", err)
		goos.Exit(1)
	}
	fmt.Printf("RLIMIT_NPROC: %v\n", lim)
}

func setVM() {
	vm := []byte("1048575")
	err := ioutil.WriteFile(mapCnt, vm, 0644)
	if err != nil {
		fmt.Printf("WARN: could not adjust vm.max_map_count: %v\n", err)
		fmt.Println("WARN: this is normal if running in an unprivileged container")
	}

	c, err := ioutil.ReadFile(mapCnt)
	if err != nil {
		fmt.Fprintf(stderr, "could not read vm.max_map_count: %v\n", err)
	}
	i, err := strconv.Atoi(strings.TrimSpace((string(c))))
	if err != nil {
		fmt.Fprintf(stderr, "could not convert vm.max_map_count to integer: %v\n", err)
	}
	if i < 1048575 {
		fmt.Printf("vm.max_map_count is %v and must be greater or equal to 1048575\n", i)
		fmt.Println("for Docker, use --privileged with the run command ")
		fmt.Println("or run sudo sysctl -w vm.max_map_count=1048575 on the host system")
		log.Fatalln("for Kubernetes, run the initContainer")
	} else {
		fmt.Printf("vm.max_map_count: %v", string(c))
	}
}

func setupDataDirectory() {
	if _, err := goos.Stat(dataDir); goos.IsNotExist(err) {
		err = goos.MkdirAll(dataDir, 0755)
		if err != nil {
			fmt.Fprintf(stderr, "could not create data directory: %v\n", err)
			goos.Exit(1)
		}
	}
	err := chownRecursive(dataDir, user.UID, user.GID)
	if err != nil {
		fmt.Fprintf(goos.Stderr, "could not change data directory ownership: %v\n", err)
		goos.Exit(1)
	}
}

func chownRecursive(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info goos.FileInfo, err error) error {
		if err == nil {
			err = goos.Chown(name, uid, gid)
		}
		return err
	})
}

func setFDPermission() {
	err := syscall.Chown("/proc/self/fd/1", user.UID, user.GID)
	if err != nil {
		fmt.Fprintf(goos.Stderr, "could not change fd/1 ownership: %v\n", err)
		goos.Exit(1)
	}
}

func updateYAML() {
	cas, err := ioutil.ReadFile(casYAML)
	if err != nil {
		log.Fatalf("failed to read cassandra.yaml: %v\n", err)
	}
	if _, err := goos.Stat(customYAML); err != nil {
		fmt.Println("using default cassandra.yaml")
		conf := config.CassandraYAML{}
		if err = yaml.Unmarshal(cas, &conf); err != nil {
			log.Fatalf("failed to unmarshal default conf: %v\n", err)
		}
		conf.ListenAddress = ip
		conf.BroadcastRPCAddress = ip
		conf.RPCAddress = "0.0.0.0"
		conf.SeedProvider[0].Parameters[0].Seeds = ip
		o, err := yaml.Marshal(&conf)
		if err != nil {
			log.Fatalf("failed to marshal cassandra.yaml: %v\n", err)
		}
		err = ioutil.WriteFile(casYAML, []byte(o), 0644)
		if err != nil {
			log.Fatalf("failed to write cassandra.yaml: %v\n", err)
		}
		return
	}
	cus, err := ioutil.ReadFile(customYAML)
	if err != nil {
		log.Fatalf("failed to read custom cassandra.yaml: %v\n", err)
	}
	casConf := config.CassandraYAML{}
	if err := yaml.Unmarshal(cas, &casConf); err != nil {
		log.Fatalf("failed to unmarshal cassandra conf: %v\n", err)
	}
	cusConf := config.CassandraYAML{}
	if err := yaml.Unmarshal(cus, &cusConf); err != nil {
		log.Fatalf("failed to unmarshal custom conf: %v\n", err)
	}
	if err := mergo.Merge(&cusConf, casConf); err != nil {
		log.Fatalf("failed to merge configuration: %v\n", err)
	}
	cusConf.ListenAddress = ip
	cusConf.BroadcastRPCAddress = ip
	cusConf.RPCAddress = "0.0.0.0"
	cusConf.ServerEncryptionOptions.Keystore = cert.KeystoreFile
	cusConf.ServerEncryptionOptions.KeystorePassword = cert.Pass
	cusConf.ServerEncryptionOptions.Truststore = cert.TruststoreFile
	cusConf.ServerEncryptionOptions.TruststorePassword = cert.Pass
	cusConf.ServerEncryptionOptions.Protocol = "TLSv1.2"
	cusConf.ServerEncryptionOptions.CipherSuites = []string{"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"}
	cusConf.ClientEncryptionOptions.Keystore = cert.KeystoreFile
	cusConf.ClientEncryptionOptions.KeystorePassword = cert.Pass
	cusConf.ClientEncryptionOptions.Truststore = cert.TruststoreFile
	cusConf.ClientEncryptionOptions.TruststorePassword = cert.Pass
	cusConf.ClientEncryptionOptions.Protocol = "TLSv1.2"
	cusConf.ClientEncryptionOptions.CipherSuites = []string{"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"}
	h, err := goos.Hostname()
	if err != nil {
		log.Fatalf("could not get hostname: %v\n", err)
	}
	cusConf.SeedProvider[0].Parameters[0].Seeds = seedHosts(h)
	o, err := yaml.Marshal(&cusConf)
	if err != nil {
		log.Fatalf("failed to marshal cassandra.yaml: %v\n", err)
	}
	err = ioutil.WriteFile(casYAML, []byte(o), 0644)
	if err != nil {
		log.Fatalf("failed to write cassandra.yaml: %v\n", err)
	}
	log.Printf("Final YAML: %s\n", string(o))
}

func seedHosts(h string) string {
	seeds := goos.Getenv("SEEDS")
	if seeds != "" {
		return seeds
	}
	r, err := regexp.Compile("-\\d")
	if err != nil {
		log.Fatalf("could not create regex: %v\n", err)
	}
	s := r.FindString(h)
	if s == "" {
		log.Printf("could not find valid host name, using %v for seed value\n", ip)
		return ip
	}
	seq := strings.TrimPrefix(s, "-")
	n, err := strconv.Atoi(seq)
	if err != nil {
		log.Fatalf("could not find valid host number, got error %v\n", err)
	}
	domain := ""
	parts := r.FindStringIndex(h)
	if len(parts) < 1 {
		log.Fatalf("could not find valid host part, got error %v\n", err)
	}
	i := parts[0]
	runes := []rune(h)
	svc := string(runes[0:i])
	if ns != uk {
		domain = fmt.Sprintf(".%s.%s.svc.cluster.local", svc, ns)
	}
	seeds = fmt.Sprintf("%s%s", strings.Replace(h, s, "-0", 1), domain)
	if n >= maxSeedHost {
		s1 := fmt.Sprintf("%s%s", strings.Replace(h, s, "-1", 1), domain)
		s2 := fmt.Sprintf("%s%s", strings.Replace(h, s, "-2", 1), domain)
		seeds = fmt.Sprintf("%v,%v,%v", seeds, s1, s2)
	}
	return seeds
}

func run() {
	fmt.Printf("UID: %v, GID: %v\n", user.UID, user.GID)
	cmd := exec.Command("java", javaArgs()...)
	cmd.Credential(&syscall.Credential{Uid: uint32(user.UID), Gid: uint32(user.GID)})
	cmd.Stdout(stdout)
	cmd.Stderr(stderr)
	c := make(chan goos.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Println("shutting down cassandra")
			cmd.Process().Signal(sig)
		}
	}()
	cmd.Run()
	fmt.Println("cassandra stopped")
}

func javaArgs() []string {
	s := jvm.Args()
	s = append(s, "-XX:CompileCommandFile=/etc/cassandra/hotspot_compiler")
	s = append(s, jvm.Xms())
	s = append(s, jvm.Xmx())
	s = append(s, jvm.Xmn())
	s = append(s, "org.apache.cassandra.service.CassandraDaemon")
	return s
}
