package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Validate checks a particular value against this schema, updating values when
// possible so that they pass validation.  If the provided value is not valid
// (and cannot be coerced into being valid) then it returns an error.
func Validate(schema Schema, value any) (any, bool, error) {

	const location = "schema.Schema.Validate"

	// RULE: Schema element cannot be nil
	if isNil(schema.Element) {
		return value, false, derp.Internal(location, "Schema must not be nil")
	}

	// Validate the value using the schema's element
	result, rewrites, err := validate(schema.Element, value)

	// On error, report the ORIGINAL value; validate() returns nil results on error.
	if err != nil {
		return result, false, derp.Wrap(err, location, "Value is not valid for this schema", value)
	}

	value = result

	// RULE: A value that had to be rewritten (clamped or formatted) is not valid as-given.
	// Set() can rewrite values in place; Validate() answers whether the value already conforms,
	// so any required modification means the value fails validation.  The error details name
	// each rewritten property with its before/after values.
	if len(rewrites) > 0 {
		return value, false, derp.Validation("Value is not valid for this schema", rewrites.details()...)
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return value, false, derp.Wrap(err, location, "Validating `required-if` fields", value)
	}

	return value, false, nil
}

// validate verifies that the provided value meets the requirements of the schema element,
// and updates the value if necessary.  The returned rewriteList records every value that
// had to be modified; leaf validators report a simple changed flag, which is converted
// here (where both the before and after values are in hand) into a rewrite record.
func validate(element Element, value any) (any, rewriteList, error) {

	const location = "schema.validate"

	switch typedElement := element.(type) {

	case Any:
		result, changed, err := validate_Any(typedElement, value)
		return result, newRewriteList(changed, value, result), err

	case Array:
		return validate_Array(typedElement, value)

	case Boolean:
		result, changed, err := validate_Boolean(typedElement, value)
		return result, newRewriteList(changed, value, result), err

	case Integer:
		result, changed, err := validate_Integer(typedElement, value)
		return result, newRewriteList(changed, value, result), err

	case Number:
		result, changed, err := validate_Number(typedElement, value)
		return result, newRewriteList(changed, value, result), err

	case Object:
		return validate_Object(typedElement, value)

	case String:
		stringValue := convert.String(value)
		result, changed, err := validate_String(typedElement, stringValue)
		return result, newRewriteList(changed, stringValue, result), err
	}

	// This is an invalid element type (this should never happen)
	return nil, nil, derp.Internal(location, "Invalid element type", element)
}
