package schema

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/null"
)

// Boolean represents a boolean data type within a JSON-Schema.
type Boolean struct {
	Default  null.Bool `json:"default"`
	Required bool
}

// Type returns the data type of this Element
func (element Boolean) Type() reflect.Type {
	return reflect.TypeOf(true)
}

// IsRequired returns TRUE if this element is a required field
func (element Boolean) IsRequired() bool {
	return element.Required
}

func (element Boolean) Get(object any, path string) (any, Element, error) {
	return element.GetReflect(convert.ReflectValue(object), path)
}

// Find locates a child of this element
func (element Boolean) GetReflect(object reflect.Value, path string) (any, Element, error) {

	if path != "" {
		return nil, element, derp.NewInternalError("schema.Boolean.Find", "Can't find sub-properties on a 'boolean' type", path)
	}

	return convert.BoolDefault(object, element.Default.Bool()), element, nil
}

// Set formats a value and applies it to the provided object/path
func (element Boolean) Set(object any, path string, value any) error {

	// Shortcut if the object is a PathSetter.  Just call the SetPath function and we're good.
	if setter, ok := object.(PathSetter); ok {
		return setter.SetPath(path, value)
	}

	return element.SetReflect(convert.ReflectValue(object), path, value)
}

// Set formats a value and applies it to the provided object/path
func (element Boolean) SetReflect(object reflect.Value, path string, value any) error {

	// Cannot set sub-properties of a boolean
	if path != "" {
		return derp.NewInternalError("schema.Boolean.Set", "Can't set sub-properties on a boolean", path, value)
	}

	// Convert and set the value
	return setWithReflection(object, convert.BoolDefault(value, element.Default.Bool()))
}

// Validate validates a generic value using this schema
func (element Boolean) Validate(object any) error {

	boolValue := convert.BoolDefault(object, element.Default.Bool())

	if element.Required && (!boolValue) {
		return ValidationError{Message: "field is required"}
	}

	return nil
}

// DefaultValue returns the default value for this element type
func (element Boolean) DefaultValue() any {
	return element.Required
}

// MarshalJSON implements the json.Marshaler interface
func (element Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(element.MarshalMap())
}

// MarshalMap populates object data into a map[string]any
func (element Boolean) MarshalMap() map[string]any {

	result := map[string]any{
		"type": TypeBoolean,
	}

	if element.Default.IsPresent() {
		result["default"] = element.Default.Bool()
	}

	if element.Required {
		result["required"] = true
	}

	return result
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Boolean) UnmarshalMap(data map[string]any) error {

	if convert.String(data["type"]) != "boolean" {
		return derp.New(500, "schema.Boolean.UnmarshalMap", "Data is not type 'boolean'", data)
	}

	element.Default = convert.NullBool(data["default"])
	element.Required = convert.Bool(data["required"])

	return nil
}
