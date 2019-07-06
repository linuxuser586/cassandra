package cert

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	client "github.com/linuxuser586/apis/grpc/pki/client/v1"

	"github.com/linuxuser586/apis/grpc/pki/client/v1/mock"
	"github.com/linuxuser586/common/pkg/os"
)

var keyContents = "testkey"
var crtContents = "testcrt"

func TestAddress(t *testing.T) {
	d := []struct {
		name string
		exp  string
		env  func(key string) string
	}{
		{"no env defined", "127.0.0.1:10041", func(key string) string {
			return ""
		}},
		{"host only defined", "10.0.1.2:10041", func(key string) string {
			if key == hostEnv {
				return "10.0.1.2"
			}
			return ""
		}},
		{"port only defined", "127.0.0.1:8888", func(key string) string {
			if key == portEnv {
				return "8888"
			}
			return ""
		}},
		{"host and port defined", "10.10.10.10:1024", func(key string) string {
			if key == hostEnv {
				return "10.10.10.10"
			}
			return "1024"
		}},
	}
	for _, tt := range d {
		t.Run(tt.name, func(t *testing.T) {
			os.Getenv = tt.env
			a := address()
			if a != tt.exp {
				t.Errorf("got: %v, want: %v", a, tt.exp)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	os.Hostname = func() (string, error) {
		return "fake-0.host", nil
	}
	os.Getenv = func(key string) string {
		return "ns10"
	}
	podIP = func() (string, error) {
		return "10.11.12.13", nil
	}
	fs = afero.NewMemMapFs()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockClientServiceClient(ctrl)
	c.EXPECT().NewCert(gomock.Any(), gomock.Any()).Return(&client.CertResponse{Key: keyContents, Cert: crtContents}, nil)
	if err := load(c); err != nil {
		t.Fatal(err)
	}
	kf, err := fs.Open(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	defer kf.Close()
	kb, err := afero.ReadAll(kf)
	if err != nil {
		t.Fatal(err)
	}
	cf, err := fs.Open(crtFile)
	if err != nil {
		t.Fatal(err)
	}
	defer cf.Close()
	cb, err := afero.ReadAll(cf)
	if err != nil {
		t.Fatal(err)
	}
	ks := string(kb)
	if ks != keyContents {
		t.Errorf("got: %v, want: %v", ks, keyContents)
	}
	cs := string(cb)
	if cs != crtContents {
		t.Errorf("got: %v, want: %v", cs, crtFile)
	}
}

func TestPodDNS(t *testing.T) {
	os.Getenv = func(key string) string {
		if key == "DOCKER_HOSTNAME" {
			return ""
		}
		return "ns1"
	}
	dns := "cassandra-test-0.cassandra-test.ns1.svc.cluster.local"
	a, err := podDNS("cassandra-test-0")
	if err != nil {
		t.Fatal(err)
	}
	if a != dns {
		t.Errorf("got: %v, want: %v", a, dns)
	}
}
