package schema

import (
	"math"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"golang.org/x/exp/constraints"
)

// validate_Integer checks that the provided value meets the requirements of the Integer schema element, and updates the value if necessary.
func validate_Integer(element Integer, value any) (any, bool, error) {

	const location = "schema.validate_Integer"

	switch typedValue := value.(type) {

	case int:
		return validate_Integer_Generic(element, typedValue)
	case int8:
		return validate_Integer_Generic(element, typedValue)
	case int16:
		return validate_Integer_Generic(element, typedValue)
	case int32:
		return validate_Integer_Generic(element, typedValue)
	case int64:
		return validate_Integer_Generic(element, typedValue)
	}

	return nil, false, derp.Internal(location, "Value must be an integer", value)
}

// validate_Integer_Generic checks that the provided value meets the requirements of the schema element.
func validate_Integer_Generic[T constraints.Integer](element Integer, value T) (T, bool, error) {

	// RULE: Required value cannot be zero
	if element.Required && (value == 0) {
		return 0, false, derp.Validation("Value is required")
	}

	// RULE: Rewrite value if it is below the minimum
	if element.Minimum.IsPresent() && (value < T(element.Minimum.Int64())) {
		return T(element.Minimum.Int64()), true, nil
	}

	// RULE: Rewrite value if it is above the maximum
	if element.Maximum.IsPresent() && (value > T(element.Maximum.Int64())) {
		return T(element.Maximum.Int64()), true, nil
	}

	// RULE: Value must be a multiple of the specified value
	if element.MultipleOf.IsPresent() && notMultipleOfInteger(value, T(element.MultipleOf.Int64())) {
		return value, false, derp.Validation("Must be a multiple of " + convert.String(element.MultipleOf))
	}

	// RULE: Value must be one of the specified values
	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, value) {
		return value, false, derp.Validation("Must be one of the specified values")
	}

	// Return the value converted back to the target type
	return value, false, nil
}

// narrowIntegerBitSize narrows a 64-bit integer to the element's declared bit size,
// returning a validation error if the value falls outside the range of that width.
func narrowIntegerBitSize(bitSize int, value int64) (any, error) {

	// This guards against a raw conversion silently wrapping an out-of-range value
	// (for example int8(300) == 44) before the minimum/maximum rules are applied.

	switch bitSize {

	case 8:
		if (value < math.MinInt8) || (value > math.MaxInt8) {
			return nil, derp.Validation("Value must fit within an 8-bit integer", value)
		}
		return int8(value), nil

	case 16:
		if (value < math.MinInt16) || (value > math.MaxInt16) {
			return nil, derp.Validation("Value must fit within a 16-bit integer", value)
		}
		return int16(value), nil

	case 32:
		if (value < math.MinInt32) || (value > math.MaxInt32) {
			return nil, derp.Validation("Value must fit within a 32-bit integer", value)
		}
		return int32(value), nil

	case 64:
		return value, nil

	default:
		// A bare "int" follows the platform width; on 32-bit platforms it matches int32.
		if (strconv.IntSize == 32) && ((value < math.MinInt32) || (value > math.MaxInt32)) {
			return nil, derp.Validation("Value must fit within an integer", value)
		}
		return int(value), nil
	}
}
