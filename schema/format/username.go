package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// usernamePattern matches ASCII letters, numbers, and underscores
var usernamePattern = regexp.MustCompile(
	// List both cases: (?i) folds the Kelvin sign (U+212A) and long s (U+017F) into K and s
	`^[A-Za-z0-9_]*$`,
)

// Username validates a simple token string suitable for use as URL identifiers
func Username(_ string) StringFormat {

	return func(value string) (string, error) {

		if usernamePattern.MatchString(value) {
			return value, nil
		}

		return "", derp.Internal("schema.format.Username", "Usernames can only contain letters, numbers, and underscores.", value)
	}
}
