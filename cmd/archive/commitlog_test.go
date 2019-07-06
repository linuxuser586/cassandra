package archive

import (
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/spf13/afero"
)

const (
	origPath   = "/var/lib/cassandra/data/commitlog"
	backupPath = "/var/lib/cassandra/data/backup/commitlog"
	file       = "/in.txt"
	src        = origPath + file
	dst        = backupPath + file
)

var once sync.Once

func before(t *testing.T) {
	once.Do(func() {
		fs = afero.NewMemMapFs()
		link = fakeLink
		loadFile(t, "testdata/in.txt", src)
	})
}

func fakeLink(s, d string) error {
	f, err := fs.Open(s)
	if err != nil {
		return err
	}
	defer f.Close()
	t, err := fs.OpenFile(d, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer t.Close()
	_, err = io.Copy(t, f)
	if err != nil {
		return err
	}
	return nil
}

func loadFile(t *testing.T, f, p string) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	afero.WriteFile(fs, p, b, 0644)
}

func equals(t *testing.T, f1, f2 string) (eq bool, p1 string, p2 string) {
	b1, err := afero.ReadFile(fs, f1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := afero.ReadFile(fs, f2)
	if err != nil {
		t.Fatal(err)
	}
	p1 = string(b1)
	p2 = string(b2)
	return p1 == p2, p1, p2
}

func TestCopy(t *testing.T) {
	before(t)
	if err := Copy(src, dst); err != nil {
		t.Fatal(err)
	}
	eq, p1, p2 := equals(t, src, dst)
	if !eq {
		t.Fatalf("want: %v, got: %v", p1, p2)
	}
}

func TestLink(t *testing.T) {
	before(t)
	if err := Link(src, dst); err != nil {
		t.Fatal(err)
	}
	eq, p1, p2 := equals(t, src, dst)
	if !eq {
		t.Fatalf("want: %v, got: %v", p1, p2)
	}
}
