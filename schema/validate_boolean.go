package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// validate_Boolean checks that the provided value meets the requirements of the schema element, and updates the value if necessary.
func validate_Boolean[T any](element Boolean, value T) (T, bool, error) {

	// Make a boolean version of this value to compare schema rules
	boolValue, isBoolean := convert.BoolOk(value, false)

	// RULE: Value must be a boolean
	if !isBoolean {
		return value, false, derp.Validation("Must be a boolean")
	}

	// RULE: Required value cannot be false
	if element.Required && (!boolValue) {
		return value, false, derp.Validation("Value is required")
	}

	// Return the value converted back to the target type
	return value, false, nil
}
