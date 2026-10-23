package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Validate checks a value against the schema, returning an error when it does not conform
// as given, including when a value had to be rewritten to conform.
func Validate(schema Schema, value any) (any, bool, error) {

	const location = "schema.Schema.Validate"

	// RULE: Schema element cannot be nil
	if isNil(schema.Element) {
		return value, false, derp.Internal(location, "Schema must not be nil")
	}

	// Validate the value using the schema's element
	rootValue, restore := addressableRoot(schema.Element, value)
	result, rewrites, err := validate(schema.Element, rootValue)

	// On error, name the ORIGINAL value's type; validate() returns nil results on error.
	if err != nil {
		return result, false, derp.Wrap(
			err, location, "Value is not valid for this schema", typeName(value),
		)
	}

	value = restore(result)

	// RULE: A value that had to be rewritten (clamped or formatted) is not valid as-given, because
	// Validate() answers whether the value already conforms.  The error details name each
	// rewritten property with its before/after values.
	if len(rewrites) > 0 {
		return value, false, derp.Validation("Value is not valid for this schema", rewrites.details()...)
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return value, false, derp.Wrap(
			err, location, "Validating `required-if` fields", typeName(value),
		)
	}

	return value, false, nil
}

// addressableRoot prepares the top-level value for validation: a map held by value is validated
// through a pointer to a copy, as a nested map is, and anything else is returned unchanged
func addressableRoot(element Element, value any) (any, func(any) any) {

	// RULE: Only a map.  A struct held by value stays unwritable, because Emissary's
	// newStartupStream relies on SetAll refusing a whole model.Content
	if reflect.ValueOf(value).Kind() != reflect.Map {
		return value, func(validated any) any { return validated }
	}

	return addressableItem(element, value)
}

// validate checks a value against a schema element and updates it where necessary,
// returning a rewriteList that records every value that had to be modified
func validate(element Element, value any) (any, rewriteList, error) {

	const location = "schema.validate"

	// Leaf validators report only a changed flag.  It becomes a rewrite record here, where
	// both the before and after values are in hand.

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

		// RULE: An Object describes a map or a struct, so a single value or a list is refused
		if err := requireMapKind(value); err != nil {
			return nil, nil, err
		}

		return validate_Object(typedElement, value)

	case String:
		stringValue := convert.String(value)
		result, changed, err := validate_String(typedElement, stringValue)
		return result, newRewriteList(changed, stringValue, result), err
	}

	// This is an invalid element type (this should never happen)
	return nil, nil, derp.Internal(location, "Invalid element type", element)
}
