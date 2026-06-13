package schema

import (
	"unicode/utf8"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
)

// validate_String checks that the provided value meets the requirements of the String schema element, and updates the value if necessary.
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

	// Validate minimum length (measured in runes, not bytes)
	if (element.MinLength > 0) && (utf8.RuneCountInString(value) < element.MinLength) {
		return value, false, derp.Validation("Minimum length is " + convert.String(element.MinLength))
	}

	// Validate maximum length (ACTUALLY TRUNCATES THE STRING TO THE MAXIMUM LENGTH).
	// Length is measured in runes, and truncation happens on a rune boundary so that
	// multi-byte characters are never split into invalid UTF-8.
	if (element.MaxLength > 0) && (utf8.RuneCountInString(value) > element.MaxLength) {
		value = string([]rune(value)[:element.MaxLength])
		changed = true
	}

	// Validate enumerated values
	if len(element.Enum) > 0 {
		if (value != "") && (compare.NotContains(element.Enum, value)) {
			return value, false, derp.Validation("Must be one of the specified values", value, element.Enum)
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
