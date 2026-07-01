package format

import "github.com/benpate/rosetta/html"

// NoHTML strips all HTML tags from a string and collapses whitespace into a single space character.
func NoHTML(arg string) StringFormat {
	return func(value string) (string, error) {
		value = html.RemoveTags(value)
		value = html.CollapseWhitespace(value)
		return value, nil
	}
}
