package format

import "github.com/benpate/rosetta/html"

// NoHTML strips all HTML tags from a string.
func NoHTML(arg string) StringFormat {
	return func(value string) (string, error) {
		value = html.RemoveTags(value)
		value = html.CollapseWhitespace(value)
		return value, nil
	}
}
