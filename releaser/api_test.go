package releaser

import (
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
