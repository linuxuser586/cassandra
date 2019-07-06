package cassandra

import (
	goos "os"
	"testing"
)

func TestSeedHostsNoDomain(t *testing.T) {
	ip = "127.0.0.1"
	setNamespace()
	d := []struct {
		h string
		t string
		x string
	}{
		{"test-host", "IP", ip},
		{"test-host-1", "host not 0 or max", "test-host-0"},
		{"test-host-0.local.srv", "host 0", "test-host-0.local.srv"},
		{"test-host-2.local.srv", "host max", "test-host-0.local.srv"},
		{"test-host-3.local.srv", "host grater than max", "test-host-0.local.srv,test-host-1.local.srv,test-host-2.local.srv"},
	}
	for _, tbl := range d {
		r := seedHosts(tbl.h)
		if r != tbl.x {
			t.Errorf("Test: %v, seed host was incorrect, got: %v, want: %v", tbl.t, r, tbl.x)
		}
	}
}

func TestSeedHostsWithDomain(t *testing.T) {
	ip = "127.0.0.1"
	goos.Setenv("POD_NAMESPACE", "my-test")
	setNamespace()
	d := []struct {
		h string
		t string
		x string
	}{
		{"test-host", "IP", ip},
		{"test-host-1", "host not 0 or max", "test-host-0.test-host.my-test.svc.cluster.local"},
		{"test-host-0", "host 0", "test-host-0.test-host.my-test.svc.cluster.local"},
		{"test-host-2", "host max", "test-host-0.test-host.my-test.svc.cluster.local"},
		{"test-host-3", "host grater than max", "test-host-0.test-host.my-test.svc.cluster.local,test-host-1.test-host.my-test.svc.cluster.local,test-host-2.test-host.my-test.svc.cluster.local"},
	}
	for _, tbl := range d {
		r := seedHosts(tbl.h)
		if r != tbl.x {
			t.Errorf("Test: %v, seed host was incorrect, got: %v, want: %v", tbl.t, r, tbl.x)
		}
	}
}
