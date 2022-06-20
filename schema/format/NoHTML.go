package format

import (
	"github.com/benpate/htmlconv"
)

// NoHTML strips all HTML tags from a string.
func NoHTML(arg string) StringFormat {
	return func(value string) (string, error) {
		value = htmlconv.RemoveTags(value)
		value = htmlconv.CollapseWhitespace(value)
		return value, nil
	}
}
