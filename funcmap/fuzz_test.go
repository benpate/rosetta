package funcmap

import (
	"html/template"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzSafeAttr asserts the invariant behind the `attr` template function: whatever it
// returns is emitted into an HTML attribute with NO escaping, so the result must either
// be empty, or be the input unchanged AND free of every character that could close the
// attribute, open a tag, or start an adjacent one.
func FuzzSafeAttr(f *testing.F) {

	f.Add("button-primary")
	f.Add("")
	f.Add(`" onmouseover="alert(1)`)
	f.Add("a b")
	f.Add("a=b")
	f.Add("<script>")
	f.Add("&amp;")
	f.Add("a\tb")
	f.Add("a b")
	f.Add("héllo")
	f.Add("\x00")

	f.Fuzz(func(t *testing.T, value string) {

		result := string(safeAttr(value))

		if result == "" {
			return
		}

		// A non-empty result must be the input, verbatim and unmodified
		require.Equal(t, value, result)

		// ...and must not carry a single attribute-breakout character
		for _, character := range result {
			require.False(t, isUnsafeAttrRune(character), "unsafe rune %q survived in %q", character, result)
		}

		require.False(t, strings.ContainsAny(result, "\"'`<>=& \t\r\n"))

		// The guard must be idempotent
		require.Equal(t, result, string(safeAttr(result)))
	})
}

// FuzzSafeCSSValue asserts the invariant behind the `cssValue` template function: the
// result is emitted as trusted CSS, so it must either be empty, or be the input unchanged,
// hold only allowlisted characters, and never name the two dangerous CSS functions.
func FuzzSafeCSSValue(f *testing.F) {

	f.Add("#ff0000")
	f.Add("")
	f.Add("linear-gradient(90deg, #fff, #000)")
	f.Add("12.5%")
	f.Add("red; background: url(x)")
	f.Add("url(javascript:alert(1))")
	f.Add("URL(x)")
	f.Add("expression(alert(1))")
	f.Add("EXPRESSION(1)")
	f.Add("</style>")
	f.Add("a\\6c ert")
	f.Add("héllo")

	f.Fuzz(func(t *testing.T, value string) {

		result := string(safeCSSValue(value))

		if result == "" {
			return
		}

		// A non-empty result must be the input, verbatim and unmodified
		require.Equal(t, value, result)

		// ...built only from allowlisted characters...
		for _, character := range result {
			require.False(t, isUnsafeCSSValueRune(character), "unsafe rune %q survived in %q", character, result)
		}

		// ...with none of the characters that break out of a value, a property, or the style attribute
		require.False(t, strings.ContainsAny(result, `"';<>{}@\/*:!`))

		// ...and never the two dangerous functions, in any casing
		lowercased := strings.ToLower(result)
		require.False(t, strings.Contains(lowercased, "url"))
		require.False(t, strings.Contains(lowercased, "expression"))

		// The guard must be idempotent
		require.Equal(t, result, string(safeCSSValue(result)))
	})
}

// FuzzSafeURL asserts the invariant behind the `safeURL` template function: the result is
// used as a navigation target, so it must either be empty, or be the input unchanged and
// parse as one of exactly two safe shapes -- a same-site relative path (no scheme AND no
// host), or an absolute http(s) URL.
func FuzzSafeURL(f *testing.F) {

	f.Add("/stream/123")
	f.Add("")
	f.Add("https://example.com/path")
	f.Add("http://example.com")
	f.Add("javascript:alert(1)")
	f.Add("JavaScript:alert(1)")
	f.Add("\tjavascript:alert(1)")
	f.Add("data:text/html,<script>alert(1)</script>")
	f.Add("//evil.com/path")
	f.Add("urn:isbn:1234")
	f.Add("mailto:ben@pate.org")
	f.Add("../relative")
	f.Add("?query=1")
	f.Add("#fragment")

	f.Fuzz(func(t *testing.T, value string) {

		result := safeURL(value)

		if result == "" {
			return
		}

		// A non-empty result must be the input, verbatim and unmodified
		require.Equal(t, value, result)

		// ...and must re-parse into one of the two permitted shapes
		parsed, err := url.Parse(result)
		require.NoError(t, err)

		if parsed.Scheme == "" {
			require.Equal(t, "", parsed.Host, "scheme-relative URL survived: %q", result)
		} else {
			require.Contains(t, []string{"http", "https"}, parsed.Scheme, "dangerous scheme survived: %q", result)
		}

		// The guard must be idempotent
		require.Equal(t, result, safeURL(result))
	})
}

// FuzzHTMLFuncs_NeverPanic feeds arbitrary strings through the escaping and string-munging
// template functions to confirm that none of them panics, and that the ones which promise
// escaping actually neutralize the markup they are handed.
func FuzzHTMLFuncs_NeverPanic(f *testing.F) {

	f.Add("<script>alert(1)</script>", "script")
	f.Add("", "")
	f.Add("héllo wörld", "ö")
	f.Add("\x00\xff", "\xff")
	f.Add("https://example.com/a/b", "example")
	f.Add("a&b<c>d", "&")

	functions := All()

	f.Fuzz(func(t *testing.T, value string, search string) {

		// The "text" function escapes its input, so no live markup may survive
		text := string(functions["text"].(func(string) template.HTML)(value))
		require.False(t, strings.Contains(text, "<script"))

		// "highlight" escapes both arguments, so the only live markup is its own <b> wrapper
		highlighted := string(functions["highlight"].(func(string, string) template.HTML)(value, search))
		require.False(t, strings.Contains(strings.ReplaceAll(highlighted, `<b class="highlight">`, ""), "<script"))

		// The remaining string helpers must simply never panic
		functions["js"].(func(string) string)(value)
		functions["queryEscape"].(func(string) string)(value)
		functions["domainOnly"].(func(string) string)(value)
		functions["stripProtocol"].(func(string) string)(value)
		functions["textOnly"].(func(string) string)(value)
		functions["hasImage"].(func(string) bool)(value)
		functions["addQueryParams"].(func(string, string) string)(search, value)
	})
}
