package releaser

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type FromCommandLine struct {
	fs     FileSystem
	git    Git
	github GitHub
}

func (f *FromCommandLine) MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error {
	owner, repo, err := f.git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get remote as github repo: %w", err)
	}
	return f.github.MergePullRequest(ctx, owner, repo, prNumber)
}

func (f *FromCommandLine) ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error {
	owner, repo, err := f.git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get remote as github repo: %w", err)
	}
	if approvalMessage == "" {
		approvalMessage = "Approved by cresta-releaser"
	}
	return f.github.AcceptPullRequest(ctx, approvalMessage, owner, repo, prNumber)
}

func (f *FromCommandLine) CheckForPROnCurrentBranch(ctx context.Context) (int64, error) {
	branch, err := f.git.CurrentBranchName(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}

	owner, repo, err := f.git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get remote as github repo: %w", err)
	}

	pr, err := f.github.FindPRForBranch(ctx, owner, repo, branch)
	if err != nil {
		return 0, fmt.Errorf("failed to find pr for branch: %w", err)
	}

	return pr, nil
}

func (f *FromCommandLine) GithubWhoami(ctx context.Context) (string, error) {
	return f.github.Self(ctx)
}

func (f *FromCommandLine) PullRequestCurrent(ctx context.Context) error {
	currentBranch, err := f.git.CurrentBranchName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	owner, repo, err := f.git.GetRemoteAsGithubRepo(ctx)
	if err != nil {
		return fmt.Errorf("unable to parse remote URL: %w", err)
	}
	info, err := f.github.RepositoryInfo(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("unable to get repository info for %s/%s: %w", owner, repo, err)
	}
	if err := f.github.CreatePullRequest(ctx, info.Repository.ID, string(info.Repository.DefaultBranchRef.Name), currentBranch, "master", "Update release notes"); err != nil {
		return fmt.Errorf("unable to create pull request: %w", err)
	}
	return nil
}

func (f *FromCommandLine) ForcePushCurrentBranch(ctx context.Context) error {
	currentBranch, err := f.git.CurrentBranchName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	if currentBranch == "master" || currentBranch == "main" {
		return fmt.Errorf("cannot force push master or main branch")
	}
	return f.git.ForcePushHead(ctx, "origin", currentBranch)
}

func (f *FromCommandLine) CommitForRelease(ctx context.Context, application string, release string) error {
	msg := fmt.Sprintf("cresta-releaser: %s:%s", application, release)
	return f.git.CommitAll(ctx, msg)
}

func (f *FromCommandLine) FreshGitBranch(ctx context.Context, application string, release string, forcedName string) error {
	if err := f.git.VerifyFresh(ctx); err != nil {
		return fmt.Errorf("git is not clean: %w", err)
	}
	branchName := forcedName
	if branchName == "" {
		branchName = fmt.Sprintf("creta-release-%s-%s", application, release)
	}
	if err := f.git.CheckoutNewBranch(ctx, branchName); err != nil {
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

func (f *FromCommandLine) isReleaseSymlink(application string, release string) bool {
	releaseDirectory := filepath.Join("apps", application, "releases", release)
	fi, err := os.Stat(releaseDirectory)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
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
	// If this release is a symlink, then we never substitue during a promotion
	if f.isReleaseSymlink(application, release) {
		return thisRelease, thisRelease, nil
	}

	previousReleaseName := releases[thisReleaseIndex-1]
	prevRelease, err := f.GetRelease(application, previousReleaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get previous release %s: %w", previousReleaseName, err)
	}
	nextRelease := describeNewRelease(prevRelease, previousReleaseName, release)
	return thisRelease, nextRelease, nil
}

func filterReleaseOutput(filename string, output string, oldReleaseName string, newReleaseName string) string {
	ret := strings.ReplaceAll(output, oldReleaseName, newReleaseName)
	ret = replaceAutoPromoteTags(filename, output)
	return ret
}

var autoPromoteFilter = regexp.MustCompile(`(.*) # filter-auto-promote .*$`)

func replaceAutoPromoteTags(filename string, output string) string {
	if !strings.HasSuffix(filename, ".yaml") && !strings.HasSuffix(filename, ".yml") {
		return output
	}
	lines := strings.Split(output, "\n")
	newLines := make([]string, 0, len(lines))
	hasMatch := false
	for _, line := range lines {
		if autoPromoteFilter.MatchString(line) {
			hasMatch = true
			newLines = append(newLines, autoPromoteFilter.ReplaceAllString(line, "$1"))
		} else {
			newLines = append(newLines, line)
		}
	}
	if hasMatch {
		return strings.Join(newLines, "\n")
	}
	return output
}

func describeNewRelease(promoteFrom *Release, previousName string, newName string) *Release {
	ret := &Release{}
	for _, f := range promoteFrom.Files {
		ret.Files = append(ret.Files, ReleaseFile{
			Name:    f.Name,
			Content: filterReleaseOutput(f.Name, f.Content, previousName, newName),
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

func NewFromCommandLine(ctx context.Context, logger *zap.Logger, githubCfg *NewGQLClientConfig) (*FromCommandLine, error) {
	gh, err := NewGQLClient(ctx, logger, githubCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	return &FromCommandLine{
		fs: &OSFileSystem{},
		git: &GitCli{
			Logger: logger,
		},
		github: gh,
	}, nil
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
	// FreshGitBranch will create a fresh git branch for releasing
	FreshGitBranch(ctx context.Context, application string, release string, forcedName string) error
	// CommitForRelease will commit the release to the git branch
	CommitForRelease(ctx context.Context, application string, release string) error
	// ForcePushCurrentBranch will force push the current branch to the remote repostiory as a branch with the same name.
	// Fails on branches master or main.
	ForcePushCurrentBranch(ctx context.Context) error
	// PullRequestCurrent creates a pull request for the current branch
	PullRequestCurrent(ctx context.Context) error
	// CheckForPROnCurrentBranch will check if there is a pull request on the current branch
	CheckForPROnCurrentBranch(ctx context.Context) (int64, error)
	// GithubWhoami returns who the CLI thinks you are on github
	GithubWhoami(ctx context.Context) (string, error)
	// ApprovePullRequestForCurrentRemote will approve the pull request on the current remote
	ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error
	// MergePullRequestForCurrentRemote will merge an approved PR
	MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error
}
