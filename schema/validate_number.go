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
		return value64, false, derp.Validation("Value must be a valid floating point number")
	}

	// RULE: Number cannot be NaN
	if math.IsNaN(value64) {
		return value64, false, derp.Validation("Floating point value must not be NaN")
	}

	// RULE: Number cannot be "positive Infinity"
	if math.IsInf(value64, 1) {

		// If possible, clamp the value to the maximum
		if element.Maximum.IsPresent() {
			return element.Maximum.Float(), true, nil
		}
		return value64, false, derp.Validation("Floating point value must not be positive infinity")
	}

	// RULE: Number cannot be "negative Infinity"
	if math.IsInf(value64, -1) {

		// If possible, clamp the value to the minimum
		if element.Minimum.IsPresent() {
			return element.Minimum.Float(), true, nil
		}
		return value64, false, derp.Validation("Floating point value must not be negative infinity")
	}

	// RULE: Required value cannot be zero
	if element.Required && (value64 == 0) {
		return value64, false, derp.Validation("Floating point value is required")
	}

	// RULE: Value must be a multiple of the specified value
	if element.MultipleOf.IsPresent() && notMultipleOfFloat(value64, element.MultipleOf.Float()) {
		return value64, false, derp.Validation("Floating point value must be a multiple of " + convert.String(element.MultipleOf))
	}

	// RULE: Value must be one of the specified values
	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, value64) {
		return value64, false, derp.Validation("Floating point value must be one of the specified values")
	}

	// RULE: Rewrite value if it is below the minimum
	if element.Minimum.IsPresent() && (value64 < element.Minimum.Float()) {
		return element.Minimum.Float(), true, nil
	}

	// RULE: Rewrite value if it is above the maximum
	if element.Maximum.IsPresent() && (value64 > element.Maximum.Float()) {
		return element.Maximum.Float(), true, nil
	}

	// The value satisfies every rule; return it unchanged
	return value64, false, nil
}

// isMultipleOfFloat reports whether value is a multiple of multipleOf.
// A multipleOf of zero is treated as "no constraint".
func isMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {

	// A zero divisor is treated as "no constraint" and must not panic.
	if multipleOf == 0 {
		return true
	}

	// Round the quotient to the nearest integer and confirm it reproduces the value. A direct
	// math.Mod(value, multipleOf) == 0 check fails for exact decimal multiples (math.Mod(0.3, 0.1)
	// is ~0.0999... not 0) because neither operand is representable exactly in float64, so we
	// compare with a tolerance scaled to the magnitude of the value instead.
	v := float64(value)
	step := float64(multipleOf)
	rounded := math.Round(v/step) * step
	return math.Abs(v-rounded) <= 1e-9*math.Max(1, math.Abs(v))
}

// notMultipleOfFloat returns TRUE when the value is not an exact multiple of multipleOf.
func notMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {
	return !isMultipleOfFloat(value, multipleOf)
}
