package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// usernamePattern matches letters, numbers, and underscores, case-insensitively.
// It is compiled once, at package scope, because the constructor below runs on every validation.
var usernamePattern = regexp.MustCompile(`(?i)^[A-Z0-9_]*$`)

// Username validates a simple token string suitable for use as URL identifiers
func Username(_ string) StringFormat {

	return func(value string) (string, error) {

		if usernamePattern.MatchString(value) {
			return value, nil
		}

		return "", derp.Internal("schema.format.Username", "Usernames can only contain letters, numbers, and underscores.", value)
	}
}
