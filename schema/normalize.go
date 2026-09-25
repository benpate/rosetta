package schema

import (
	"github.com/benpate/derp"
)

// Normalize rewrites a value in place so that it conforms to the schema, returning the path
// of every rewritten value, or an error when the value cannot be made to conform.
func Normalize(schema Schema, value any) ([]string, error) {

	const location = "schema.Schema.Normalize"

	// RULE: Schema element cannot be nil
	if isNil(schema.Element) {
		return nil, derp.Internal(location, "Schema must not be nil")
	}

	// Validate the value, applying rewrites in place.  Unlike Validate, a value that had to be
	// formatted, clamped, or truncated is not an error.
	_, rewrites, err := validate(schema.Element, value)

	if err != nil {
		return nil, derp.Wrap(err, location, "Value is not valid for this schema", typeName(value))
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return nil, derp.Wrap(err, location, "Validating `required-if` fields", typeName(value))
	}

	return rewrites.paths(), nil
}

// Normalize rewrites a value in place so that it conforms to this schema, returning the path
// of every rewritten value, or an error when the value cannot be made to conform.
func (schema Schema) Normalize(value any) ([]string, error) {
	return Normalize(schema, value)
}
