package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// validate_Boolean checks that the provided value meets the requirements of the Boolean schema element.
func validate_Boolean[T any](element Boolean, value T) (T, bool, error) {

	// Coerce to a bool so the schema rules can be checked against it
	boolValue, isBoolean := convert.BoolOk(value, false)

	// RULE: Value must be a boolean
	if !isBoolean {
		return value, false, derp.Validation("Value must be a boolean")
	}

	// RULE: Required value cannot be false
	if element.Required && (!boolValue) {
		return value, false, derp.Validation("Value is required")
	}

	// Return the original value unchanged; only its boolean form was needed to check the rules
	return value, false, nil
}
