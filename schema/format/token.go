package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// Token validates a simple token string suitable for use as URL identifiers
func Token(_ string) StringFormat {

	token := regexp.MustCompile("(?i)^[A-Z0-9-_]+$")

	return func(value string) (string, error) {

		// Allow empty tokens
		if value == "" {
			return value, nil
		}

		// Non-empty IDs must be 24 hexadecimal characters
		if token.MatchString(value) {
			return value, nil
		}

		return "", derp.NewInternalError("schema.format.Token", "Value is not a valid Token", value)
	}
}
