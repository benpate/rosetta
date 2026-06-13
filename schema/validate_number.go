package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
)

func validate_Number(element Number, value float64) (float64, bool, error) {

	// RULE: Required value cannot be zero
	if element.Required && (value == 0) {
		return value, false, derp.Validation(" float field is required")
	}

	// RULE: Rewrite value if it is below the minimum
	if element.Minimum.IsPresent() && (value < element.Minimum.Float()) {
		return element.Minimum.Float(), true, nil
	}

	// RULE: Rewrite value if it is above the maximum
	if element.Maximum.IsPresent() && (value > element.Maximum.Float()) {
		return element.Maximum.Float(), true, nil
	}

	// RULE: Value must be a multiple of the specified value
	if element.MultipleOf.IsPresent() && !isMultipleOf(value, element.MultipleOf.Float()) {
		return value, false, derp.Validation(" float must be a multiple of " + convert.String(element.MultipleOf))
	}

	// RULE: Value must be one of the specified values
	if (len(element.Enum) > 0) && !compare.Contains(element.Enum, value) {
		return value, false, derp.Validation(" float must contain one of the specified values")
	}

	// Return the value converted back to the target type
	return value, false, nil
}
