package releaser

import (
	"context"
	"sigs.k8s.io/yaml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseConfigMergeFrom(t *testing.T) {
	c1 := `
searchReplace:
  - search: "name"
    replace: "jack"
`
	c2 := `
searchReplace:	
  - search: "job"
    replace: "nothing"
`
	var r1, r2 ReleaseConfig
	require.NoError(t, yaml.Unmarshal([]byte(c1), &r1))
	require.NoError(t, yaml.Unmarshal([]byte(c2), &r2))
	r1.mergeFrom(r2)
	releaserFile := ReleaseFile{
		Name:      "test",
		Directory: ".",
		Content:   "I am name and I do job",
	}
	newContent, err := r1.ApplyToFile(releaserFile, "00-head", "01-staging")
	require.NoError(t, err)
	require.Equal(t, "I am jack and I do nothing", newContent)
}

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

func TestReleaserReadOnly(t *testing.T) {
	WithEmptyExampleApplication(t, func(inst Api) {
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

func RequireRelease(t *testing.T, ctx context.Context, inst Api, application string, release string) {
	prev, newVersion, err := inst.PreviewRelease(ctx, application, release, false)
	require.NoError(t, err)
	require.NoError(t, inst.ApplyRelease(application, release, prev, newVersion))
}

func TestPromotion(t *testing.T) {
	ctx := context.Background()
	NewComplexSetup().WithLayout(ctx, t, func(inst Api) {
		t.Run("promote fully a2", func(t *testing.T) {
			t.Run("promote first", func(t *testing.T) {
				RequireRelease(t, ctx, inst, "a2", "01-staging")
				RequireFileMatches(t, "a2", "01-staging", "config.yaml", "hello world 01-staging YOU-ARE-01-staging")
				t.Run("promote second", func(t *testing.T) {
					RequireRelease(t, ctx, inst, "a2", "02-prod")
					RequireFileMatches(t, "a2", "02-prod", "config.yaml", "replace self YOU-ARE-02-prod")
				})
			})
		})
	})
	NewComplexSetup().WithLayout(ctx, t, func(inst Api) {
		t.Run("promote fully a3", func(t *testing.T) {
			t.Run("promote first", func(t *testing.T) {
				RequireRelease(t, ctx, inst, "a3", "01-staging")
				RequireFileMatches(t, "a3", "01-staging", "config.yaml", "basic promotion 01-staging 01-staging")
				RequireFileMissing(t, "a3", "01-staging", "unknown")
				t.Run("promote second", func(t *testing.T) {
					RequireRelease(t, ctx, inst, "a3", "02-prod")
					RequireFileMatches(t, "a3", "02-prod", "config.yaml", "basic promotion 02-prod 02-prod")
					RequireFileMissing(t, "a3", "02-prod", "unknown")
				})
			})
		})
	})
	NewComplexSetup().WithLayout(ctx, t, func(inst Api) {
		t.Run("promote fully a4", func(t *testing.T) {
			t.Run("promote first", func(t *testing.T) {
				RequireRelease(t, ctx, inst, "a4", "01-staging")
				RequireFileMatches(t, "a4", "01-staging", "config.yaml", "AWESOME PROJECT 01-staging 01-staging")
				RequireFileMissing(t, "a4", "01-staging", "unknown")
			})
		})
	})
}
