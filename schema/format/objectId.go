package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// objectIDPattern matches a 24-character hexadecimal identifier, with or without case.
// It is compiled once, at package scope, because the constructor below runs on every validation.
var objectIDPattern = regexp.MustCompile("(?i)^[A-F0-9]{24}$")

// ObjectID validates a mongodb-style identifier (24 hexadecimal characters)
func ObjectID(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty objectIds
		if value == "" {
			return value, nil
		}

		// Non-empty IDs must be 24 hexadecimal characters
		if objectIDPattern.MatchString(value) {
			return value, nil
		}

		return "", derp.Internal("schema.format.ObjectID", "Value is not a valid ObjectID", value)
	}
}
