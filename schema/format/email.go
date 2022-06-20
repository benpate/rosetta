package format

import (
	"net/mail"

	"github.com/benpate/derp"
)

// Email validates an email address using Go's built-in system email parser.
func Email(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty addresses
		if value == "" {
			return "", nil
		}

		if _, err := mail.ParseAddress(value); err != nil {
			return "", derp.Wrap(err, "schema.format.Email", "Invalid email Address", value)
		}

		return value, nil
	}
}
