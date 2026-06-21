package html

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// htmlSeeds are shared starting inputs for the HTML-processing fuzzers: well-formed tags, malformed
// tags, unterminated tags, nested and doubled brackets, entities, and multi-byte runes near brackets.
var htmlSeeds = []string{
	"",
	"plain text",
	"<p>hello</p>",
	"<a href=\"x\">link</a>",
	"<<br>>",
	"<unterminated",
	"unopened>",
	"<>",
	"<<>>",
	"a < b > c",
	"<p>日本語です</p>",
	"日本<語>です",
	"<style>body{}</style>",
	"&lt;&gt;&amp;&nbsp;",
	"<br><br><br>",
	"&#169; &copy; &mdash;",
	strings.Repeat("<b>x</b>", 100),
	strings.Repeat("é", 300),
}

// addHTMLSeeds registers the shared seed corpus on a fuzzer.
func addHTMLSeeds(f *testing.F) {
	for _, seed := range htmlSeeds {
		f.Add(seed)
	}
}

// FuzzRemoveTags asserts that RemoveTags strips every angle bracket and never panics or corrupts
// UTF-8. RemoveTags walks byte indices from a rune range loop, so malformed input is the danger zone.
func FuzzRemoveTags(f *testing.F) {
	addHTMLSeeds(f)

	f.Fuzz(func(t *testing.T, html string) {

		result := RemoveTags(html)

		// No angle brackets may survive — that is the whole job of RemoveTags.
		if strings.ContainsAny(result, "<>") {
			t.Fatalf("RemoveTags(%q) = %q still contains an angle bracket", html, result)
		}

		// Valid input must stay valid: the function must not slice through a multi-byte rune.
		// (Invalid input is passed through as-is, which is acceptable.)
		if utf8.ValidString(html) && !utf8.ValidString(result) {
			t.Fatalf("RemoveTags(%q) = %q corrupted valid UTF-8", html, result)
		}
	})
}

// FuzzSummary asserts Summary's length bound and tag-free, valid-UTF-8 output.
func FuzzSummary(f *testing.F) {
	addHTMLSeeds(f)

	f.Fuzz(func(t *testing.T, html string) {

		result := Summary(html)

		// Summary truncates to 200 runes plus an optional "..." ellipsis.
		if runes := utf8.RuneCountInString(result); runes > 203 {
			t.Fatalf("Summary(%q) returned %d runes, want <= 203", html, runes)
		}

		// Summary strips tags, so no angle brackets may survive.
		if strings.ContainsAny(result, "<>") {
			t.Fatalf("Summary(%q) = %q still contains an angle bracket", html, result)
		}

		// Rune-slicing at 200 must never split a multi-byte character of otherwise-valid input.
		if utf8.ValidString(html) && !utf8.ValidString(result) {
			t.Fatalf("Summary(%q) = %q corrupted valid UTF-8", html, result)
		}
	})
}

// FuzzCollapseWhitespace asserts the normalization invariants: no leading/trailing space, no runs of
// two or more spaces, valid UTF-8, and idempotence.
func FuzzCollapseWhitespace(f *testing.F) {
	addHTMLSeeds(f)
	f.Add("   ")
	f.Add("a\t\tb")
	f.Add("\n\n\nx\n\n\n")
	f.Add("trailing   ")

	f.Fuzz(func(t *testing.T, text string) {

		result := CollapseWhitespace(text)

		// No leading or trailing single space should remain after trimming.
		if strings.HasPrefix(result, " ") || strings.HasSuffix(result, " ") {
			t.Fatalf("CollapseWhitespace(%q) = %q has a leading/trailing space", text, result)
		}

		// Whitespace runs must collapse to a single space, so "  " must never appear.
		if strings.Contains(result, "  ") {
			t.Fatalf("CollapseWhitespace(%q) = %q contains a double space", text, result)
		}

		if utf8.ValidString(text) && !utf8.ValidString(result) {
			t.Fatalf("CollapseWhitespace(%q) = %q corrupted valid UTF-8", text, result)
		}

		// Collapsing an already-collapsed string must change nothing.
		if again := CollapseWhitespace(result); again != result {
			t.Fatalf("CollapseWhitespace is not idempotent: %q -> %q -> %q", text, result, again)
		}
	})
}

// FuzzHTMLDecoders is a no-panic + valid-UTF-8 sweep over the remaining HTML helpers, all of which
// accept untrusted markup.
func FuzzHTMLDecoders(f *testing.F) {
	addHTMLSeeds(f)

	f.Fuzz(func(t *testing.T, html string) {

		outputs := map[string]string{
			"ToText":                  ToText(html),
			"ToSearchText":            ToSearchText(html),
			"RemoveSpecialCharacters": RemoveSpecialCharacters(html),
			"RemoveAnchors":           RemoveAnchors(html),
			"Minimal":                 Minimal(html),
		}

		// Valid input must not be corrupted into invalid UTF-8 (invalid input may pass through).
		if utf8.ValidString(html) {
			for name, output := range outputs {
				if !utf8.ValidString(output) {
					t.Fatalf("%s(%q) = %q corrupted valid UTF-8", name, html, output)
				}
			}
		}

		// IsHTML is a predicate; just confirm it never panics on arbitrary input.
		_ = IsHTML(html)
	})
}
