package managedgitrepo

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cresta/cresta-releaser/releaser"
)

type Repo struct {
	DiskLocation string
	URL          string
	Fs           releaser.FileSystem
	Gh           releaser.GitHub
	G            releaser.Git
}

func NewRepo(ctx context.Context, diskLocation string, url string, fs releaser.FileSystem, gh releaser.GitHub, g releaser.Git) (*Repo, error) {
	r := &Repo{
		DiskLocation: diskLocation,
		URL:          url,
		Fs:           fs,
		Gh:           gh,
		G:            g,
	}
	if releaser.IsGitCheckout(r.Fs, diskLocation) {
		return r, r.ResetExistingToOrigin(ctx)
	}
	if err := r.Clone(ctx); err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}
	return r, nil
}

func (r *Repo) VerifyOrSetAuthorInfo(ctx context.Context, name string, email string) error {
	if isSet, err := r.G.IsAuthorConfigured(ctx); err != nil {
		return fmt.Errorf("failed to check if author info is set: %w", err)
	} else if isSet {
		return nil
	}
	if err := r.G.SetLocalAuthor(ctx, name, email); err != nil {
		return fmt.Errorf("failed to set author info: %w", err)
	}
	return nil
}

func (r *Repo) ResetExistingToOrigin(ctx context.Context) error {
	cloneURL, err := r.urlWithToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get clone url: %w", err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current dir: %w", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			panic(err)
		}
	}()
	if os.Chdir(r.DiskLocation) != nil {
		return fmt.Errorf("failed to chdir to repo: %w", err)
	}
	if err := r.G.ChangeOrigin(ctx, cloneURL); err != nil {
		return fmt.Errorf("failed to change origin: %w", err)
	}
	if err := r.G.ResetClean(ctx); err != nil {
		return fmt.Errorf("failed to reset clean: %w", err)
	}
	if err := r.G.ResetToOriginalBranch(ctx); err != nil {
		return fmt.Errorf("failed to reset to original branch: %w", err)
	}
	return nil
}

func (r *Repo) Clone(ctx context.Context) error {
	if isGitCheckout(ctx, r.Fs, r.DiskLocation) {
		return nil
	}
	cloneURL, err := r.urlWithToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get clone url with token added: %w", err)
	}
	return r.G.CloneURL(ctx, cloneURL, r.DiskLocation)
}

func (r *Repo) UpdateCheckout(ctx context.Context) error {
	return r.G.FetchAllFromRemote(ctx)
}

func (r *Repo) urlWithToken(ctx context.Context) (string, error) {
	x, err := url.Parse(r.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %w", r.URL, err)
	}
	token, err := r.Gh.GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	fmt.Println(x.String())
	x.User = url.UserPassword("x-access-token", token)
	fmt.Println(x.String())
	return x.String(), nil
}

func isGitCheckout(_ context.Context, fs releaser.FileSystem, path string) bool {
	exists, err := fs.DirectoryExists(filepath.Join(path, ".git"))
	return err != nil && exists
}
