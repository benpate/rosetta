package format

import (
	"regexp"
	"strings"

	"github.com/benpate/derp"
)

// uuidPattern matches a canonical, lowercase UUID: 32 hexadecimal digits in groups of 8-4-4-4-12.
// It is compiled once, at package scope, because the constructor below runs on every validation.
var uuidPattern = regexp.MustCompile(
	"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
)

// UUID validates a UUID, and stores it in its canonical lowercase form.  A URL whose last
// path segment is a UUID (such as a MusicBrainz page) is accepted, and only the UUID is kept.
func UUID(_ string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty UUIDs
		value = strings.TrimSpace(value)

		if value == "" {
			return "", nil
		}

		// Keep only the UUID from a pasted URL
		value = strings.ToLower(lastPathSegment(value))

		if uuidPattern.MatchString(value) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid UUID", value)
	}
}

// lastPathSegment returns the last non-empty path segment of a URL, without its query or
// fragment.  A value with no slash is returned unchanged.
func lastPathSegment(value string) string {

	// Remove any fragment, then any query string
	value, _, _ = strings.Cut(value, "#")
	value, _, _ = strings.Cut(value, "?")

	// Ignore trailing slashes, then take what follows the last one
	value = strings.TrimRight(value, "/")

	if index := strings.LastIndex(value, "/"); index >= 0 {
		return value[index+1:]
	}

	return value
}
