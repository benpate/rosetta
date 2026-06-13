package format

import (
	"regexp"

	"github.com/benpate/derp"
)

var iso8601 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(:\d{2}(Z[+-]\d{2}:\d{2})?)?$`)

// ISO8601 returns a StringFormat that validates a value as an ISO-8601 date string.
func ISO8601(arg string) StringFormat {
	return func(value string) (string, error) {
		if iso8601.Match([]byte(value)) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid date string", value)
	}
}

// Date returns a StringFormat intended to validate date strings; validation is not yet implemented.
func Date(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}

// DateTime returns a StringFormat intended to validate date-time strings; validation is not yet implemented.
func DateTime(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}

// Time returns a StringFormat intended to validate time strings; validation is not yet implemented.
func Time(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}
