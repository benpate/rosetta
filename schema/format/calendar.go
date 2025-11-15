package format

import (
	"regexp"

	"github.com/benpate/derp"
)

var iso8601 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(:\d{2}(Z[+-]\d{2}:\d{2})?)?$`)

func ISO8601(arg string) StringFormat {
	return func(value string) (string, error) {
		if iso8601.Match([]byte(value)) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid date string", value)
	}
}

func Date(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}

func DateTime(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}

func Time(arg string) StringFormat {

	return func(value string) (string, error) {
		// TODO: HIGH: Implement or delete this
		return value, nil
	}
}
