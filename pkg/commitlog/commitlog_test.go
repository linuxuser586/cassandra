package commitlog

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/spf13/afero"
)

const (
	origFile   = "testdata/commitlog_archiving.properties"
	customFile = "testdata/custom.properties"
	mergeFile  = "testdata/merge.properties"
)

var once sync.Once

func before(t *testing.T) {
	once.Do(func() {
		fs = afero.NewMemMapFs()
		loadTestFiles(t, origFile, propFile)
		loadTestFiles(t, customFile, customPropFile)
		loadTestFiles(t, mergeFile, mergeFile)
		// allows to restore to the original
		loadTestFiles(t, origFile, origFile)
		loadTestFiles(t, customFile, customFile)
	})
}

func loadTestFiles(t *testing.T, name, path string) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	afero.WriteFile(fs, path, b, 0644)
}

func reset(t *testing.T) {
	before(t)
	restore(t, origFile, propFile)
	restore(t, customFile, customPropFile)
}

func restore(t *testing.T, name, path string) {
	b, err := afero.ReadFile(fs, name)
	if err != nil {
		t.Fatal(err)
	}
	afero.WriteFile(fs, path, b, 0644)
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

func TestNoCustomProperties(t *testing.T) {
	reset(t)
	fs.Remove(customPropFile)
	Update()
	eq, p1, p2 := equals(t, origFile, propFile)
	if !eq {
		t.Fatalf("want: %v, got: %v", p1, p2)
	}
}

func TestCustomProperties(t *testing.T) {
	reset(t)
	Update()
	eq, p1, p2 := equals(t, propFile, mergeFile)
	if !eq {
		t.Fatalf("want: %v, got: %v", p1, p2)
	}
}
