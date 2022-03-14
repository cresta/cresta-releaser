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
	VerifyFresh(ctx context.Context) error
	CheckoutNewBranch(ctx context.Context, branch string) error
	CommitAll(ctx context.Context, message string) error
	CurrentBranchName(ctx context.Context) (string, error)
	ForcePushHead(ctx context.Context, repository string, ref string) error
	GetRemoteAsGithubRepo(ctx context.Context) (string, string, error)
}

type GitCli struct {
	Logger *zap.Logger
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
	// Will either look like
	// * https://github.com/cresta/cresta-releaser.git
	// * git@github.com:cresta/cresta-releaser.git
	remote := strings.TrimSpace(stdout.String())
	remote = strings.TrimPrefix(remote, "https://github.com/")
	remote = strings.TrimPrefix(remote, "git@github.com:")
	remote = strings.TrimSuffix(remote, ".git")
	parts := strings.Split(remote, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("failed to parse remote URL %s", remote)
	}
	return parts[0], parts[1], nil
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

func (g *GitCli) VerifyFresh(ctx context.Context) error {
	var stdout, stderr bytes.Buffer
	err := pipe.Shell("git status --short").Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		return fmt.Errorf("git status failed: %w", err)
	}
	if stdout.Len() > 0 {
		return fmt.Errorf("checkout is not fresh: %s", stdout.String())
	}
	return nil
}

func (g *GitCli) CommitAll(ctx context.Context, message string) error {
	var stdout, stderr bytes.Buffer
	if err := pipe.Shell("git add .").Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("git add failed (%s %s): %w", stdout.String(), stderr.String(), err)
	}
	if err := pipe.NewPiped("git", "commit", "-a", "-m", message).Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("git commit failed (%s %s): %w", stdout.String(), stderr.String(), err)
	}
	return nil
}

func (g *GitCli) CheckoutNewBranch(ctx context.Context, branch string) error {
	var stdout, stderr bytes.Buffer
	err := pipe.NewPiped("git", "checkout", "-b", branch, "origin/master").Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}
	return nil
}

var _ Git = &GitCli{}
