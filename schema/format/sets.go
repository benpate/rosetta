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

		// RULE: Return the VALUE, never `arg`. A format returns what the caller supplied, and
		// schema.Set writes that result back into the object -- so returning the option list
		// here overwrote every accepted field with its own configuration.
		return value, nil
	}
}
