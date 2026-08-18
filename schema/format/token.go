package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// tokenPattern matches letters, numbers, dashes, and underscores, case-insensitively.
// It is compiled once, at package scope: this pattern is by far the most expensive in
// the package to compile, and the StringFormat constructors below run on EVERY validation.
var tokenPattern = regexp.MustCompile(`(?i)^[\p{L}\p{N}-_]+$`)

// Token validates a simple token string suitable for use as URL identifiers
func Token(_ string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty tokens
		if value == "" {
			return value, nil
		}

		// Non-empty IDs must look like a token (characters, numbers, dashes, and underscores)
		if tokenPattern.MatchString(value) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid Token", value)
	}
}
