package format

import (
	"errors"
	"strings"
)

// In returns a StringFormat that requires the value to be one of the comma-separated options in arg.
func In(arg string) StringFormat {

	options := strings.Split(arg, ",")

	return func(value string) (string, error) {

		for _, option := range options {
			if value == option {
				return value, nil
			}
		}

		return "", errors.New(value + " is not an allowed value")
	}
}

// NotIn returns a StringFormat that rejects the value if it is one of the comma-separated options in arg.
func NotIn(arg string) StringFormat {

	options := strings.Split(arg, ",")

	return func(value string) (string, error) {

		for _, option := range options {
			if value == option {
				return "", errors.New(value + " is not an allowed value")
			}
		}

		return arg, nil
	}
}
