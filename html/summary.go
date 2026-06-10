package html

// var firstParagraph *regexp.Regexp

// func init() {
// 	firstParagraph = regexp.MustCompile("^(.*?)<(br|/p|/div)>")
// }

// Summary returns the first few sentences of content from an HTML document
func Summary(html string) string {

	// Remove all tags and whitespace
	result := CollapseWhitespace(RemoveTags(html))

	// Truncate to 200 runes (plus ellipsis)
	if runes := []rune(result); len(runes) > 200 {
		result = string(runes[:200]) + "..."
	}

	return result
}
