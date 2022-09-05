package schema

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/null"
)

// Boolean represents a boolean data type within a JSON-Schema.
type Boolean struct {
	Default  null.Bool `json:"default"`
	Required bool
}

/***********************************
 * ELEMENT META-DATA
 ***********************************/

// Type returns the data type of this Element
func (element Boolean) Type() reflect.Type {
	return reflect.TypeOf(bool(true))
}

// DefaultValue returns the default value for this element type
func (element Boolean) DefaultValue() any {
	return element.Required
}

// IsRequired returns TRUE if this element is a required field
func (element Boolean) IsRequired() bool {
	return element.Required
}

/***********************************
 * PRIMARY INTERFACE METHODS
 ***********************************/

// Find locates a child of this element
func (element Boolean) Get(object reflect.Value, path list.List) (reflect.Value, error) {

	// RULE: Cannot get sub-properties of a boolean
	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Boolean.Find", "Can't find sub-properties on a 'boolean' type", path)
	}

	// Try to convert and return the new value
	return reflect.ValueOf(convert.BoolDefault(object, element.Default.Bool())), nil
}

// GetElement returns the element definition for a given path
func (element Boolean) GetElement(path list.List) (Element, error) {

	if path.IsEmpty() {
		return element, nil
	}

	return nil, derp.NewInternalError("schema.Boolean.GetElement", "Can't find sub-properties on a boolean", path)
}

// Set formats a value and applies it to the provided object/path
func (element Boolean) Set(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	// RULE: Cannot set sub-properties of a boolean
	if !path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError("schema.Boolean.Set", "Can't set sub-properties on a boolean", path, value)
	}

	// Convert and return the value
	return reflect.ValueOf(convert.BoolDefault(value, element.Default.Bool())), nil
}

// Remove removes a value from the provided object/path.  In the case of booleans, this is a no-op.
func (element Boolean) Remove(_ reflect.Value, _ list.List) (reflect.Value, error) {
	return reflect.ValueOf(nil), derp.NewInternalError("schema.Boolean.Remove", "Can't remove properties from a boolean.  This should never happen.")
}

// Validate validates a generic value using this schema
func (element Boolean) Validate(object any) error {

	boolValue := convert.BoolDefault(object, element.Default.Bool())

	if element.Required && (!boolValue) {
		return ValidationError{Message: "field is required"}
	}

	return nil
}

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

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
		return derp.NewInternalError("schema.Boolean.UnmarshalMap", "Data is not type 'boolean'", data)
	}

	element.Default = convert.NullBool(data["default"])
	element.Required = convert.Bool(data["required"])

	return nil
}
