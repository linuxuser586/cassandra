package filesystem

import (
	"sync"

	"github.com/spf13/afero"
)

var once sync.Once
var fs afero.Fs

// Singleton is the filesystem singleton.
func Singleton() afero.Fs {
	once.Do(func() {
		fs = afero.NewOsFs()
	})
	return fs
}
