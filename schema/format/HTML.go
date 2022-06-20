package format

import (
	"github.com/microcosm-cc/bluemonday"
)

var htmlPolicy *bluemonday.Policy

func init() {
	htmlPolicy = bluemonday.UGCPolicy()
	htmlPolicy.AddTargetBlankToFullyQualifiedLinks(true)
	htmlPolicy.RequireNoFollowOnFullyQualifiedLinks(true)
	htmlPolicy.RequireParseableURLs(true)
}

// HTML allows basic HTML tags, but strips iframes, object, embed, style, script tags.
func HTML(arg string) StringFormat {
	return func(value string) (string, error) {
		return htmlPolicy.Sanitize(value), nil
	}
}
