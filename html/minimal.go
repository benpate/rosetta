package html

import "github.com/microcosm-cc/bluemonday"

var minimalHTMLPolicy *bluemonday.Policy

func init() {
	minimalHTMLPolicy = bluemonday.UGCPolicy()
	minimalHTMLPolicy.AllowElements("br", "p", "b", "i", "u", "img", "div", "pre", "code", "ol", "ul", "li")
}

// Minimal sanitizes the provided HTML, allowing only a minimal set of formatting elements.
func Minimal(text string) string {
	return minimalHTMLPolicy.Sanitize(text)
}
