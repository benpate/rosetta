package html

import (
	"regexp"
)

var firstParagraph *regexp.Regexp

func init() {
	firstParagraph = regexp.MustCompile("^(.*?)<(br|/p|/div)>")
}

// Summary returns the first few sentences of content from an HTML document
func Summary(html string) string {

	// Try to get just the first paragraph of content.
	result := firstParagraph.FindString(html)

	// If there is no "paragraph" marker (e.g. this is plain text) then just use the original
	if result == "" {
		result = html
	}

	// Remove all tags and whitespace
	result = CollapseWhitespace(RemoveTags(result))

	// Truncate to 200 characters (plus ellipsis)
	if len(result) > 200 {
		result = result[:200] + "..."
	}

	return result
}
