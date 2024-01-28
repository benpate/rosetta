package html

import "github.com/microcosm-cc/bluemonday"

var minimalHTMLPolicy *bluemonday.Policy

func init() {
	minimalHTMLPolicy = bluemonday.UGCPolicy()
	minimalHTMLPolicy.AllowElements("br", "p", "b", "i", "u", "img", "div", "pre", "code", "ol", "ul", "li")
}

func Minimal(text string) string {
	return minimalHTMLPolicy.Sanitize(text)
}
