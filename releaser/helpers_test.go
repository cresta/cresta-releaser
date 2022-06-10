package releaser

import (
	"bytes"
	"context"
	"github.com/cresta/magehelper/pipe"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type RepositoryLayout struct {
	Files              map[string]string
	RepositoryRoot     string
	PreviousWorkingDir string
}

func NewExampleRepository() *RepositoryLayout {
	return &RepositoryLayout{
		Files: map[string]string{
			filepath.Join("apps", "a1", "releases", "00-head", "config.yaml"):    `release\nfrom/00-head`,
			filepath.Join("apps", "a1", "releases", "01-staging", "config.yaml"): ``,
		},
	}
}

func NewComplexSetup() *RepositoryLayout {
	return &RepositoryLayout{
		Files: map[string]string{
			filepath.Join("apps", "a1", "releases", "00-head", "config.yaml"):    `release\nfrom/00-head`,
			filepath.Join("apps", "a1", "releases", "01-staging", "config.yaml"): ``,

			filepath.Join("apps", "a2", "releases", "00-head", "config.yaml"):       `hello world 00-head I-AM-00-head`,
			filepath.Join("apps", "a2", "releases", "00-head", ".releaser.yaml"):    "searchReplace:\n  - search: I-AM\n    replace: YOU-ARE",
			filepath.Join("apps", "a2", "releases", "01-staging", "config.yaml"):    ``,
			filepath.Join("apps", "a2", "releases", "01-staging", ".releaser.yaml"): "searchReplace:\n  - search: hello world 02-prod\n    replace: replace self",
			filepath.Join("apps", "a2", "releases", "02-prod", "config.yaml"):       ``,

			filepath.Join("apps", "a3", "releases", "00-head", "config.yaml"): `basic promotion 00-head 01-staging`,
			filepath.Join("apps", "a3", "releases", "01-staging", "unused"):   ``,
			filepath.Join("apps", "a3", "releases", "02-prod", "unused"):      ``,

			filepath.Join("apps", "a4", "releases", "00-head", "config.yaml"):    `basic promotion 00-head 01-staging`,
			filepath.Join("apps", "a4", ".releaser.yaml"):                        "searchReplace:\n  - search: promotion\n    replace: PROJECT",
			filepath.Join("apps", "a4", "releases", "00-head", ".releaser.yaml"): "searchReplace:\n  - search: basic\n    replace: AWESOME",
			filepath.Join("apps", "a4", "releases", "01-staging", "unused"):      ``,
		},
	}
}

func WithEmptyExampleApplication(t *testing.T, innerFunction func(inst Api)) {
	NewExampleRepository().WithLayout(context.Background(), t, innerFunction)
}

func (d *RepositoryLayout) WithLayout(ctx context.Context, t *testing.T, innerFunction func(inst Api)) {
	defer d.Cleanup(t)
	d.Setup(t)
	l := zap.NewProductionConfig()
	logger, err := l.Build()
	require.NoError(t, err)
	inst, err := NewFromCommandLine(ctx, logger, &NewGQLClientConfig{
		Token: "unset",
	})
	require.NoError(t, err)
	innerFunction(inst)
}

func RequireFileMatches(t *testing.T, application string, release string, filename string, expectedContent string) {
	fullPath := filepath.Join("apps", application, "releases", release, filename)
	content, err := ioutil.ReadFile(fullPath)
	require.NoError(t, err)
	require.Equal(t, expectedContent, string(content))
}

func RequireFileMissing(t *testing.T, application string, release string, filename string) {
	fullPath := filepath.Join("apps", application, "releases", release, filename)
	_, err := os.Stat(fullPath)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func (d *RepositoryLayout) Setup(t *testing.T) {
	require.Empty(t, d.RepositoryRoot)
	dir, err := os.MkdirTemp("", "releaser-test")
	require.NoError(t, err)
	d.RepositoryRoot = dir
	s, err := os.Getwd()
	require.NoError(t, err)
	d.PreviousWorkingDir = s
	require.NoError(t, os.Chdir(dir))

	for path, content := range d.Files {
		dirToMake, _ := filepath.Split(path)
		require.NoError(t, os.MkdirAll(dirToMake, 0755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	}
	MustExec(t, pipe.Shell("git init"))
	MustExec(t, pipe.Shell("git config --local user.email example@example.com"))
	MustExec(t, pipe.Shell("git config --local user.name John Doe"))
	MustExec(t, pipe.Shell("git add ."))
	MustExec(t, pipe.Shell("git commit -am 'init'"))
}

func MustExec(t *testing.T, cmd *pipe.PipedCmd) {
	var stdout, stderr bytes.Buffer
	ctx := context.Background()
	err := cmd.Execute(ctx, nil, &stdout, &stderr)
	if err != nil {
		t.Log(stdout.String())
		t.Log(stderr.String())
		require.NoError(t, err)
	}
}

func (d *RepositoryLayout) Cleanup(t *testing.T) {
	if d.PreviousWorkingDir != "" {
		require.NoError(t, os.Chdir(d.PreviousWorkingDir))
	}
	if d.RepositoryRoot != "" {
		require.NoError(t, os.RemoveAll(d.RepositoryRoot))
	}
}
