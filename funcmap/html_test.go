package funcmap

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTMLFuncs_URLsAndStrings(t *testing.T) {

	f := All()

	addQueryParams := f["addQueryParams"].(func(string, string) string)
	require.Equal(t, "http://x.com?a=1", addQueryParams("a=1", "http://x.com"))
	require.Equal(t, "http://x.com?b=2&a=1", addQueryParams("a=1", "http://x.com?b=2"))

	domainOnly := f["domainOnly"].(func(string) string)
	require.Equal(t, "example.com", domainOnly("https://example.com/path"))
	require.Equal(t, "example.com", domainOnly("http://example.com"))

	stripProtocol := f["stripProtocol"].(func(string) string)
	require.Equal(t, "example.com", stripProtocol("https://example.com"))
	require.Equal(t, "example.com", stripProtocol("http://example.com"))
	require.Equal(t, "example.com", stripProtocol("example.com"))

	safeURL := f["safeURL"].(func(string) string)
	// Same-site relative paths and absolute http(s) URLs pass through unchanged.
	require.Equal(t, "/stream/123", safeURL("/stream/123"))
	require.Equal(t, "/stream/123?q=1#top", safeURL("/stream/123?q=1#top"))
	require.Equal(t, "relative/path", safeURL("relative/path"))
	require.Equal(t, "", safeURL(""))
	require.Equal(t, "https://remote.example/@author", safeURL("https://remote.example/@author"))
	require.Equal(t, "http://remote.example/news", safeURL("http://remote.example/news"))
	// Dangerous or non-navigational schemes are rejected.
	require.Equal(t, "", safeURL("javascript:alert(1)"))
	require.Equal(t, "", safeURL("JavaScript:alert(1)")) // scheme is case-insensitive
	require.Equal(t, "", safeURL("data:text/html,<script>alert(1)</script>"))
	require.Equal(t, "", safeURL("urn:loaded")) // AS object IDs can be urn: — not navigable
	require.Equal(t, "", safeURL("vbscript:msgbox(1)"))
	require.Equal(t, "", safeURL("file:///etc/passwd"))
	// Scheme-relative URLs have an empty scheme but a host, so a browser treats
	// them as off-site absolute URLs — reject them.
	require.Equal(t, "", safeURL("//evil.example/phish"))
	// Leading-whitespace tricks fail to parse and must fail closed.
	require.Equal(t, "", safeURL("\tjavascript:alert(1)"))

	require.Equal(t, "a+b", f["queryEscape"].(func(string) string)("a b"))
	// JSEscapeString escapes angle brackets to their unicode escapes
	require.Equal(t, `\u003Cb\u003E`, f["js"].(func(string) string)("<b>"))
}

func TestHTMLFuncs_Markup(t *testing.T) {

	f := All()

	attr := f["attr"].(func(string) template.HTMLAttr)
	// Legitimate attribute values (tokens, data values, ARIA states) pass through.
	require.Equal(t, template.HTMLAttr("data-x"), attr("data-x"))
	require.Equal(t, template.HTMLAttr("btn-primary"), attr("btn-primary"))
	require.Equal(t, template.HTMLAttr("true"), attr("true"))
	// Anything that could break out of the attribute or inject a new one is rejected.
	require.Equal(t, template.HTMLAttr(""), attr(`x" onmouseover="alert(1)`)) // quote closes the attribute
	require.Equal(t, template.HTMLAttr(""), attr("x onmouseover=alert(1)"))   // whitespace starts a new attribute
	require.Equal(t, template.HTMLAttr(""), attr("x=y"))                      // `=` starts a value
	require.Equal(t, template.HTMLAttr(""), attr("a<b"))                      // `<` opens a tag
	require.Equal(t, template.HTMLAttr(""), attr("a&amp;b"))                  // `&` is entity-ambiguous
	require.Equal(t, template.CSS("color:red"), f["css"].(func(string) template.CSS)("color:red"))

	cssValue := f["cssValue"].(func(string) template.CSS)
	// Legitimate computed values pass through unchanged.
	require.Equal(t, template.CSS("#ff0000"), cssValue("#ff0000"))
	require.Equal(t, template.CSS("rgba(255, 0, 0, 0.50)"), cssValue("rgba(255, 0, 0, 0.50)"))
	require.Equal(t, template.CSS("linear-gradient(120deg, rgba(1, 2, 3, 1.00), rgba(4, 5, 6, 0.50))"), cssValue("linear-gradient(120deg, rgba(1, 2, 3, 1.00), rgba(4, 5, 6, 0.50))"))
	// Anything with a context-breakout character is rejected wholesale.
	require.Equal(t, template.CSS(""), cssValue("red; background:url(javascript:alert(1))")) // `;` and `:`
	require.Equal(t, template.CSS(""), cssValue(`red" onmouseover="alert(1)`))               // `"` closes the attribute
	require.Equal(t, template.CSS(""), cssValue("red}body{display:none"))                    // `}` breaks out of a rule
	require.Equal(t, template.CSS(""), cssValue("@import 'evil.css'"))                       // `@` and `'`
	require.Equal(t, template.CSS(""), cssValue(`\3c script\3e`))                            // `\` CSS escape
	// Dangerous functions whose spelling is all-allowlisted are rejected by name.
	require.Equal(t, template.CSS(""), cssValue("expression(alert(1))"))
	require.Equal(t, template.CSS(""), cssValue("EXPRESSION(alert(1))")) // case-insensitive
	require.Equal(t, template.CSS(""), cssValue("url(x)"))               // url() would still parse under the charset allowlist
	require.Equal(t, template.HTML("<p>hi</p>"), f["html"].(func(string) template.HTML)("<p>hi</p>"))

	highlight := f["highlight"].(func(string, string) template.HTML)
	require.Equal(t, template.HTML(`a <b class="highlight">b</b> c`), highlight("a b c", "b"))
	require.Equal(t, template.HTML("a b c"), highlight("a b c", "")) // empty search returns text unchanged

	// Untrusted text and search terms are HTML-escaped; only the <b> wrapper is live markup.
	require.Equal(t, template.HTML("&lt;img src=x onerror=alert(1)&gt;"), highlight("<img src=x onerror=alert(1)>", ""))
	require.Equal(t, template.HTML(`&lt;script&gt;<b class="highlight">x</b>&lt;/script&gt;`), highlight("<script>x</script>", "x"))
	require.Equal(t, template.HTML(`a <b class="highlight">&lt;b&gt;</b> c`), highlight("a <b> c", "<b>")) // markup in the search term is escaped too

	hasImage := f["hasImage"].(func(string) bool)
	require.True(t, hasImage(`<img src="x">`))
	require.True(t, hasImage(`<picture>`))
	require.False(t, hasImage("plain text"))
}

func TestHTMLFuncs_JSON(t *testing.T) {

	f := All()

	jsonFn := f["json"].(func(any) string)
	require.Equal(t, `{"a":1}`, jsonFn(map[string]int{"a": 1}))

	jsonIndent := f["jsonIndent"].(func(any) string)
	require.Contains(t, jsonIndent(map[string]int{"a": 1}), "\n")
}

func TestHTMLFuncs_Markdown(t *testing.T) {

	markdown := All()["markdown"].(func(any) template.HTML)
	result := string(markdown("# Heading"))
	require.Contains(t, result, "<h1")
	require.Contains(t, result, "Heading")
}

func TestHTMLFuncs_Text(t *testing.T) {

	f := All()

	require.NotNil(t, f["htmlMinimal"].(func(string) template.HTML))
	require.NotNil(t, f["summary"])
	require.NotNil(t, f["textOnly"])

	text := f["text"].(func(string) template.HTML)
	require.NotEmpty(t, string(text("hello")))
}
