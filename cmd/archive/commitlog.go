package archive

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

// makes testing easier
var link = osLink

func osLink(s, d string) error {
	return os.Link(s, d)
}

// Copy the commitlog archive to the backup directory
func Copy(s, d string) error {
	if err := createDir(d); err != nil {
		return err
	}
	b, err := afero.ReadFile(fs, s)
	if err != nil {
		return err
	}
	err = afero.WriteFile(fs, d, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Link the commitlog archive to the backup directory
func Link(s, d string) error {
	if err := createDir(d); err != nil {
		return err
	}
	if err := link(s, d); err != nil {
		return err
	}
	return nil
}

func createDir(p string) error {
	dir := path.Dir(p)
	isDir, err := afero.Exists(fs, dir)
	if err != nil {
		return err
	}
	if !isDir {
		if err := fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
