package format

import (
	"regexp"

	"github.com/benpate/derp"
)

// colorPattern matches a six-digit hexadecimal color, with or without case.
// It is compiled once, at package scope: the StringFormat constructors below run on
// EVERY validation, so compiling inside one would recompile the pattern each time.
var colorPattern = regexp.MustCompile("(?i)^#[0-9a-f]{6}$")

// Color validates an email address using Go's built-in system email parser.
func Color(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty addresses
		if value == "" {
			return "", nil
		}

		// Colors must match the regular expression.
		if !colorPattern.Match([]byte(value)) {
			return "", derp.BadRequest("schema.format.Color", "Value is not a valid color", value)
		}

		return value, nil
	}
}
