package managedgitrepo

import (
	"context"
	"fmt"
	"net/url"
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
	if err := r.Clone(ctx); err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}
	return r, nil
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
