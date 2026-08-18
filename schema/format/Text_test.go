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

	// Odd and malformed input must not panic or mangle the text
	check("unclosed tag", "<b>unclosed", "unclosed")
	check("stray closing tag", "orphan</b>", "orphan")
	check("bare less-than is not a tag", "5 < 6", "5 < 6")
	check("bare greater-than is not a tag", "6 > 5", "6 > 5")
	check("empty tag", "<>text", "<>text")
	check("html entities are left encoded", "&amp;&lt;&gt;", "&amp;&lt;&gt;")
	check("comment", "a<!-- hidden -->b", "ab")
	check("unicode survives", "héllo <b>wörld</b> 🎉", "héllo wörld 🎉")
	check("invalid utf-8 does not panic", "bad\xff\xfebytes", "bad\xff\xfebytes")
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
