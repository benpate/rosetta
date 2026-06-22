package html

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRemoveTags_EdgeCases covers malformed and boundary markup that the original tests missed.
func TestRemoveTags_EdgeCases(t *testing.T) {

	cases := map[string]struct {
		input  string
		expect string
	}{
		"empty":             {"", ""},
		"no tags":           {"plain text", "plain text"},
		"unterminated tag":  {"before <unterminated", "before "},
		"unopened tag":      {"a > b", " b"},
		"empty tag":         {"x<>y", "xy"},
		"doubled brackets":  {"<<br>>", ""},
		"adjacent tags":     {"<a><b>text</b></a>", "text"},
		"only a tag":        {"<p>", ""},
		"multibyte content": {"<p>日本語</p>", "日本語"},
	}

	for name, c := range cases {
		result := RemoveTags(c.input)
		require.Equal(t, c.expect, result, name)
		require.NotContains(t, result, "<", name)
		require.NotContains(t, result, ">", name)
	}
}

// TestRemoveSpecialCharacters covers the entity-decoding table directly (previously only exercised
// indirectly through ToText).
func TestRemoveSpecialCharacters(t *testing.T) {

	cases := map[string]struct {
		input  string
		expect string
	}{
		"named lt/gt":    {"&lt;tag&gt;", "<tag>"},
		"numeric lt/gt":  {"&#60;tag&#62;", "<tag>"},
		"ampersand":      {"a &amp; b", "a & b"},
		"nbsp to space":  {"a&nbsp;b", "a b"},
		"copyright":      {"&copy;", "(C)"},
		"em dash":        {"a&mdash;b", "a-b"},
		"smart quotes":   {"&ldquo;hi&rdquo;", `"hi"`},
		"no entities":    {"plain text", "plain text"},
		"unknown entity": {"&unknown;", "&unknown;"},
	}

	for name, c := range cases {
		require.Equal(t, c.expect, RemoveSpecialCharacters(c.input), name)
	}
}

// TestToSearchText confirms tags are replaced with spaces (so adjacent words don't merge) and
// entities are decoded — the behavior that distinguishes it from ToText.
func TestToSearchText(t *testing.T) {

	// Tags are replaced with whitespace so the words on either side stay separate (not "onetwo").
	require.Equal(t, "one two", ToSearchText("one<br>two"))

	// Empty input is handled.
	require.Equal(t, "", ToSearchText(""))

	// Entities are decoded.
	require.Equal(t, "a & b", ToSearchText("a &amp; b"))
}

// TestMinimal confirms the sanitizer keeps allowed formatting tags and strips disallowed/dangerous ones.
func TestMinimal(t *testing.T) {

	// Allowed formatting tags survive.
	require.Equal(t, "<b>bold</b>", Minimal("<b>bold</b>"))
	require.Equal(t, "<p>para</p>", Minimal("<p>para</p>"))

	// Script tags and their contents are removed.
	require.NotContains(t, Minimal("<script>alert(1)</script>"), "<script>")

	// Event-handler attributes are stripped.
	require.NotContains(t, Minimal(`<p onclick="evil()">x</p>`), "onclick")

	// Empty input is handled.
	require.Equal(t, "", Minimal(""))
}

// TestSummary_RuneBoundary confirms truncation counts runes (not bytes) and never splits a
// multi-byte character at the 200-rune boundary.
func TestSummary_RuneBoundary(t *testing.T) {

	// Exactly 200 runes: no truncation, no ellipsis.
	exactly200 := strings.Repeat("a", 200)
	require.Equal(t, exactly200, Summary(exactly200))

	// 201 runes: truncated to 200 plus an ellipsis.
	over := strings.Repeat("a", 201)
	result := Summary(over)
	require.Equal(t, strings.Repeat("a", 200)+"...", result)

	// 201 multi-byte runes must truncate cleanly on a rune boundary (valid UTF-8, 200 runes + "...").
	multibyte := strings.Repeat("é", 201)
	mbResult := Summary(multibyte)
	require.Equal(t, strings.Repeat("é", 200)+"...", mbResult)
}
