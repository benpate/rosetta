package html

import (
	"regexp"
	"strings"
)

// whitespace matches a run of one or more whitespace characters.
var whitespace *regexp.Regexp

// init compiles the whitespace pattern.
func init() {
	whitespace = regexp.MustCompile(`\s+`)
}

// CollapseWhitespace converts each run of ASCII whitespace into a single SPACE character, and trims both ends
func CollapseWhitespace(text string) string {
	result := whitespace.ReplaceAllString(text, " ")

	result = strings.TrimPrefix(result, " ")
	result = strings.TrimSuffix(result, " ")
	return result
}
