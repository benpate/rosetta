package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"golang.org/x/exp/constraints"
)

// validate_Integer checks that the provided value meets the requirements of the Integer schema element, and updates the value if necessary.
func validate_Integer(element Integer, value any) (any, bool, error) {

	const location = "schema.validate_Integer"

	// Convert the value into an integer of the specified bit size.
	coercedValue, lossless, inBounds := convert.IntBitsizeOk(value, 0, element.BitSize)

	if !lossless {
		return value, false, derp.Validation("Value must be an integer")
	}

	if !inBounds {
		return value, false, derp.Validation("Value must fit within the specified bit size")
	}

	switch typedValue := coercedValue.(type) {

	case int:
		return validate_Integer_Generic(element, typedValue, !lossless)
	case int8:
		return validate_Integer_Generic(element, typedValue, !lossless)
	case int16:
		return validate_Integer_Generic(element, typedValue, !lossless)
	case int32:
		return validate_Integer_Generic(element, typedValue, !lossless)
	case int64:
		return validate_Integer_Generic(element, typedValue, !lossless)
	}

	// This should never happen
	return nil, false, derp.Internal(location, "Value must be an integer. This should never happen.", value)
}

// validate_Integer_Generic checks that the provided value meets the requirements of the schema element.
func validate_Integer_Generic[T constraints.Integer](element Integer, value T, changed bool) (T, bool, error) {

	const location = "schema.validate_Integer_Generic"

	// Get this value as an int64 so we can compare it correctly.
	value64 := int64(value)

	// RULE: Required value cannot be zero
	if element.Required && (value == 0) {
		return 0, false, derp.Validation("Value is required")
	}

	// RULE: Value must be a multiple of the specified value
	if element.MultipleOf.IsPresent() && notMultipleOfInteger(value64, element.MultipleOf.Int64()) {
		return T(value), false, derp.Validation("Must be a multiple of " + convert.String(element.MultipleOf))
	}

	// RULE: Value must be one of the specified values
	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, value) {
		return T(value), false, derp.Validation("Must be one of the specified values")
	}

	// RULE: Rewrite value if it is below the minimum
	if element.Minimum.IsPresent() {

		minValue := element.Minimum.Int64()

		// Verify the minimum value fits within the target bitsize
		if _, _, ok := convert.IntBitsizeOk(minValue, 0, element.BitSize); !ok {
			return 0, false, derp.Internal(location, "Minimum value is out of bounds for specified bit size", minValue, element.BitSize)
		}

		// Clamp the value to the minimum
		if value64 < minValue {
			value64 = minValue
			changed = true
		}
	}

	// RULE: Rewrite the value if it is above the maximum
	if element.Maximum.IsPresent() {

		maxValue := element.Maximum.Int64()

		// Verify the maximum value fits within the target bitsize
		if _, _, ok := convert.IntBitsizeOk(maxValue, 0, element.BitSize); !ok {
			return 0, false, derp.Internal(location, "Maximum value is out of bounds for specified bit size", maxValue, element.BitSize)
		}

		// Clamp the value to the maximum
		if value64 > maxValue {
			value64 = maxValue
			changed = true
		}
	}

	// Return the value converted back to the target type
	return T(value64), changed, nil
}

// isMultipleOfInteger reports whether value is an exact integer multiple of
// multipleOf.
func isMultipleOfInteger[T constraints.Integer](value, multipleOf T) bool {

	// A multipleOf of zero is treated as "no constraint" (and also avoids a divide-by-zero panic).
	if multipleOf == 0 {
		return true
	}

	// Using integer modulo so that large values are not corrupted by a detour through float64.
	return value%multipleOf == 0
}

// notMultipleOfInteger returns TRUE when the value is not an exact integer multiple of multipleOf.
func notMultipleOfInteger[T constraints.Integer](value, multipleOf T) bool {
	return !isMultipleOfInteger(value, multipleOf)
}
