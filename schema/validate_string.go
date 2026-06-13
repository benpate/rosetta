package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
)

func validate_String(element String, value string) (string, bool, error) {

	changed := false

	// Verify required fields (after format functions are applied)
	if element.Required && (value == "") {
		return value, false, derp.Validation("Value is required")
	}

	// Validate minimum value
	if (element.MinValue != "") && (value < element.MinValue) {
		return value, false, derp.Validation("Minimum value is " + element.MinValue)
	}

	// Validate maximum value
	if (element.MaxValue != "") && (value > element.MaxValue) {
		return value, false, derp.Validation("Maximum value is " + element.MaxValue)
	}

	// Validate minimum length
	if (element.MinLength > 0) && (len(value) < element.MinLength) {
		return value, false, derp.Validation("Minimum length is " + convert.String(element.MinLength))
	}

	// Validate maximum length (ACTUALLY TRUNCATES THE STRING TO THE MAXIMUM LENGTH)
	if (element.MaxLength > 0) && (len(value) > element.MaxLength) {
		value = value[:element.MaxLength]
		changed = true
	}

	// Validate enumerated values
	if len(element.Enum) > 0 {
		if (value != "") && (compare.NotContains(element.Enum, value)) {
			return value, false, derp.Validation("Value must match one of...", value, element.Enum)
		}
	}

	// Validate against all formatting functions, tracking changes
	return validate_String_Formats(element, value, changed)
}

func validate_String_Formats(element String, value string, changed bool) (string, bool, error) {

	const location = "schema.validate_String_Formats"

	// Validate against all formatting functions
	for _, formatFunc := range element.formatFunctions() {

		// Try formatting the value with this function
		newValue, err := formatFunc(value)

		if err != nil {
			return value, false, derp.Wrap(err, location, "Applying format function", value)
		}

		// If the value has been changed, then flag it so.
		if newValue != value {
			changed = true
			value = newValue
		}
	}

	return value, changed, nil
}
