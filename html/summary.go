package html

// Summary returns the plain text of an HTML document, cut to 200 runes plus "..." when longer
func Summary(html string) string {

	// Remove all tags and collapse whitespace
	result := CollapseWhitespace(RemoveTags(html))

	// Truncate to 200 runes (plus ellipsis)
	if runes := []rune(result); len(runes) > 200 {
		result = string(runes[:200]) + "..."
	}

	return result
}
