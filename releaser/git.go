package releaser

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/cresta/magehelper/pipe"
	"go.uber.org/zap"
)

type Git interface {
	AreThereUncommittedChanges(ctx context.Context) (bool, error)
	CheckoutNewBranch(ctx context.Context, branch string) error
	CommitAll(ctx context.Context, message string) error
	CurrentBranchName(ctx context.Context) (string, error)
	ForcePushHead(ctx context.Context, repository string, ref string) error
	GetRemoteAsGithubRepo(ctx context.Context) (string, string, error)
	CloneURL(ctx context.Context, url string, into string) error
	ResetClean(ctx context.Context) error
	FetchAllFromRemote(ctx context.Context) error
	ResetToOriginalBranch(ctx context.Context) error
	ChangeOrigin(ctx context.Context, newOrigin string) error
}

type GitCli struct {
	Logger *zap.Logger
}

func (g *GitCli) ChangeOrigin(ctx context.Context, newOrigin string) error {
	return pipe.NewPiped("git", "remote", "set-url", "origin", newOrigin).Run(ctx)
}

func (g *GitCli) ResetToOriginalBranch(ctx context.Context) error {
	if err := pipe.Shell("git checkout master").Run(ctx); err != nil {
		return fmt.Errorf("failed to checkout master: %w", err)
	}
	return g.ResetClean(ctx)
}

func (g *GitCli) FetchAllFromRemote(ctx context.Context) error {
	return pipe.Shell("git fetch --all -v").Run(ctx)
}

func (g *GitCli) ResetClean(ctx context.Context) error {
	if err := pipe.Shell("Git clean -ffdx").Run(ctx); err != nil {
		return fmt.Errorf("Git clean failed: %w", err)
	}
	return pipe.Shell("git reset --hard").Run(ctx)
}

func (g *GitCli) CloneURL(ctx context.Context, url string, into string) error {
	g.Logger.Debug("starting to run command clone")
	defer g.Logger.Debug("done with command clone")
	return pipe.NewPiped("git", "clone", url, into).Run(ctx)
}

func (g *GitCli) runAndLogOutput(ctx context.Context, cmd *pipe.PipedCmd) (bytes.Buffer, bytes.Buffer, error) {
	g.Logger.Debug("starting to run command")
	var stdout, stderr bytes.Buffer
	err := cmd.Execute(ctx, nil, &stdout, &stderr)
	g.Logger.Debug("ran command", zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
	return stdout, stderr, err
}

func (g *GitCli) GetRemoteAsGithubRepo(ctx context.Context) (string, string, error) {
	g.Logger.Debug("GetRemoteAsGithubRepo")
	defer g.Logger.Debug("GetRemoteAsGithubRepo done")
	var stdout, stderr bytes.Buffer
	stdout, stderr, err := g.runAndLogOutput(ctx, pipe.NewPiped("git", "remote", "get-url", "origin"))
	if err != nil {
		return "", "", fmt.Errorf("failed to get remote URL %s: %s", stderr.String(), err)
	}
	remoteURL := stdout.String()
	remoteURL = strings.TrimSpace(remoteURL)
	remoteURL = strings.ToLower(remoteURL)
	remoteURL = strings.TrimSuffix(remoteURL, ".git")
	parts := strings.Split(remoteURL, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("failed to parse remote URL %s", remoteURL)
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}

func (g *GitCli) CurrentBranchName(ctx context.Context) (string, error) {
	var stdout bytes.Buffer
	stdout, _, err := g.runAndLogOutput(ctx, pipe.Shell("git rev-parse --abbrev-ref HEAD"))
	if err != nil {
		return "", fmt.Errorf("failed to get current branch name: %w", err)
	}
	if stdout.String() == "" {
		return "", fmt.Errorf("got an empty branch name")
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (g *GitCli) ForcePushHead(ctx context.Context, repository string, ref string) error {
	var stdout, stderr bytes.Buffer
	if err := pipe.NewPiped("git", "push", "--force", repository, fmt.Sprintf("HEAD:%s", ref)).Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("failed to force push head (%s %s): %w", stdout.String(), stderr.String(), err)
	}
	return nil
}

func (g *GitCli) AreThereUncommittedChanges(ctx context.Context) (bool, error) {
	var stdout, stderr bytes.Buffer
	err := pipe.Shell("git status --short").Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		return false, fmt.Errorf("git status failed: %w", err)
	}
	if stdout.Len() > 0 {
		return false, nil
	}
	return true, nil
}

func (g *GitCli) CommitAll(ctx context.Context, message string) error {
	var stdout, stderr bytes.Buffer
	if err := pipe.Shell("git add .").Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("Git add failed (%s %s): %w", stdout.String(), stderr.String(), err)
	}
	if err := pipe.NewPiped("git", "commit", "-a", "-m", message).Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("Git commit failed (%s %s): %w", stdout.String(), stderr.String(), err)
	}
	return nil
}

func (g *GitCli) CheckoutNewBranch(ctx context.Context, branch string) error {
	var stdout, stderr bytes.Buffer
	err := pipe.NewPiped("git", "checkout", "-b", branch, "origin/master").Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		return fmt.Errorf("Git checkout failed: %w", err)
	}
	return nil
}

var _ Git = &GitCli{}
