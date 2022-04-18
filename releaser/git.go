package releaser

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
	ForceDeleteLocalBranch(ctx context.Context, branch string) error
	ChangeOrigin(ctx context.Context, newOrigin string) error
	DoesBranchExist(ctx context.Context, branch string) (bool, error)
	IsAuthorConfigured(ctx context.Context) (bool, error)
	SetLocalAuthor(ctx context.Context, name string, email string) error
	ForceRemoteRefresh(ctx context.Context) error
	CurrentGitSha(ctx context.Context) (string, error)
}

type refreshInterval struct {
	triggerAt time.Time
	interval  time.Duration
	mu        sync.RWMutex
}

func (i *refreshInterval) Execute(ctx context.Context, f func(ctx context.Context) error) error {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if time.Now().Before(i.triggerAt) {
		return nil
	}
	if err := f(ctx); err != nil {
		return fmt.Errorf("failed to execute refresh interval: %w", err)
	}
	interval := i.interval
	if interval == 0 {
		interval = time.Minute
	}
	i.triggerAt = time.Now().Add(interval)
	return nil
}

func (i *refreshInterval) AlwaysExecute(ctx context.Context, f func(ctx context.Context) error) error {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if err := f(ctx); err != nil {
		return fmt.Errorf("failed to execute refresh interval: %w", err)
	}
	interval := i.interval
	if interval == 0 {
		interval = time.Minute
	}
	i.triggerAt = time.Now().Add(interval)
	return nil
}

type GitCli struct {
	Logger       *zap.Logger
	fetchRefresh refreshInterval
}

func (g *GitCli) CurrentGitSha(ctx context.Context) (string, error) {
	var stdout bytes.Buffer
	if err := pipe.Shell("git rev-parse --verify HEAD").Execute(ctx, nil, &stdout, nil); err != nil {
		return "", fmt.Errorf("failed to get current git sha: %w", err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (g *GitCli) ForceRemoteRefresh(ctx context.Context) error {
	return g.refreshWithFunction(ctx, g.fetchRefresh.AlwaysExecute)
}

func (g *GitCli) refreshWithFunction(ctx context.Context, f func(context.Context, func(context.Context) error) error) error {
	return f(ctx, func(ctx context.Context) error {
		return pipe.Shell("git fetch --all -v").Run(ctx)
	})
}

func (g *GitCli) IsAuthorConfigured(ctx context.Context) (bool, error) {
	if err := pipe.Shell("git config --get user.name").Run(ctx); err != nil {
		return false, nil
	}
	if err := pipe.Shell("git config --get user.email").Run(ctx); err != nil {
		return false, nil
	}
	return true, nil
}

func (g *GitCli) SetLocalAuthor(ctx context.Context, name string, email string) error {
	if err := pipe.NewPiped("git", "config", "user.email", email).Run(ctx); err != nil {
		return err
	}
	if err := pipe.NewPiped("git", "config", "user.name", name).Run(ctx); err != nil {
		return err
	}
	return nil
}

func (g *GitCli) DoesBranchExist(ctx context.Context, branch string) (bool, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := pipe.NewPiped("git", "show-ref", "--quiet", "--verify", "refs/heads/"+branch).Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		if strings.Contains(stderr.String(), "fatal:") {
			return false, fmt.Errorf("unable to check if branch exists: %w", err)
		}
		return false, nil
	}
	return true, nil
}

func (g *GitCli) ForceDeleteLocalBranch(ctx context.Context, branch string) error {
	return pipe.NewPiped("git", "branch", "-D", branch).Run(ctx)
}

func (g *GitCli) ChangeOrigin(ctx context.Context, newOrigin string) error {
	return pipe.NewPiped("git", "remote", "set-url", "origin", newOrigin).Run(ctx)
}

func (g *GitCli) ResetToOriginalBranch(ctx context.Context) error {
	if err := g.FetchAllFromRemote(ctx); err != nil {
		return fmt.Errorf("failed to fetch all from remote: %w", err)
	}
	currentBranch, err := g.CurrentBranchName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch name: %w", err)
	}
	if currentBranch != "master" {
		// Ignore error because we don't care if the branch doesn't exist
		_ = pipe.Shell("git branch -D master").Run(ctx)
		if err := pipe.Shell("git checkout -b master origin/master").Run(ctx); err != nil {
			return fmt.Errorf("failed to checkout master: %w", err)
		}
	}
	if err := g.ResetClean(ctx); err != nil {
		return fmt.Errorf("failed to reset clean: %w", err)
	}
	if err := pipe.Shell("git reset --hard origin/master").Run(ctx); err != nil {
		return fmt.Errorf("failed to reset to origin/master: %w", err)
	}
	return nil
}

func (g *GitCli) FetchAllFromRemote(ctx context.Context) error {
	return g.refreshWithFunction(ctx, g.fetchRefresh.Execute)
}

func (g *GitCli) ResetClean(ctx context.Context) error {
	if err := pipe.Shell("git clean -ffdx").Run(ctx); err != nil {
		return fmt.Errorf("git clean failed: %w", err)
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
	remoteURL = strings.TrimPrefix(remoteURL, "git@github.com:")
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
	g.Logger.Debug("AreThereUncommittedChanges")
	defer g.Logger.Debug("AreThereUncommittedChanges done")
	var stdout, stderr bytes.Buffer
	err := pipe.Shell("git status --short").Execute(ctx, nil, &stdout, &stderr)
	g.Logger.Debug("ran git status", zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
	if err != nil {
		return false, fmt.Errorf("git status failed: %w", err)
	}
	if stdout.Len() > 0 {
		return true, nil
	}
	return false, nil
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
		return fmt.Errorf("git checkout failed (%s:%s): %w", stdout.String(), stderr.String(), err)
	}
	return nil
}

var _ Git = &GitCli{}
