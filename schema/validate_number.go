package schema

import (
	"math"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"golang.org/x/exp/constraints"
)

// validate_Number checks that the provided value meets the requirements of the Number schema element, and updates the value if necessary.
func validate_Number(element Number, value any) (float64, bool, error) {

	// Convert the value to a Float64 (we don't do Float32)
	value64, ok := convert.FloatOk(value, 0)

	if !ok {
		return value64, false, derp.Validation("Value must be a number")
	}

	// RULE: Required value cannot be zero
	if element.Required && (value == 0) {
		return value64, false, derp.Validation("Value is required")
	}

	// RULE: Rewrite value if it is below the minimum
	if element.Minimum.IsPresent() && (value64 < element.Minimum.Float()) {
		return element.Minimum.Float(), true, nil
	}

	// RULE: Rewrite value if it is above the maximum
	if element.Maximum.IsPresent() && (value64 > element.Maximum.Float()) {
		return element.Maximum.Float(), true, nil
	}

	// RULE: Value must be a multiple of the specified value
	if element.MultipleOf.IsPresent() && notMultipleOfFloat(value64, element.MultipleOf.Float()) {
		return value64, false, derp.Validation("Must be a multiple of " + convert.String(element.MultipleOf))
	}

	// RULE: Value must be one of the specified values
	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, value64) {
		return value64, false, derp.Validation("Must be one of the specified values")
	}

	// Return the value converted back to the target type
	return value64, false, nil
}

// isMultipleOfFloat reports whether value is an exact multiple of multipleOf.
// A multipleOf of zero is treated as "no constraint".
func isMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {
	if multipleOf == 0 {
		return true
	}
	return math.Mod(float64(value), float64(multipleOf)) == 0
}

// notMultipleOfFloat returns TRUE when the value is not an exact multiple of multipleOf.
func notMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {
	return !isMultipleOfFloat(value, multipleOf)
}
