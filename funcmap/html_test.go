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

	require.Equal(t, "a+b", f["queryEscape"].(func(string) string)("a b"))
	// JSEscapeString escapes angle brackets to their unicode escapes
	require.Equal(t, `\u003Cb\u003E`, f["js"].(func(string) string)("<b>"))
}

func TestHTMLFuncs_Markup(t *testing.T) {

	f := All()

	require.Equal(t, template.HTMLAttr("data-x"), f["attr"].(func(string) template.HTMLAttr)("data-x"))
	require.Equal(t, template.CSS("color:red"), f["css"].(func(string) template.CSS)("color:red"))
	require.Equal(t, template.HTML("<p>hi</p>"), f["html"].(func(string) template.HTML)("<p>hi</p>"))

	highlight := f["highlight"].(func(string, string) template.HTML)
	require.Equal(t, template.HTML(`a <b class="highlight">b</b> c`), highlight("a b c", "b"))
	require.Equal(t, template.HTML("a b c"), highlight("a b c", "")) // empty search returns text unchanged

	hasImage := f["hasImage"].(func(string) bool)
	require.True(t, hasImage(`<img src="x">`))
	require.True(t, hasImage(`<picture>`))
	require.False(t, hasImage("plain text"))

	removeLinks := f["removeLinks"].(func(string) template.HTML)
	require.Equal(t, template.HTML(`<span href="x">y</span>`), removeLinks(`<a href="x">y</a>`))
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
