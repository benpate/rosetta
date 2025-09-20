package schema

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

// Element interface wraps all of the methods required for schema elements.
type Element interface {

	// Default returns the default value for this element
	DefaultValue() any

	// IsRequired returns true if this a value is required for this element
	IsRequired() bool

	// Validate validates the provided value
	Validate(value any) error

	// ValidateRequiredIf handles conditional validation of a required field
	ValidateRequiredIf(schema Schema, path list.List, globalValue any) error

	// MarshalMap populates the object data into a map[string]any
	MarshalMap() map[string]any

	// getElement returns a named sub-element of this element, if it exists.
	GetElement(string) (Element, bool)

	Inherit(Element)

	AllProperties() ElementMap
}

// WritableElement represents an Element (usually a pointer to a concrete type) whose value can be changed.
type WritableElement interface {

	// UnmarshalMap tries to populate this object using data from a map[string]any
	UnmarshalMap(map[string]any) error

	Element
}

// UnmarshalJSON tries to parse a []byte into a schema.Element
func UnmarshalJSON(data []byte) (Element, error) {

	var result map[string]any

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, derp.Wrap(err, "schema.UnmarshalJSON", "Unable to unmarshal JSON", string(data))
	}

	element, err := UnmarshalMap(result)

	if err != nil {
		return nil, derp.Wrap(err, "schema.UnmarshalJSON", "Unable to unmarshal map", string(data))
	}

	if element == nil {
		return nil, derp.InternalError("schema.UnmarshalJSON", "Unmarshalled element is nil", string(data))
	}

	return element, nil
}

// UnmarshalMap tries to parse a map[string]any into a schema.Element
func UnmarshalMap(data any) (Element, error) {

	const location = "schema.UnmarshalMap"

	// NILCHECK: Data cannot be nil
	if data == nil {
		return nil, derp.InternalError(location, "Element is nil")
	}

	// Confirm that the data is a map[string]any
	dataMap, isMap := data.(map[string]any)

	if !isMap {
		return nil, derp.InternalError(location, "data is not map[string]any", data)
	}

	// Convert the map value into the correct element.
	switch Type(convert.String(dataMap["type"])) {

	case TypeArray:
		result := Array{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil

	case TypeBoolean:
		result := Boolean{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil

	case TypeInteger:
		result := Integer{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil

	case TypeNumber:
		result := Number{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil

	case TypeObject:
		result := Object{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil

	case TypeString:
		result := String{}
		if err := result.UnmarshalMap(dataMap); err != nil {
			return nil, err
		}
		return result, nil
	}

	// Fall through to failure.  You should be disappointed.
	return nil, derp.InternalError(location, "Unrecognized data type", data)
}
