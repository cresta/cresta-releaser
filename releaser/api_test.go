package releaser

import (
	"context"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseConfigRegex(t *testing.T) {
	makeTest := func(regex string, oldContent string, newContent string) func(t *testing.T) {
		return func(t *testing.T) {
			x := ReleaseConfig{
				RegexSearchReplace: []regexSearchReplace{
					{
						LineRegexMatch: regex,
						ReplaceWith:    ``,
					},
				},
			}
			seenContent, err := x.ApplyToFile(ReleaseFile{
				Name:      "ignored",
				Directory: "ignored",
				Content:   oldContent,
			}, "previousName", "newName")
			require.NoError(t, err)
			require.Equal(t, newContent, seenContent)
		}
	}
	t.Run("simple", makeTest("empty", "hello", "hello"))
	t.Run("oneline", makeTest("abc", "cabcc", "cc"))
	t.Run("simplemulti", makeTest("a.*b", `
a
file
with:
  content
line abad
line abaaabd`, `
a
file
with:
  content
line ad
line d`))
	t.Run("complexmulti", makeTest(" # .*:autoupdate.*", `
    a-server:
      image:
        tag: master-branch-abcd # a-server:tag-autopush:autoupdate
    b-server:
      image:
        tag: another # not-auto-push
`, `
    a-server:
      image:
        tag: master-branch-abcd
    b-server:
      image:
        tag: another # not-auto-push
`))
}

func WithEmptyExampleApplication(t *testing.T, f func(directory string, inst Api)) {
	dir, err := os.MkdirTemp("", "releaser-test")
	require.NoError(t, err)
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	defer func() {
		require.NoError(t, os.Chdir(currentDir))
		require.NoError(t, os.RemoveAll(dir))
	}()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "apps"), 0755))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "apps", "a1"), 0755))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "apps", "a1", "releases"), 0755))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "apps", "a1", "releases", "00-head"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "apps", "a1", "releases", "00-head", "config.yaml"), []byte(`release\nfrom/00-head`), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "apps", "a1", "releases", "01-staging"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "apps", "a1", "releases", "00-head", "config.yaml"), []byte(``), 0644))
	l := zap.NewProductionConfig()
	logger, err := l.Build()
	require.NoError(t, err)
	inst, err := NewFromCommandLine(context.Background(), logger, &NewGQLClientConfig{
		Token: "unset",
	})
	require.NoError(t, err)
	require.NoError(t, err)
	f(dir, inst)
}

func TestReleaserReadOnly(t *testing.T) {
	WithEmptyExampleApplication(t, func(_ string, inst Api) {
		t.Run("ListApplications", func(t *testing.T) {
			s, err := inst.ListApplications()
			require.NoError(t, err)
			require.Equal(t, []string{"a1"}, s)
		})
		t.Run("ListReleases", func(t *testing.T) {
			t.Run("existing", func(t *testing.T) {
				s, err := inst.ListReleases("a1")
				require.NoError(t, err)
				require.Equal(t, []string{"00-head", "01-staging"}, s)
			})
			t.Run("missing", func(t *testing.T) {
				_, err := inst.ListReleases("missing")
				require.Errorf(t, err, "application missing not found")
			})
		})
	})
}

func CreateApplicationMirrorRelease(t *testing.T) {
	WithEmptyExampleApplication(t, func(_ string, inst Api) {
		t.Run("CreateApplicationMirrorRelease", func(t *testing.T) {
			err := inst.CreateApplicationMirrorRelease("a3", "a1")
			require.NoError(t, err)
			s, err := inst.ListApplications()
			require.NoError(t, err)
			require.Equal(t, []string{"a1", "a3"}, s)
		})
	})
}
