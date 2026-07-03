package schema

import (
	"github.com/benpate/derp"
)

// Normalize applies the schema to a value, rewriting it in place (formatting,
// clamping, truncating) as needed so that it conforms.  Unlike Validate, a value
// that had to be rewritten is not an error; Normalize returns the paths of every
// rewritten value so that callers can report the changes.  It returns an error
// only when a value cannot be made to conform.
func Normalize(schema Schema, value any) ([]string, error) {

	const location = "schema.Schema.Normalize"

	// RULE: Schema element cannot be nil
	if isNil(schema.Element) {
		return nil, derp.Internal(location, "Schema must not be nil")
	}

	// Validate the value, applying rewrites in place
	_, rewrites, err := validate(schema.Element, value)

	if err != nil {
		return nil, derp.Wrap(err, location, "Value is not valid for this schema", value)
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return nil, derp.Wrap(err, location, "Validating `required-if` fields", value)
	}

	return rewrites.paths(), nil
}

// Normalize applies this schema to the provided value, rewriting it in place so
// that it conforms.  It returns the paths of every rewritten value, and an error
// only when a value cannot be made to conform.
func (schema Schema) Normalize(value any) ([]string, error) {
	return Normalize(schema, value)
}
