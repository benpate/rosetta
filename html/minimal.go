package html

import "github.com/microcosm-cc/bluemonday"

// minimalHTMLPolicy sanitizes HTML down to a small set of formatting elements.
var minimalHTMLPolicy *bluemonday.Policy

// init builds the minimal sanitizer policy.
func init() {
	minimalHTMLPolicy = bluemonday.UGCPolicy()
	minimalHTMLPolicy.AllowElements("br", "p", "b", "i", "u", "img", "div", "pre", "code", "ol", "ul", "li")
}

// Minimal sanitizes the provided HTML, allowing only a minimal set of formatting elements.
func Minimal(text string) string {
	return minimalHTMLPolicy.Sanitize(text)
}
