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
		derp.Report(err)
		return nil, derp.Wrap(err, "schema.UnmarshalJSON", "Error unmarshalling JSON", string(data))
	}

	element, err := UnmarshalMap(result)

	if err != nil {
		return nil, derp.Wrap(err, "schema.UnmarshalJSON", "Error unmarshalling map", string(data))
	}

	return element, nil
}

// UnmarshalMap tries to parse a map[string]any into a schema.Element
func UnmarshalMap(data any) (Element, error) {

	if data == nil {
		return nil, derp.InternalError("schema.UnmarshalMap", "Element is nil")
	}

	dataMap, ok := data.(map[string]any)

	if !ok {
		return nil, derp.InternalError("schema.UnmarshalMap", "data is not map[string]any", data)
	}

	switch Type(convert.String(dataMap["type"])) {

	case TypeArray:
		result := Array{}
		err := result.UnmarshalMap(dataMap)
		return result, err

	case TypeBoolean:
		result := Boolean{}
		err := result.UnmarshalMap(dataMap)
		return result, err

	case TypeInteger:
		result := Integer{}
		err := result.UnmarshalMap(dataMap)
		return result, err

	case TypeNumber:
		result := Number{}
		err := result.UnmarshalMap(dataMap)
		return result, err

	case TypeObject:
		result := Object{}
		err := result.UnmarshalMap(dataMap)
		return result, err

	case TypeString:
		result := String{}
		err := result.UnmarshalMap(dataMap)
		return result, err
	}

	// Fall through to failure.  You should be sad.
	return nil, derp.InternalError("schema.UnmarshalElement", "Unrecognized data type", data)

}
