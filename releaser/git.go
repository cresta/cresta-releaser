package releaser

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cresta/magehelper/pipe"
	"strings"
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
}

func (g *GitCli) GetRemoteAsGithubRepo(ctx context.Context) (string, string, error) {
	var stdout, stderr bytes.Buffer
	if err := pipe.NewPiped("git", "remote", "get-url", "origin").Execute(ctx, nil, &stdout, &stderr); err != nil {
		return "", "", fmt.Errorf("failed to get remote URL %s: %s", stderr.String(), err)
	}
	// Will either look like
	// * https://github.com/cresta/cresta-releaser.git
	// * git@github.com:cresta/cresta-releaser.git
	remote := stdout.String()
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
	var stdout, stderr bytes.Buffer
	if err := pipe.Shell("git rev-parse --abbrev-ref HEAD").Execute(ctx, nil, &stdout, &stderr); err != nil {
		return "", fmt.Errorf("failed to get current branch name: %w", err)
	}
	if stdout.String() == "" {
		return "", fmt.Errorf("got an empty branch name")
	}
	return stdout.String(), nil
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
	err := pipe.NewPiped("git", "checkout", "-b", branch).Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}
	return nil
}

var _ Git = &GitCli{}
