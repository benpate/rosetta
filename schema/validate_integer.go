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
