package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// Color validates an email address using Go's built-in system email parser.
func Color(arg string) StringFormat {

	color := regexp.MustCompile("(?i)^#[0-9a-f]{6}$")

	return func(value string) (string, error) {

		// Allow empty addresses
		if value == "" {
			return "", nil
		}

		// Colors must match the regular expression.
		if !color.Match([]byte(value)) {
			return "", derp.BadRequestError("schema.format.Color", "Value is not a valid color", value)
		}

		return value, nil
	}
}
