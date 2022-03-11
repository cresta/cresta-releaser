package releaser

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/shurcooL/githubv4"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type GitHub interface {
	// CreatePullRequest creates a PR of your current branch.  It assumes there is a remote branch with the
	// exact same name.  It will fail if you're already on master or main.
	CreatePullRequest(ctx context.Context, remoteRepositoryId graphql.ID, baseRefName string, remoteRefName string, title string, body string) error
	// RepositoryInfo returns special information about a remote repository
	RepositoryInfo(ctx context.Context, owner string, name string) (*RepositoryInfo, error)
	// Self returns the current user
	Self(ctx context.Context) (string, error)
}

type RepositoryInfo struct {
	Repository struct {
		ID               githubv4.ID
		DefaultBranchRef struct {
			Name githubv4.String
			ID   githubv4.ID
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type createPullRequest struct {
	CreatePullRequest struct {
		// Note: This is unused, but the library requires at least something to be read for the mutation to happen
		ClientMutationID githubv4.ID
	} `graphql:"createPullRequest(input: $input)"`
}

type GithubGraphqlAPI struct {
	ClientV4 *githubv4.Client
}

type NewGQLClientConfig struct {
	Rt             http.RoundTripper
	AppID          int64
	InstallationID int64
	PEMKeyLoc      string
	Token          string
}

func clientFromToken(_ context.Context, token string) (GitHub, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	gql := githubv4.NewClient(httpClient)
	return &GithubGraphqlAPI{
		ClientV4: gql,
	}, nil
}

func clientFromPEM(ctx context.Context, baseRoundTripper http.RoundTripper, appID int64, installID int64, pemLoc string) (GitHub, error) {
	if baseRoundTripper == nil {
		baseRoundTripper = http.DefaultTransport
	}
	trans, err := ghinstallation.NewKeyFromFile(baseRoundTripper, appID, installID, pemLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to find key file: %w", err)
	}
	_, err = trans.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to validate token: %w", err)
	}
	gql := githubv4.NewClient(&http.Client{Transport: trans})
	return &GithubGraphqlAPI{
		ClientV4: gql,
	}, nil
}

func tokenFromGithubCLI() string {
	s, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configPath := filepath.Join(s, ".config", "gh", "hosts.yml")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ""
	}
	var out map[string]configFileAuths
	if err := yaml.Unmarshal(b, &out); err != nil {
		return ""
	}
	if _, exists := out["github.com"]; !exists {
		return ""
	}
	return out["github.com"].Token
}

type configFileAuths struct {
	Token string `yaml:"oauth_token"`
}

func NewGQLClient(ctx context.Context, cfg *NewGQLClientConfig) (GitHub, error) {
	if cfg != nil && cfg.Token != "" {
		return clientFromToken(ctx, cfg.Token)
	}
	if cfg != nil && cfg.PEMKeyLoc != "" {
		return clientFromPEM(ctx, cfg.Rt, cfg.AppID, cfg.InstallationID, cfg.PEMKeyLoc)
	}
	if os.Getenv("GITHUB_TOKEN") != "" {
		return clientFromToken(ctx, os.Getenv("GITHUB_TOKEN"))
	}
	if token := tokenFromGithubCLI(); token != "" {
		return clientFromToken(ctx, token)
	}
	return nil, fmt.Errorf("no token provided")
}

func (g *GithubGraphqlAPI) Self(ctx context.Context) (string, error) {
	var q struct {
		Viewer struct {
			Login githubv4.String
			ID    githubv4.ID
		}
	}
	if err := g.ClientV4.Query(ctx, &q, nil); err != nil {
		return "", fmt.Errorf("unable to run graphql query: %w", err)
	}
	return string(q.Viewer.Login), nil
}

func (g *GithubGraphqlAPI) CreatePullRequest(ctx context.Context, remoteRepositoryId graphql.ID, baseRefName string, remoteRefName string, title string, body string) error {
	var ret createPullRequest
	if err := g.ClientV4.Mutate(ctx, &ret, githubv4.CreatePullRequestInput{
		RepositoryID: remoteRepositoryId,
		BaseRefName:  githubv4.String(baseRefName),
		HeadRefName:  githubv4.String(remoteRefName),
		Title:        githubv4.String(title),
		Body:         githubv4.NewString(githubv4.String(body)),
	}, nil); err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}
	return nil
}

func (g *GithubGraphqlAPI) RepositoryInfo(ctx context.Context, owner string, name string) (*RepositoryInfo, error) {
	var repoInfo RepositoryInfo
	if err := g.ClientV4.Query(ctx, &repoInfo, map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}); err != nil {
		return nil, fmt.Errorf("unable to query graphql for repository info: %w", err)
	}
	return &repoInfo, nil
}

var _ GitHub = &GithubGraphqlAPI{}
