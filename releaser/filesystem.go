package releaser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type FileSystem interface {
	DirectoriesInsideDirectory(dir string) ([]string, error)
	DirectoryExists(dir string) (bool, error)
	CreateDirectory(dir string) error
	DeleteFile(dir string, name string) error
	ModifyFileContent(dir string, name string, content string) error
	CreateFile(dir string, name string, content string, perms os.FileMode) error
	ChangeFileMode(dir string, name string, perms os.FileMode) error
	FilesInsideDirectory(dir string) ([]File, error)
	ReadFile(dir string, name string) ([]byte, error)
	FileExists(dir string, name string) (bool, error)
}

func IsGitCheckout(fs FileSystem, dir string) bool {
	exists, err := fs.DirectoryExists(filepath.Join(dir, ".git"))
	return exists && err == nil
}

func FilesAtRoot(fs FileSystem, dir string) ([]File, error) {
	var ret []File
	files, err := fs.FilesInsideDirectory(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting files inside directory %s: %w", dir, err)
	}
	ret = append(ret, files...)
	subdirs, err := fs.DirectoriesInsideDirectory(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting subdirectories inside directory %s: %w", dir, err)
	}
	for _, subdir := range subdirs {
		files, err := FilesAtRoot(fs, filepath.Join(dir, subdir))
		if err != nil {
			return nil, fmt.Errorf("error getting files inside directory %s: %w", dir, err)
		}
		ret = append(ret, files...)
	}

	return ret, nil
}

type File struct {
	RelativePath string
	Name         string
	Content      string
	Mode         os.FileMode
}

type OSFileSystem struct {
	Logger *zap.Logger
}

func (O *OSFileSystem) FileExists(dir string, name string) (bool, error) {
	stats, err := os.Stat(filepath.Join(dir, name))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error getting file stats: %w", err)
	}
	return !stats.IsDir(), nil
}

func (O *OSFileSystem) ReadFile(dir string, name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(dir, name))
}

func (O *OSFileSystem) CreateDirectory(dir string) error {
	exists, err := O.DirectoryExists(dir)
	if err != nil {
		return fmt.Errorf("error checking if directory %s exists: %w", dir, err)
	}
	if exists {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	return nil
}

func (O *OSFileSystem) DeleteFile(dir string, name string) error {
	O.Logger.Debug("deleting file", zap.String("dir", dir), zap.String("name", name))
	if err := os.Remove(filepath.Join(dir, name)); err != nil {
		return fmt.Errorf("error deleting file %s: %s", name, err)
	}
	return nil
}

func (O *OSFileSystem) ModifyFileContent(dir string, name string, content string) error {
	O.Logger.Debug("modifying file content", zap.String("dir", dir), zap.String("name", name))
	if err := ioutil.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		return fmt.Errorf("error modifying file %s: %s", name, err)
	}
	return nil
}

func (O *OSFileSystem) CreateFile(dir string, name string, content string, perms os.FileMode) error {
	O.Logger.Debug("creating file", zap.String("dir", dir), zap.String("name", name))
	if err := ioutil.WriteFile(filepath.Join(dir, name), []byte(content), perms); err != nil {
		return fmt.Errorf("error creating file %s: %s", name, err)
	}
	return nil
}

func (O *OSFileSystem) ChangeFileMode(dir string, name string, perms os.FileMode) error {
	O.Logger.Debug("changing file mode", zap.String("dir", dir), zap.String("name", name))
	if err := os.Chmod(filepath.Join(dir, name), perms); err != nil {
		return fmt.Errorf("error changing file mode for file %s: %s", name, err)
	}
	return nil
}

func (O *OSFileSystem) FilesInsideDirectory(dir string) ([]File, error) {
	O.Logger.Debug("getting files inside directory", zap.String("dir", dir))
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
			RelativePath: dir,
			Name:         ent.Name(),
			Content:      string(b),
			Mode:         fi.Mode(),
		})
	}
	return ret, nil
}

func (O *OSFileSystem) DirectoryExists(dir string) (bool, error) {
	O.Logger.Debug("checking if directory exists", zap.String("dir", dir))
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
	O.Logger.Debug("getting directories inside directory", zap.String("dir", dir))
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
