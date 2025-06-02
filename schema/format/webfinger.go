package format

import (
	"net/mail"
	"strings"

	"github.com/benpate/derp"
)

// WebFinger validates an email address using Go's built-in system email parser.
func WebFinger(arg string) StringFormat {

	const location = "schema.format.WebFinger"

	return func(value string) (string, error) {

		// Allow empty addresses
		if value == "" {
			return "", nil
		}

		if !strings.HasPrefix(value, "@") {
			return "", derp.BadRequestError(location, "Invalid WebFinger Handle", value, "WebFinger handles must start with '@'")
		}

		value = strings.TrimPrefix(value, "@")

		if _, err := mail.ParseAddress(value); err != nil {
			return "", derp.Wrap(err, location, "Invalid email Address", value)
		}

		return value, nil
	}
}
