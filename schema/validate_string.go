package schema

import (
	"regexp"
	"unicode/utf8"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
)

// validate_String checks that the provided value meets the requirements of the String schema element, and updates the value if necessary.
func validate_String(element String, value string) (string, bool, error) {

	// Steps that REWRITE the value run first, in an order that reaches a stable result:
	// format the value, truncate it to the maximum length, then format once more (truncation
	// can leave an artifact, such as a trailing space, that the formatter would still collapse).
	// Only then do the accept/reject rules run, against the value as it will actually be stored,
	// so that a value produced by Set always passes a subsequent Validate.

	const location = "schema.validate_String"

	// Remember the original so we can report whether any rewrite step changed the value.
	original := value

	// REWRITE: apply formatting functions (they may shrink or rewrite the value).
	value, err := applyStringFormats(element, value)
	if err != nil {
		return value, false, derp.Wrap(err, location, "Applying format functions", value)
	}

	// REWRITE: truncate to the maximum length (measured in runes, not bytes). Truncation
	// happens on a rune boundary so that multi-byte characters are never split into invalid UTF-8.
	if (element.MaxLength > 0) && (utf8.RuneCountInString(value) > element.MaxLength) {
		value = string([]rune(value)[:element.MaxLength])
	}

	// REWRITE: format again, because truncation can re-introduce something the formatter
	// removes (e.g. a trailing space at the cut boundary). After this the value is stable.
	value, err = applyStringFormats(element, value)
	if err != nil {
		return value, false, derp.Wrap(err, location, "Applying format functions", value)
	}

	changed := (value != original)

	// All rewriting is complete; the remaining rules only accept or reject the final value.

	// Verify required fields
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

	// Validate enumerated values
	if len(element.Enum) > 0 {
		if (value != "") && (compare.NotContains(element.Enum, value)) {
			return value, false, derp.Validation("Must be one of the specified values", value, element.Enum)
		}
	}

	// Validate regex Pattern (accept/reject only; it never rewrites the value)
	if element.Pattern != "" {
		if matched, err := regexp.MatchString(element.Pattern, value); err != nil {
			return value, false, derp.Wrap(err, location, "Evaluating pattern", element.Pattern)
		} else if !matched {
			return value, false, derp.Validation("Value must match pattern", value, element.Pattern)
		}
	}

	// Successfully validated/updated the value.
	return value, changed, nil
}

// applyStringFormats runs every formatting function for the element over the value in order,
// returning the rewritten value or the first formatting error.
func applyStringFormats(element String, value string) (string, error) {
	for _, formatFunc := range element.formatFunctions() {
		formatted, err := formatFunc(value)
		if err != nil {
			return value, err
		}
		value = formatted
	}
	return value, nil
}
