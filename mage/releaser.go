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

func ListReleases(_ context.Context, application string) error {
	releases, err := MustGetInstance().ListReleases(application)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteStringSlice(os.Stdout, releases)
}

func ListApplications(_ context.Context) error {
	apps, err := MustGetInstance().ListApplications()
	if err != nil {
		return err
	}
	return getOutputFormat().WriteStringSlice(os.Stdout, apps)
}

func GetRelease(_ context.Context, application string, release string) error {
	out, err := MustGetInstance().GetRelease(application, release)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteObject(os.Stdout, out)
}

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

func ApplyRelease(_ context.Context, application string, release string) error {
	oldRelease, newRelease, err := MustGetInstance().PreviewRelease(application, release)
	if err != nil {
		return err
	}
	return MustGetInstance().ApplyRelease(application, release, oldRelease, newRelease)
}

func FreshGitBranch(ctx context.Context, application string, release string) error {
	forcedName := os.Getenv("FORCED_NAME")
	return MustGetInstance().FreshGitBranch(ctx, application, release, forcedName)
}

func CommitForRelease(ctx context.Context, application string, release string) error {
	return MustGetInstance().CommitForRelease(ctx, application, release)
}

func ForcePushCurrentBranch(ctx context.Context) error {
	return MustGetInstance().ForcePushCurrentBranch(ctx)
}

func PullRequestCurrent(ctx context.Context) error {
	return MustGetInstance().PullRequestCurrent(ctx)
}

func CheckForPROnCurrentBranch(ctx context.Context) error {
	i, err := MustGetInstance().CheckForPROnCurrentBranch(ctx)
	if err != nil {
		return err
	}
	fmt.Println(i)
	return nil
}

func GithubWhoami(ctx context.Context) error {
	s, err := MustGetInstance().GithubWhoami(ctx)
	if err != nil {
		return err
	}
	return getOutputFormat().WriteString(os.Stdout, s)
}

func ApprovePullRequestForCurrentRemote(ctx context.Context, approvalMessage string, prNumber int64) error {
	return MustGetInstance().ApprovePullRequestForCurrentRemote(ctx, approvalMessage, prNumber)
}

func MergePullRequestForCurrentRemote(ctx context.Context, prNumber int64) error {
	return MustGetInstance().MergePullRequestForCurrentRemote(ctx, prNumber)
}
