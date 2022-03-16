package mage

import (
	"context"
	"fmt"
	"github.com/cresta/cresta-releaser/cmd/cresta-releaser/commands"
	"github.com/cresta/cresta-releaser/releaser"
	"github.com/magefile/mage/mg"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
	"os"
	"sync"
)

var Instance releaser.Api
var once sync.Once

func MustGetInstance() releaser.Api {
	once.Do(func() {
		l := zap.NewProductionConfig()
		if mg.Verbose() {
			l.Level.SetLevel(zap.DebugLevel)
		}
		logger, err := l.Build()
		if err != nil {
			panic(err)
		}
		var ret releaser.Api
		ret, err = releaser.NewFromCommandLine(context.Background(), logger, nil)
		if err != nil {
			panic(err)
		}
		Instance = ret
	})
	if Instance == nil {
		panic("Unable to make releaser instance")
	}
	return Instance
}

func getOutputFormat() commands.OutputFormatter {
	mageOutput := os.Getenv("MAGE_OUTPUT")
	switch mageOutput {
	case "":
		return &commands.NewlineFormatter{}
	case "auto":
		return &commands.NewlineFormatter{}
	case "json":
		return &commands.JSONFormatter{}
	default:
		panic("Invalid formatter " + mageOutput)
	}
}

// ListReleases will list all releases for an application
func ListReleases(_ context.Context, application string) error {
	releases, err := MustGetInstance().ListReleases(application)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteStringSlice(os.Stdout, releases)
}

// ListApplications will list all applications
func ListApplications(_ context.Context) error {
	apps, err := MustGetInstance().ListApplications()
	if err != nil {
		return err
	}
	return getOutputFormat().WriteStringSlice(os.Stdout, apps)
}

// GetRelease will get a release for an application
func GetRelease(_ context.Context, application string, release string) error {
	out, err := MustGetInstance().GetRelease(application, release)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteObject(os.Stdout, out)
}

// PreviewRelease will show what a new release will look like, promoting from the previous version.  It returns the
// old release and the new release.
func PreviewRelease(_ context.Context, application string, release string) error {
	oldRelease, newRelease, err := MustGetInstance().PreviewRelease(application, release)
	if err != nil {
		return err
	}
	oldContent, newContent := oldRelease.Yaml(), newRelease.Yaml()
	d := diffmatchpatch.New()
	diffs := d.DiffMain(oldContent, newContent, true)
	return getOutputFormat().WriteString(os.Stdout, d.DiffPrettyText(diffs))
}

// ApplyRelease will promote a release to be the current version by applying the previously
// fetched PreviewRelease
func ApplyRelease(_ context.Context, application string, release string) error {
	oldRelease, newRelease, err := MustGetInstance().PreviewRelease(application, release)
	if err != nil {
		return err
	}
	return MustGetInstance().ApplyRelease(application, release, oldRelease, newRelease)
}

// FreshGitBranch will create a fresh git branch for releasing.  The name of the branch will somewhat match the
// release + application name.
func FreshGitBranch(ctx context.Context, application string, release string) error {
	forcedName := os.Getenv("FORCED_NAME")
	return MustGetInstance().FreshGitBranch(ctx, application, release, forcedName)
}

// CommitForRelease will commit the release to the git branch.  It assumes you've already called ApplyRelease
func CommitForRelease(ctx context.Context, application string, release string) error {
	return MustGetInstance().CommitForRelease(ctx, application, release)
}

// ForcePushCurrentBranch will force push the current branch to the remote repository as a branch with the same name.
// Fails on branches master or main.
func ForcePushCurrentBranch(ctx context.Context) error {
	return MustGetInstance().ForcePushCurrentBranch(ctx)
}

// PullRequestCurrent creates a pull request for the current branch
func PullRequestCurrent(ctx context.Context) error {
	return MustGetInstance().PullRequestCurrent(ctx)
}

// CheckForPROnCurrentBranch will check if there is a pull request on the current branch.  Returns 0 if there is no
// PR, otherwise the PR number
func CheckForPROnCurrentBranch(ctx context.Context) error {
	i, err := MustGetInstance().CheckForPROnCurrentBranch(ctx)
	if err != nil {
		return err
	}
	fmt.Println(i)
	return nil
}

// GithubWhoami returns who the CLI thinks you are on github
func GithubWhoami(ctx context.Context) error {
	s, err := MustGetInstance().GithubWhoami(ctx)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteString(os.Stdout, s)
}

// ApprovePullRequestForCurrentRemote will approve the pull request on the current remote
func ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error {
	return MustGetInstance().ApprovePullRequestForCurrentRemote(ctx, approvalMessage, prNumber)
}

// MergePullRequestForCurrentRemote will merge an approved PR
func MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error {
	return MustGetInstance().MergePullRequestForCurrentRemote(ctx, prNumber)
}
