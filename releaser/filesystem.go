package releaser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileSystem interface {
	DirectoriesInsideDirectory(dir string) ([]string, error)
	DirectoryExists(dir string) (bool, error)
	FilesInsideDirectory(dir string) ([]File, error)
}

type File struct {
	Name    string
	Content string
	Mode    os.FileMode
}

type OSFileSystem struct {
}

func (O *OSFileSystem) FilesInsideDirectory(dir string) ([]File, error) {
	exists, err := O.DirectoryExists(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to check if directory exists: %s", err)
	}
	if !exists {
		return nil, fmt.Errorf("directory %s does not exist", dir)
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %s", dir, err)
	}
	var ret []File
	for _, ent := range ents {
		if ent.IsDir() {
			continue
		}
		fi, err := ent.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting file info for %s: %s", ent.Name(), err)
		}
		b, err := ioutil.ReadFile(filepath.Join(dir, ent.Name()))
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %s", ent.Name(), err)
		}
		ret = append(ret, File{
			Name:    ent.Name(),
			Content: string(b),
			Mode:    fi.Mode(),
		})
	}
	return ret, nil
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
