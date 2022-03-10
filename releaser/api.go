package releaser

import (
	"fmt"
	"path/filepath"
)

type FromCommandLine struct {
	fs FileSystem
}

func (f *FromCommandLine) ListApplications() ([]string, error) {
	exists, err := f.fs.DirectoryExists("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to check if apps directory exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("apps directory does not exist")
	}
	dirs, err := f.fs.DirectoriesInsideDirectory("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	return dirs, nil
}

func NewFromCommandLine() *FromCommandLine {
	return &FromCommandLine{
		fs: &OSFileSystem{},
	}
}

func (f *FromCommandLine) ListReleases(application string) ([]string, error) {
	exists, err := f.fs.DirectoryExists("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to check if apps directory exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("apps directory does not exist")
	}
	existsApp, err := f.fs.DirectoryExists(filepath.Join("apps", application))
	if err != nil {
		return nil, fmt.Errorf("failed to check if application directory exists %s: %w", application, err)
	}
	if !existsApp {
		return nil, fmt.Errorf("application %s does not exist", application)
	}
	existsReleases, err := f.fs.DirectoryExists(filepath.Join("apps", application, "releases"))
	if err != nil {
		return nil, fmt.Errorf("failed to check if releases directory exists %s: %w", application, err)
	}
	if !existsReleases {
		return nil, nil
	}
	dirs, err := f.fs.DirectoriesInsideDirectory(filepath.Join("apps", application, "releases"))
	if err != nil {
		return nil, fmt.Errorf("failed to list releases for application %s: %w", application, err)
	}
	return dirs, nil
}

var _ Api = &FromCommandLine{}

type Api interface {
	ListReleases(application string) ([]string, error)
	ListApplications() ([]string, error)
}
