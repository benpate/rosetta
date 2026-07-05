package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// Token validates a simple token string suitable for use as URL identifiers
func Token(_ string) StringFormat {

	// A token is a string that contains only letters, numbers, dashes, and underscores. It is case-insensitive.
	token := regexp.MustCompile(`(?i)^[\p{L}\p{N}-_]+$`)

	return func(value string) (string, error) {

		// Allow empty tokens
		if value == "" {
			return value, nil
		}

		// Non-empty IDs must look like a token (characters, numbers, dashes, and underscores)
		if token.MatchString(value) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid Token", value)
	}
}
