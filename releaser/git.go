package releaser

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cresta/magehelper/pipe"
)

type Git interface {
	VerifyFresh(ctx context.Context) error
	CheckoutNewBranch(ctx context.Context, branch string) error
	CommitAll(ctx context.Context, message string) error
}

type GitCli struct {
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
		return fmt.Errorf("git add failed: %w", err)
	}
	if err := pipe.NewPiped("git", "commit", "-a", "-m", message).Execute(ctx, nil, &stdout, &stderr); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
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
