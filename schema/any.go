package schema

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

type Any struct {
	Required bool
}

// Default returns the default value for this element
func (element Any) DefaultValue() any {
	return nil
}

// IsRequired returns true if this a value is required for this element
func (element Any) IsRequired() bool {
	return element.Required
}

// Validate validates the provided value
func (element Any) Validate(value any) error {
	return nil
}

// Clean updates a value to match the schema.  The value must be a pointer.
func (element Any) Clean(value any) error {
	return nil
}

// MarshalMap populates the object data into a map[string]any
func (element Any) MarshalMap() map[string]any {

	return map[string]any{
		"type":     TypeAny,
		"required": element.Required,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Any) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "any" {
		return derp.NewInternalError("schema.String.UnmarshalMap", "Data is not type 'string'", data)
	}

	element.Required = convert.Bool(data["required"])

	return err
}
