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
	value, updated, err := validate(schema.Element, value)

	if err != nil {
		return value, false, derp.Wrap(err, location, "Value is not valid for this schema", value)
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return value, false, derp.Wrap(err, location, "Unable to validate `required-if` fields", value)
	}

	return value, updated, nil
}

// validate verifies that the provided value meets the requirements of the schema element,
// and updates the value if necessary.
func validate(element Element, value any) (any, bool, error) {

	const location = "schema.validate"

	switch typedElement := element.(type) {

	case Any:
		return validate_Any(typedElement, value)

	case Array:
		return validate_Array(typedElement, value)

	case Boolean:
		return validate_Boolean(typedElement, convert.Bool(value))

	case Integer:
		return validate_Integer(typedElement, convert.Int(value))

	case Number:
		return validate_Number(typedElement, convert.Float(value))

	case Object:
		return validate_Object(typedElement, value)

	case String:
		return validate_String(typedElement, convert.String(value))
	}

	// This is an invalid element type (this should never happen)
	return nil, false, derp.Internal(location, "Invalid element type", element)
}
