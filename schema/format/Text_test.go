package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestText(t *testing.T) {

	format := Text("")

	check := func(name string, input string, expected string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			result, err := format(input)
			require.NoError(t, err)
			require.Equal(t, expected, result)
		})
	}

	check("plain text passes through untouched", "hello there", "hello there")
	check("empty string", "", "")
	check("strips a simple tag pair", "<b>bold</b>", "bold")
	check("strips a script element's tags but keeps its text", "<script>alert(1)</script>", "alert(1)")
	check("strips attributes along with the tag", `<a href="https://example.com">link</a>`, "link")
	check("strips nested tags", "<div><p><b>deep</b></p></div>", "deep")
	check("strips a self-closing tag", "before<br/>after", "beforeafter")

	// Unlike NoHTML, Text does NOT collapse whitespace -- that is the whole point of it.
	check("preserves runs of spaces", "a    b", "a    b")
	check("preserves newlines", "line1\nline2", "line1\nline2")
	check("preserves tabs", "a\tb", "a\tb")
	check("preserves whitespace left behind by a stripped tag", "<b>a</b>   <i>b</i>", "a   b")
	check("preserves leading and trailing whitespace", "  padded  ", "  padded  ")

	// Odd and malformed input must not panic
	check("unclosed tag", "<b>unclosed", "unclosed")
	check("stray closing tag", "orphan</b>", "orphan")
	check("empty tag", "<>text", "text")
	check("html entities are left encoded", "&amp;&lt;&gt;", "&amp;&lt;&gt;")
	check("comment", "a<!-- hidden -->b", "ab")
	check("unicode survives", "héllo <b>wörld</b> 🎉", "héllo wörld 🎉")
	check("invalid utf-8 does not panic", "bad\xff\xfebytes", "bad\xff\xfebytes")
}

// Unescaped angle brackets in plain text are MANGLED by the underlying html.RemoveTags
// scanner, not passed through. These cases are pinned so the damage is visible and so a
// fix to RemoveTags shows up here as a failure -- they are recording current behavior,
// NOT endorsing it. See TestText_AngleBracketDataLoss for the details.
func TestText_UnescapedAngleBrackets(t *testing.T) {

	format := Text("")

	check := func(name string, input string, expected string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			result, err := format(input)
			require.NoError(t, err)
			require.Equal(t, expected, result)
		})
	}

	// A bare `<` opens a tag that never closes, so the rest of the string is dropped.
	check("bare less-than drops the remainder", "5 < 6", "5 ")

	// A bare `>` resets the scanner's write cursor, dropping everything BEFORE it.
	check("bare greater-than drops the prefix", "6 > 5", " 5")
	check("two bare greater-thans drop even more", "a > b > c", " c")

	// A doubled `<` writes the preceding text twice.
	check("doubled less-than duplicates the prefix", "a<<br>b", "aab")
}

// Spells out the two html.RemoveTags defects that TestText_UnescapedAngleBrackets pins,
// so that whoever reads a failure there knows these are known bugs rather than intent.
//
// Both are reachable from user-supplied content, because Text is a schema string format.
// Neither is a sanitizer hole -- tags are still stripped -- but both corrupt plain text.
func TestText_AngleBracketDataLoss(t *testing.T) {

	format := Text("")

	// DEFECT 1: a bare `>` discards all text since the last closed tag.
	result, err := format("2 > 1 is true")
	require.NoError(t, err)
	require.NotEqual(t, "2 > 1 is true", result, "if this now passes, DEFECT 1 was fixed")
	require.Equal(t, " 1 is true", result, "KNOWN BUG: the leading \"2 \" is silently lost")

	// DEFECT 2: a doubled `<` emits the preceding text twice.
	result, err = format("x<<b>y")
	require.NoError(t, err)
	require.NotEqual(t, "xy", result, "if this now passes, DEFECT 2 was fixed")
	require.Equal(t, "xxy", result, "KNOWN BUG: the leading \"x\" is duplicated")
}

// The `arg` parameter is unused by Text -- every value must behave identically.
func TestText_IgnoresArg(t *testing.T) {

	for _, arg := range []string{"", "anything", "0", "true", "  "} {
		result, err := Text(arg)("<b>x</b>")
		require.NoError(t, err)
		require.Equal(t, "x", result, "arg %q must not change the result", arg)
	}
}

// Text never reports an error, for any input.
func TestText_NeverErrors(t *testing.T) {

	format := Text("")

	for _, input := range []string{"", "plain", "<b>x", "\x00", "<<<>>>", "\xff"} {
		_, err := format(input)
		require.NoError(t, err, "input %q must not error", input)
	}
}

// Stripping tags is idempotent: running Text over its own output changes nothing.
func TestText_IsIdempotent(t *testing.T) {

	format := Text("")

	for _, input := range []string{"<b>bold</b>", "a<br/>b", "plain   text", "<div><p>x</p></div>"} {
		once, err := format(input)
		require.NoError(t, err)

		twice, err := format(once)
		require.NoError(t, err)
		require.Equal(t, once, twice, "Text must be idempotent for %q", input)
	}
}
