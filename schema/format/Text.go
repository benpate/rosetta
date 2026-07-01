package format

import "github.com/benpate/rosetta/html"

// Text strips all HTML tags from a string without collapsing whitespace.
func Text(arg string) StringFormat {
	return func(value string) (string, error) {
		value = html.RemoveTags(value)
		return value, nil
	}
}
