package schema

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Element interface wraps all of the methods required for schema elements.
type Element interface {

	// Type returns the Type of this particular schema element
	Type() reflect.Type

	// Find uses the path to locate a value in an object along with the schema that defines it.
	Get(object reflect.Value, path string) (any, Element, error)

	// Get retrieves a value from the data object
	// Get(object reflect.Value, path string) (any, error)

	// Set formats a value and applies it to the provided object/path
	Set(object reflect.Value, path string, value any) error

	// Validate validates the provided value
	Validate(value any) error

	// MarshalMap populates the object data into a map[string]any
	MarshalMap() map[string]any

	IsRequired() bool
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

	var result WritableElement

	if data == nil {
		return nil, derp.New(500, "schema.UnmarshalMap", "Element is nil")
	}

	dataMap, ok := data.(map[string]any)

	if !ok {
		return nil, derp.New(500, "schema.UnmarshalMap", "data is not map[string]any", data)
	}

	switch Type(convert.String(dataMap["type"])) {

	case TypeAny:
		result = &Any{}

	case TypeArray:
		result = &Array{}

	case TypeBoolean:
		result = &Boolean{}

	case TypeInteger:
		result = &Integer{}

	case TypeNumber:
		result = &Number{}

	case TypeObject:
		result = &Object{}

	case TypeString:
		result = &String{}

	default:
		return nil, derp.New(500, "schema.UnmarshalElement", "Unrecognized data type", data)

	}

	err := result.UnmarshalMap(dataMap)
	return result, err

}
