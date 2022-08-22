package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Any represents a any data type within a JSON-Schema.
type Any struct {
	Required bool
}

// Type returns the data type of this Element
func (element Any) Type() reflect.Type {
	return reflect.TypeOf(interface{}(nil))
}

// IsRequired returns TRUE if this element is a required field
func (element Any) IsRequired() bool {
	return element.Required
}

// Get returns the value of this element or its children
func (element Any) Get(object reflect.Value, path string) (any, Element, error) {

	if path != "" {
		return nil, element, derp.NewInternalError("schema.Any.Find", "Can't find sub-properties on an 'any' type", path)
	}

	return convert.Interface(object), element, nil
}

// Set formats a value and applies it to the provided object/path
func (element Any) Set(object reflect.Value, path string, value any) error {

	if path != "" {
		return derp.NewInternalError("schema.Any.Set", "Can't set sub-properties on an 'any' type", path, value)
	}

	return setWithReflection(object, value)
}

// Validate validates a value against this schema
func (element Any) Validate(value any) error {

	if element.Required {
		if convert.IsZeroValue(value) {
			return ValidationError{Message: "field is required"}
		}
	}

	return nil
}

// DefaultType returns the default type for this element
func (element Any) DefaultType() reflect.Type {
	return reflect.TypeOf(interface{}(nil))
}

// DefaultValue returns the default value for this element type
func (element Any) DefaultValue() any {
	return any(nil)
}

// MarshalMap populates object data into a map[string]any
func (element Any) MarshalMap() map[string]any {
	return map[string]any{
		"type": TypeAny,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Any) UnmarshalMap(data map[string]any) error {

	if convert.String(data["type"]) != "any" {
		return derp.New(500, "schema.Any.UnmarshalMap", "Data is not type 'any'", data)
	}

	element.Required = convert.Bool(data["required"])

	return nil
}
