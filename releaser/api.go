package releaser

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

type FromCommandLine struct {
	fs FileSystem
}

func (f *FromCommandLine) ApplyRelease(application string, release string, oldRelease *Release, newRelease *Release) error {
	releaseDirectory := filepath.Join("apps", application, "releases", release)
	oldFiles := oldRelease.FilesByName()
	newFiles := newRelease.FilesByName()
	for fileName, file := range oldFiles {
		newContent, exists := newFiles[fileName]
		if !exists {
			if err := f.fs.DeleteFile(releaseDirectory, fileName); err != nil {
				return fmt.Errorf("error deleting file %s: %s", fileName, err)
			}
			continue
		}
		if file.Content != newContent.Content {
			if err := f.fs.ModifyFileContent(releaseDirectory, fileName, newContent.Content); err != nil {
				return fmt.Errorf("error modifying file %s: %s", fileName, err)
			}
		}
	}
	for fileName, file := range newFiles {
		_, exists := oldFiles[fileName]
		if exists {
			continue
		}
		if err := f.fs.CreateFile(releaseDirectory, fileName, file.Content, 0744); err != nil {
			return fmt.Errorf("error deleting file %s: %s", fileName, err)
		}
	}
	return nil
}

func (f *FromCommandLine) PreviewRelease(application string, release string) (oldRelease *Release, newRelease *Release, err error) {
	releases, err := f.ListReleases(application)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to list releases: %w", err)
	}
	thisReleaseIndex := indexOf(release, releases)
	if thisReleaseIndex == -1 {
		return nil, nil, fmt.Errorf("release %s not found", release)
	}
	if thisReleaseIndex == 0 {
		return nil, nil, fmt.Errorf("cannot preview the original release")
	}

	thisRelease, err := f.GetRelease(application, release)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get release %s: %w", release, err)
	}

	previousReleaseName := releases[thisReleaseIndex-1]
	prevRelease, err := f.GetRelease(application, previousReleaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get previous release %s: %w", previousReleaseName, err)
	}
	nextRelease := describeNewRelease(prevRelease, previousReleaseName, release)
	return thisRelease, nextRelease, nil
}

func describeNewRelease(promoteFrom *Release, previousName string, newName string) *Release {
	ret := &Release{}
	for _, f := range promoteFrom.Files {
		ret.Files = append(ret.Files, ReleaseFile{
			Name:    f.Name,
			Content: strings.ReplaceAll(f.Content, previousName, newName),
		})
	}
	return ret
}

func indexOf(s string, in []string) int {
	for i, v := range in {
		if v == s {
			return i
		}
	}
	return -1
}

func (f *FromCommandLine) GetRelease(application string, release string) (*Release, error) {
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
		if release == "" {
			return f.releaseInPath(filepath.Join("apps", application))
		}
		return nil, fmt.Errorf("releases directory does not exist for application %s", application)
	}
	existsRelease, err := f.fs.DirectoryExists(filepath.Join("apps", application, "releases", release))
	if err != nil {
		return nil, fmt.Errorf("failed to check if existing release directory exists %s: %w", application, err)
	}
	if !existsRelease {
		return nil, fmt.Errorf("release %s does not exist", release)
	}
	return f.releaseInPath(filepath.Join("apps", application, "releases", release))
}

func (f *FromCommandLine) releaseInPath(path string) (*Release, error) {
	files, err := f.fs.FilesInsideDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get files inside release directory %s: %w", path, err)
	}
	releaseFiles := make([]ReleaseFile, 0)
	for _, f := range files {
		releaseFiles = append(releaseFiles, ReleaseFile{
			Name:    f.Name,
			Content: f.Content,
		})
	}
	return &Release{Files: releaseFiles}, nil
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

type Release struct {
	Files []ReleaseFile
}

func (r *Release) Yaml() string {
	var b bytes.Buffer
	for idx, f := range r.Files {
		if idx != 0 {
			b.WriteString("---\n")
		}
		b.WriteString(f.Content)
		b.WriteString("\n")
	}
	return b.String()
}

func (r *Release) FilesByName() map[string]ReleaseFile {
	files := make(map[string]ReleaseFile)
	for _, f := range r.Files {
		files[f.Name] = f
	}
	return files
}

type ReleaseFile struct {
	Name    string
	Content string
}

type Api interface {
	// ListReleases will list all releases for an application
	ListReleases(application string) ([]string, error)
	// ListApplications will list all applications
	ListApplications() ([]string, error)
	// GetRelease will get a release for an application
	GetRelease(application string, release string) (*Release, error)
	// PreviewRelease will show what a new release will look like, promoting from the previous version
	PreviewRelease(application string, release string) (*Release, *Release, error)
	// ApplyRelease will promote a release to be the current version by applying the previously
	// fetched PreviewRelease
	ApplyRelease(application string, release string, oldRelease *Release, newRelease *Release) error
}
