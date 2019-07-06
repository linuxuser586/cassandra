package jvm

import (
	"os"
	"testing"
)

func TestYoungGenUserDefined(t *testing.T) {
	os.Setenv(jvmYG, "250m")
	m := Xmn()
	x := "-Xmn250m"
	if m != x {
		t.Errorf("User defined %v was incorrect, got: %v, want: %v", jvmYG, m, x)
	}
}

func TestYoungGenSystemDefined(t *testing.T) {
	os.Setenv(jvmYG, "")
	cpus = 4
	d := []struct {
		m uint64
		x string
	}{
		{1024, "-Xmn128m"},
		{2048, "-Xmn256m"},
		{8000, "-Xmn400m"},
	}
	for _, tbl := range d {
		mem = tbl.m
		r := Xmn()
		if r != tbl.x {
			t.Errorf("System defined young generation memory was incorrect, got: %v, want: %v", r, tbl.x)
		}
	}
}

func TestMimMemUserDefined(t *testing.T) {
	os.Setenv(jvmMem, "512m")
	m := Xms()
	x := "-Xms512m"
	if m != x {
		t.Errorf("User defined %v was incorrect, got: %v, want: %v", jvmMem, m, x)
	}
}

func TestMaxMemUserDefined(t *testing.T) {
	os.Setenv(jvmMem, "700m")
	m := Xmx()
	x := "-Xmx700m"
	if m != x {
		t.Errorf("User defined %v was incorrect, got: %v, want: %v", jvmMem, m, x)
	}
}

func TestMemory(t *testing.T) {
	os.Setenv(jvmMem, "")
	d := []struct {
		m uint64
		t string
		x string
	}{
		{1024, "-Xms", "-Xms512m"},
		{2048, "-Xms", "-Xms1024m"},
		{3788, "-Xms", "-Xms1024m"},
		{4096, "-Xms", "-Xms1024m"},
		{5120, "-Xms", "-Xms1280m"},
		{6144, "-Xms", "-Xms1536m"},
		{10240, "-Xms", "-Xms2560m"},
		{49152, "-Xms", "-Xms8192m"},
	}
	for _, tbl := range d {
		mem = tbl.m
		r := memory(tbl.t)
		if r != tbl.x {
			t.Errorf("System defined memory was incorrect, got: %v, want: %v", r, tbl.x)
		}
	}
}
