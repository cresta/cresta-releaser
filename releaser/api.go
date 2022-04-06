package releaser

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"

	"sigs.k8s.io/yaml"

	"go.uber.org/zap"
	"sigs.k8s.io/kustomize/api/types"
)

// FromCommandLine will process our API using command line execution. It assumes you have things like `Git` already
// installed.
type FromCommandLine struct {
	Fs     FileSystem
	Git    Git
	Github GitHub
	Logger *zap.Logger
}

func (f *FromCommandLine) CreateApplicationFromTemplate(templateDir string, applicationName string, data interface{}) error {
	if exists, err := f.Fs.DirectoryExists(templateDir); err != nil {
		return fmt.Errorf("unable to check if template directory %s exists: %w", templateDir, err)
	} else if !exists {
		return fmt.Errorf("template directory %s does not exist", templateDir)
	}
	applicationDir := filepath.Join("apps", applicationName)
	if exists, err := f.Fs.DirectoryExists(applicationDir); err != nil {
		return fmt.Errorf("unable to check if application directory %s exists: %w", applicationDir, err)
	} else if exists {
		return fmt.Errorf("application directory %s already exists", applicationDir)
	}
	allFiles, err := FilesAtRoot(f.Fs, templateDir)
	if err != nil {
		return fmt.Errorf("unable to get files at root of template directory %s: %w", templateDir, err)
	}
	for _, file := range allFiles {
		fileContent, err := f.Fs.ReadFile(file.RelativePath, file.Name)
		if err != nil {
			return fmt.Errorf("unable to read file %s: %w", filepath.Join(file.RelativePath, file.Name), err)
		}
		t, err := template.New("file").Funcs(sprig.TxtFuncMap()).Parse(string(fileContent))
		if err != nil {
			return fmt.Errorf("unable to parse template for file %s: %w", filepath.Join(file.RelativePath, file.Name), err)
		}
		type templateData struct {
			Name string
			Data interface{}
		}
		var buffer bytes.Buffer
		if err := t.Execute(&buffer, templateData{Name: applicationName, Data: data}); err != nil {
			return fmt.Errorf("unable to execute template for file %s: %w", filepath.Join(file.RelativePath, file.Name), err)
		}
		relPath, err := filepath.Rel(templateDir, file.RelativePath)
		if err != nil {
			return fmt.Errorf("unable to get relative path for file %s: %w", filepath.Join(file.RelativePath, file.Name), err)
		}
		newFileDirectory := filepath.Join(applicationDir, relPath)
		if err := f.Fs.MakeDirectoryAndParents(newFileDirectory); err != nil {
			return fmt.Errorf("unable to make directory %s: %w", file.RelativePath, err)
		}
		fileName, newFileContent := checkForTemplateExtensions(file.Name, buffer.String())
		if err := f.Fs.CreateFile(newFileDirectory, fileName, newFileContent, file.Mode); err != nil {
			return fmt.Errorf("unable to create file %s: %w", filepath.Join(file.RelativePath, file.Name), err)
		}
	}
	return nil
}

func checkForTemplateExtensions(fileName string, fileContent string) (string, string) {
	lines := strings.Split(fileContent, "\n")
	newFilenameRegex := regexp.MustCompile("^#\\s*filename:\\s*(.*)$")
	newFilenameMatch := newFilenameRegex.FindStringSubmatch(lines[0])
	if len(newFilenameMatch) > 1 {
		newFilename := strings.TrimSpace(newFilenameMatch[1])
		lines = lines[1:]
		return newFilename, strings.Join(lines, "\n")
	}
	return fileName, fileContent
}

func (f *FromCommandLine) AreThereUncommittedChanges(ctx context.Context) (bool, error) {
	return f.Git.AreThereUncommittedChanges(ctx)
}

func CheckForPRForRelease(ctx context.Context, a Api, application string, release string) (int64, error) {
	return a.CheckForPRForBranch(ctx, DefaultBranchNameForRelease(application, release))
}

func (f *FromCommandLine) CheckForPRForBranch(ctx context.Context, branchName string) (int64, error) {
	owner, repo, err := f.Git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get remote as Github repo: %w", err)
	}

	pr, err := f.Github.FindPRForBranch(ctx, owner, repo, branchName)
	if err != nil {
		return 0, fmt.Errorf("failed to find pr for branch: %w", err)
	}

	return pr, nil
}

func (f *FromCommandLine) CreateChildApplication(parent string, child string) error {
	doesParentExist, err := DoesApplicationExist(f, parent)
	if err != nil {
		return fmt.Errorf("failed to check if parent application %s exists: %w", parent, err)
	}
	if !doesParentExist {
		return fmt.Errorf("parent application %s does not exist", parent)
	}
	doesChildExist, err := DoesApplicationExist(f, child)
	if err != nil {
		return fmt.Errorf("failed to check if child application %s exists: %w", child, err)
	}
	if doesChildExist {
		return fmt.Errorf("child application %s already exists", child)
	}
	releasesOfParent, err := f.ListReleases(parent)
	if err != nil {
		return fmt.Errorf("failed to list releases of parent application %s: %w", parent, err)
	}
	parentKustomizationFile, err := FindKustomizationForRelease(f, parent, releasesOfParent[0])
	if err != nil {
		return fmt.Errorf("failed to find kustomization for parent application %s: %w", parent, err)
	}
	if parentKustomizationFile == "" {
		return fmt.Errorf("parent application %s does not have a kustomization file", parent)
	}
	if err := f.Fs.CreateDirectory(filepath.Join("apps", child)); err != nil {
		return fmt.Errorf("failed to create child application directory: %w", err)
	}
	const kustomizeFileContent = "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\n"
	if len(releasesOfParent) > 0 {
		if err := f.Fs.CreateDirectory(filepath.Join("apps", child, "releases")); err != nil {
			return fmt.Errorf("failed to create child application directory: %w", err)
		}
		for _, release := range releasesOfParent {
			if err := f.Fs.CreateDirectory(filepath.Join("apps", child, "releases", release)); err != nil {
				return fmt.Errorf("failed to create child application directory release %s:%s: %w", child, release, err)
			}
			if err := f.Fs.CreateFile(filepath.Join("apps", child, "releases", release), "kustomization.yaml", kustomizeFileContent, 0755); err != nil {
				return fmt.Errorf("unable to create kustomization file for child application %s:%s: %w", child, release, err)
			}
		}
	} else {
		if err := f.Fs.CreateFile(filepath.Join("apps", child), "kustomization.yaml", kustomizeFileContent, 0755); err != nil {
			return fmt.Errorf("unable to create kustomization file for child application %s: %w", child, err)
		}
	}
	var parentKustomizationPath string
	var newResourcePath string
	if len(releasesOfParent) > 0 {
		parentKustomizationPath = filepath.Join("apps", parent, "releases", releasesOfParent[0])
		newResourcePath = filepath.Join("..", "..", "..", child, "releases", releasesOfParent[0])
	} else {
		parentKustomizationPath = filepath.Join("apps", parent)
		newResourcePath = filepath.Join("..", child)
	}
	var kc types.Kustomization
	content, err := f.Fs.ReadFile(parentKustomizationPath, parentKustomizationFile)
	if err != nil {
		return fmt.Errorf("failed to read kustomization file for parent application %s: %w", parent, err)
	}
	if err := yaml.UnmarshalStrict(content, &kc); err != nil {
		return fmt.Errorf("failed to unmarshal kustomization file for parent application %s: %w", parent, err)
	}
	if kc.Resources == nil {
		kc.Resources = []string{}
	}
	kc.Resources = append(kc.Resources, newResourcePath)
	sort.Strings(kc.Resources)
	newContent, err := yaml.Marshal(kc)
	if err != nil {
		return fmt.Errorf("failed to marshal kustomization file for parent application %s: %w", parent, err)
	}
	if err := f.Fs.ModifyFileContent(parentKustomizationPath, parentKustomizationFile, string(newContent)); err != nil {
		return fmt.Errorf("failed to modify kustomization file for parent application %s: %w", parent, err)
	}
	return nil
}

func (f *FromCommandLine) MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error {
	owner, repo, err := f.Git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get remote as Github repo: %w", err)
	}
	return f.Github.MergePullRequest(ctx, owner, repo, prNumber)
}

func (f *FromCommandLine) ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error {
	owner, repo, err := f.Git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get remote as Github repo: %w", err)
	}
	if approvalMessage == "" {
		approvalMessage = "Approved by cresta-releaser"
	}
	return f.Github.AcceptPullRequest(ctx, approvalMessage, owner, repo, prNumber)
}

func (f *FromCommandLine) CheckForPROnCurrentBranch(ctx context.Context) (int64, error) {
	branch, err := f.Git.CurrentBranchName(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}
	return f.CheckForPRForBranch(ctx, branch)
}

func (f *FromCommandLine) GithubWhoami(ctx context.Context) (string, error) {
	return f.Github.Self(ctx)
}

func (f *FromCommandLine) PullRequestCurrent(ctx context.Context) (int64, error) {
	currentBranch, err := f.Git.CurrentBranchName(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}
	owner, repo, err := f.Git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to parse remote URL: %w", err)
	}
	info, err := f.Github.RepositoryInfo(ctx, owner, repo)
	if err != nil {
		return 0, fmt.Errorf("unable to get repository info for %s/%s: %w", owner, repo, err)
	}
	if prNum, err := f.Github.CreatePullRequest(ctx, info.Repository.ID, string(info.Repository.DefaultBranchRef.Name), currentBranch, fmt.Sprintf("PR from cresta-releaser for %s", currentBranch), "Deployment"); err != nil {
		return 0, fmt.Errorf("unable to create pull request: %w", err)
	} else {
		return prNum, nil
	}
}

func (f *FromCommandLine) ForcePushCurrentBranch(ctx context.Context) error {
	currentBranch, err := f.Git.CurrentBranchName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	if currentBranch == "master" || currentBranch == "main" {
		return fmt.Errorf("cannot force push master or main branch")
	}
	return f.Git.ForcePushHead(ctx, "origin", currentBranch)
}

func (f *FromCommandLine) CommitForRelease(ctx context.Context, application string, release string) error {
	msg := fmt.Sprintf("cresta-releaser: %s:%s", application, release)
	return f.Git.CommitAll(ctx, msg)
}

func DefaultBranchNameForRelease(application string, release string) string {
	return fmt.Sprintf("releaser-%s-%s", application, release)
}

func (f *FromCommandLine) FreshGitBranch(ctx context.Context, application string, release string, forcedName string) error {
	f.Logger.Debug("Creating new branch for release")
	defer f.Logger.Debug("Created new branch for release")
	if untrackedFiles, err := f.Git.AreThereUncommittedChanges(ctx); err != nil {
		return fmt.Errorf("failed to check for uncommitted changes: %w", err)
	} else if untrackedFiles {
		return fmt.Errorf("there are uncommitted changes")
	}
	branchName := forcedName
	if branchName == "" {
		branchName = DefaultBranchNameForRelease(application, release)
	}
	if err := f.Git.CheckoutNewBranch(ctx, branchName); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}

func (f *FromCommandLine) ApplyRelease(application string, release string, oldRelease *Release, newRelease *Release) error {
	releaseDirectory := filepath.Join("apps", application, "releases", release)
	oldFiles := oldRelease.FilesByName()
	newFiles := newRelease.FilesByName()
	for fileName, file := range oldFiles {
		newContent, exists := newFiles[fileName]
		if !exists {
			if err := f.Fs.DeleteFile(filepath.Join(releaseDirectory, file.Directory), fileName); err != nil {
				return fmt.Errorf("error deleting file %s: %s", fileName, err)
			}
			continue
		}
		if file.Content != newContent.Content {
			if err := f.Fs.ModifyFileContent(filepath.Join(releaseDirectory, file.Directory), fileName, newContent.Content); err != nil {
				return fmt.Errorf("error modifying file %s: %s", fileName, err)
			}
		}
	}
	for fileName, file := range newFiles {
		_, exists := oldFiles[fileName]
		if exists {
			continue
		}
		if err := f.Fs.CreateDirectory(filepath.Join(releaseDirectory, file.Directory)); err != nil {
			return fmt.Errorf("error creating directory %s: %s", file.Directory, err)
		}
		if err := f.Fs.CreateFile(filepath.Join(releaseDirectory, file.Directory), fileName, file.Content, 0744); err != nil {
			return fmt.Errorf("error creating file %s: %s", fileName, err)
		}
	}
	return nil
}

func (f *FromCommandLine) isReleaseSymlink(application string, release string) bool {
	releaseDirectory := filepath.Join("apps", application, "releases", release)
	fi, err := os.Stat(releaseDirectory)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

type searchReplace struct {
	Search  string `yaml:"search"`
	Replace string `yaml:"replace"`
}

type regexSearchReplace struct {
	LineRegexMatch string `yaml:"lineRegexMatch"`
	ReplaceWith    string `yaml:"replaceWith"`
	FileNameMatch  string `yaml:"fileNameMatch"`
}

type ReleaseConfig struct {
	SearchReplace      []searchReplace      `yaml:"searchReplace"`
	RegexSearchReplace []regexSearchReplace `yaml:"regexSearchReplace"`
}

func (c *ReleaseConfig) ApplyToFile(file ReleaseFile, previousReleaseName string, newReleaseName string) (string, error) {
	if c == nil {
		return file.Content, nil
	}
	if file.Name == releaserFileName {
		// Don't replace yourself
		return file.Content, nil
	}
	content := file.Content
	content = strings.ReplaceAll(content, previousReleaseName, newReleaseName)
	for _, sr := range c.SearchReplace {
		content = strings.ReplaceAll(content, sr.Search, sr.Replace)
	}
	for _, rs := range c.RegexSearchReplace {
		filesMatch, err := filepath.Match(rs.FileNameMatch, file.Name)
		if err != nil {
			return "", fmt.Errorf("error matching glob file name: %w", err)
		}
		if rs.FileNameMatch != "" && !filesMatch {
			continue
		}
		re, err := regexp.Compile(rs.LineRegexMatch)
		if err != nil {
			return "", fmt.Errorf("error compiling regex %s: %w", rs.LineRegexMatch, err)
		}
		content = re.ReplaceAllString(content, rs.ReplaceWith)
	}
	return content, nil
}

func (c *ReleaseConfig) mergeFrom(r ReleaseConfig) {
	r.SearchReplace = append(r.SearchReplace, r.SearchReplace...)
	r.RegexSearchReplace = append(r.RegexSearchReplace, r.RegexSearchReplace...)
}

func (f *FromCommandLine) PreviewRelease(application string, release string) (oldRelease *Release, newRelease *Release, err error) {
	f.Logger.Debug("previewing release")
	defer f.Logger.Debug("previewed release")
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
	// If this release is a symlink, then we never promote
	if f.isReleaseSymlink(application, release) {
		return thisRelease, thisRelease, nil
	}

	previousReleaseName := releases[thisReleaseIndex-1]
	prevRelease, err := f.GetRelease(application, previousReleaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get previous release %s: %w", previousReleaseName, err)
	}
	promotionConfig, err := ReleaseConfigForRelease(f.Fs, application, previousReleaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get promotion config for release %s: %w", previousReleaseName, err)
	}
	f.Logger.Debug("promotion config", zap.Any("config", promotionConfig))
	nextRelease, err := describeNewRelease(prevRelease, previousReleaseName, release, promotionConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to describe new release: %w", err)
	}
	return thisRelease, nextRelease, nil
}

const releaserFileName = ".releaser.yaml"

func ReleaseConfigForRelease(fs FileSystem, application string, release string) (*ReleaseConfig, error) {
	possibleConfigPaths := []string{
		filepath.Join("apps"),
		filepath.Join("apps", application),
		filepath.Join("apps", application, "releases", release),
	}
	var ret *ReleaseConfig
	for _, p := range possibleConfigPaths {
		exists, err := fs.FileExists(p, releaserFileName)
		if err != nil {
			return nil, fmt.Errorf("unable to check if %s:%s exists: %w", p, releaserFileName, err)
		}
		if exists {
			contents, err := fs.ReadFile(p, releaserFileName)
			if err != nil {
				return nil, fmt.Errorf("unable to read %s:%s: %w", p, releaserFileName, err)
			}
			var r ReleaseConfig
			if err := yaml.Unmarshal(contents, &r); err != nil {
				return nil, fmt.Errorf("unable to parse %s:%s .releaser.yaml: %w", p, releaserFileName, err)
			}
			if ret == nil {
				ret = &r
			} else {
				ret.mergeFrom(r)
			}
		}
	}

	return ret, nil
}

func describeNewRelease(promoteFrom *Release, previousName string, newName string, releaseConfig *ReleaseConfig) (*Release, error) {
	ret := &Release{}
	for _, f := range promoteFrom.Files {
		newContent, err := releaseConfig.ApplyToFile(f, previousName, newName)
		if err != nil {
			return nil, fmt.Errorf("unable to apply promotion config to file %s: %w", f.Name, err)
		}
		ret.Files = append(ret.Files, ReleaseFile{
			Name:      f.Name,
			Content:   newContent,
			Directory: f.Directory,
		})
	}
	return ret, nil
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
	exists, err := f.Fs.DirectoryExists("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to check if apps directory exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("apps directory does not exist")
	}
	existsApp, err := f.Fs.DirectoryExists(filepath.Join("apps", application))
	if err != nil {
		return nil, fmt.Errorf("failed to check if application directory exists %s: %w", application, err)
	}
	if !existsApp {
		return nil, fmt.Errorf("application %s does not exist", application)
	}
	existsReleases, err := f.Fs.DirectoryExists(filepath.Join("apps", application, "releases"))
	if err != nil {
		return nil, fmt.Errorf("failed to check if releases directory exists %s: %w", application, err)
	}
	if !existsReleases {
		if release == "" {
			return f.releaseInPath(filepath.Join("apps", application))
		}
		return nil, fmt.Errorf("releases directory does not exist for application %s", application)
	}
	existsRelease, err := f.Fs.DirectoryExists(filepath.Join("apps", application, "releases", release))
	if err != nil {
		return nil, fmt.Errorf("failed to check if existing release directory exists %s: %w", application, err)
	}
	if !existsRelease {
		return nil, fmt.Errorf("release %s does not exist", release)
	}
	return f.releaseInPath(filepath.Join("apps", application, "releases", release))
}

func (f *FromCommandLine) releaseInPath(path string) (*Release, error) {
	files, err := FilesAtRoot(f.Fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get files inside release directory %s: %w", path, err)
	}
	releaseFiles := make([]ReleaseFile, 0)
	for _, f := range files {
		relPath, err := filepath.Rel(path, f.RelativePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for file %s: %w", f.RelativePath, err)
		}
		releaseFiles = append(releaseFiles, ReleaseFile{
			Name:      f.Name,
			Content:   f.Content,
			Directory: relPath,
		})
	}
	return &Release{Files: releaseFiles}, nil
}

func (f *FromCommandLine) ListApplications() ([]string, error) {
	exists, err := f.Fs.DirectoryExists("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to check if apps directory exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("apps directory does not exist")
	}
	dirs, err := f.Fs.DirectoriesInsideDirectory("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	return dirs, nil
}

func NewFromCommandLine(ctx context.Context, logger *zap.Logger, githubCfg *NewGQLClientConfig) (*FromCommandLine, error) {
	gh, err := NewGQLClient(ctx, logger, githubCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Github client: %w", err)
	}
	return &FromCommandLine{
		Logger: logger,
		Fs: &OSFileSystem{
			Logger: logger,
		},
		Git: &GitCli{
			Logger: logger,
		},
		Github: gh,
	}, nil
}

func (f *FromCommandLine) ListReleases(application string) ([]string, error) {
	exists, err := f.Fs.DirectoryExists("apps")
	if err != nil {
		return nil, fmt.Errorf("failed to check if apps directory exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("apps directory does not exist")
	}
	existsApp, err := f.Fs.DirectoryExists(filepath.Join("apps", application))
	if err != nil {
		return nil, fmt.Errorf("failed to check if application directory exists %s: %w", application, err)
	}
	if !existsApp {
		return nil, fmt.Errorf("application %s does not exist", application)
	}
	existsReleases, err := f.Fs.DirectoryExists(filepath.Join("apps", application, "releases"))
	if err != nil {
		return nil, fmt.Errorf("failed to check if releases directory exists %s: %w", application, err)
	}
	if !existsReleases {
		return nil, nil
	}
	dirs, err := f.Fs.DirectoriesInsideDirectory(filepath.Join("apps", application, "releases"))
	if err != nil {
		return nil, fmt.Errorf("failed to list releases for application %s: %w", application, err)
	}
	return dirs, nil
}

var _ Api = &FromCommandLine{}

// A Release is a collection of files that we intend to change in Git
type Release struct {
	// Files is each released file
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
	// Name of the file (filename)
	Name string
	// Relative path to the file
	Directory string
	// Content of the file
	Content string
}

// Api is an interface into our release process.
type Api interface {
	// CreateChildApplication creates a child application that releases with the same cadence as the parent.  It assumes
	// the child does not exist.
	CreateChildApplication(parent string, child string) error
	// CreateApplicationFromTemplate creates a new application from a go template directory
	CreateApplicationFromTemplate(templateDir string, applicationName string, data interface{}) error
	// ListReleases will list all releases for an application
	ListReleases(application string) ([]string, error)
	// ListApplications will list all applications
	ListApplications() ([]string, error)
	// GetRelease will get a release for an application
	GetRelease(application string, release string) (*Release, error)
	// PreviewRelease will show what a new release will look like, promoting from the previous version.  It returns the
	// old release and the new release.
	PreviewRelease(application string, release string) (*Release, *Release, error)
	// ApplyRelease will promote a release to be the current version by applying the previously
	// fetched PreviewRelease
	ApplyRelease(application string, release string, oldRelease *Release, newRelease *Release) error
	// FreshGitBranch will create a fresh Git branch for releasing.  The name of the branch will somewhat match the
	// release + application name.
	FreshGitBranch(ctx context.Context, application string, release string, forcedName string) error
	// CommitForRelease will commit the release to the Git branch.  It assumes you've already called ApplyRelease
	CommitForRelease(ctx context.Context, application string, release string) error
	// AreThereUncommittedChanges will check if there are any uncommitted changes in the Git branch.
	AreThereUncommittedChanges(ctx context.Context) (bool, error)
	// ForcePushCurrentBranch will force push the current branch to the remote repository as a branch with the same name.
	// Fails on branches master or main.
	ForcePushCurrentBranch(ctx context.Context) error
	// PullRequestCurrent creates a pull request for the current branch
	PullRequestCurrent(ctx context.Context) (int64, error)
	// CheckForPROnCurrentBranch will check if there is a pull request on the current branch.  Returns 0 if there is no
	// PR, otherwise the PR number
	CheckForPROnCurrentBranch(ctx context.Context) (int64, error)
	// GithubWhoami returns who the CLI thinks you are on Github
	GithubWhoami(ctx context.Context) (string, error)
	// ApprovePullRequestForCurrentRemote will approve the pull request on the current remote
	ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error
	// MergePullRequestForCurrentRemote will merge an approved PR
	MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error
	// CheckForPRForBranch returns the PR number for a branch of the current Git repository
	CheckForPRForBranch(ctx context.Context, branchName string) (int64, error)
}
