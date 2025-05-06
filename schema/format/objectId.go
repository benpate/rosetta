package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// ObjectID validates a mongodb-style identifier (24 hexadecimal characters)
func ObjectID(arg string) StringFormat {

	objectID := regexp.MustCompile("(?i)^[A-F0-9]{24}$")

	return func(value string) (string, error) {

		// Allow empty objectIds
		if value == "" {
			return value, nil
		}

		// Non-empty IDs must be 24 hexadecimal characters
		if objectID.MatchString(value) {
			return value, nil
		}

		return "", derp.InternalError("schema.format.ObjectID", "Value is not a valid ObjectID", value)
	}
}
