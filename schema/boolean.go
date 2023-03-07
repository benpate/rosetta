package schema

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/null"
)

// Boolean represents a boolean data type within a JSON-Schema.
type Boolean struct {
	Default  null.Bool `json:"default"`
	Required bool
}

/***********************************
 * Element Interface
 ***********************************/

func (element Boolean) DefaultValue() any {
	return element.Default.Bool()
}

// IsRequired returns TRUE if this element is a required field
func (element Boolean) IsRequired() bool {
	return element.Required
}

// Validate validates a generic value using this schema
func (element Boolean) Validate(object any) error {

	boolValue, ok := object.(bool)

	if !ok {
		return derp.NewValidationError(" must be a boolean")
	}

	if element.Required && (!boolValue) {
		return derp.NewValidationError(" boolean value is required")
	}

	return nil
}

func (element Boolean) Clean(value any) error {
	// TODO: HIGH: Implement the Clean() method for the Boolean element
	return nil
}

func (element Boolean) GetElement(name string) (Element, bool) {

	if name == "" {
		return element, true
	}
	return nil, false
}

func (element Boolean) Inherit(parent Element) {
	// Do nothing
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
