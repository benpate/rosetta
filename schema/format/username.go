package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// Username validates a simple token string suitable for use as URL identifiers
func Username(_ string) StringFormat {

	token := regexp.MustCompile(`(?i)^[A-z0-9_]*$`)

	return func(value string) (string, error) {

		if token.MatchString(value) {
			return value, nil
		}

		return "", derp.NewInternalError("schema.format.Username", "Value is not a valid Username", value)
	}
}
