package releaser

import (
	"fmt"
	"os"
)

type FileSystem interface {
	DirectoriesInsideDirectory(dir string) ([]string, error)
	DirectoryExists(dir string) (bool, error)
}

type OSFileSystem struct {
}

func (O *OSFileSystem) DirectoryExists(dir string) (bool, error) {
	f, err := os.Stat(dir)
	if err == nil {
		if f.IsDir() {
			return true, nil
		}
		return false, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (O *OSFileSystem) DirectoriesInsideDirectory(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %s", dir, err)
	}
	var ret []string
	for _, ent := range ents {
		if ent.IsDir() {
			ret = append(ret, ent.Name())
		}
	}
	return ret, nil
}

var _ FileSystem = &OSFileSystem{}
