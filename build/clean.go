package build

import (
	"os"
	"path/filepath"
)

type TargetDirectory struct {
	Path string
}

func (d TargetDirectory) Clean(path string) error {
	return os.RemoveAll(path)
}

func (d TargetDirectory) Dir() string {
	return d.Path
}

func CleanBuild(dir Directory) error {
	d, err := os.Open(dir.Dir())
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = dir.Clean(filepath.Join(dir.Dir(), name))
		if err != nil {
			return err
		}
	}
	return nil
}