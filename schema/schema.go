package schema

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
)

// Schema defines a (simplified) JSON-Schema object, that can be Marshalled/Unmarshalled to JSON.
type Schema struct {
	ID      string
	Comment string
	Element Element
}

// New generates a fully initialized Schema object using the provided properties
func New(element Element) Schema {
	return Schema{
		Element: element,
	}
}

/******************************************
 * Validation Methods
 ******************************************/

// Validate checks a particular value against this schema.  If the
// provided value is not valid, then an error is returned.
func (schema Schema) Validate(value any) error {

	if isNil(schema.Element) {
		return derp.NewInternalError("schema.Schema.Validate", "Schema is nil")
	}

	if err := schema.Element.Validate(value); err != nil {
		return derp.Wrap(err, "schema.Schema.Validate", "Error validating value", value)
	}

	return nil
}

// Clean tries to force a particular value to fit this schema by updating
// it (or all of its properties) to match.  If values cannot be coerced to
// fit the schema, then an error is returned
func (schema Schema) Clean(value any) error {

	// TODO: CRITICAL: "Clean" functions are not yet implemented
	if isNil(schema.Element) {
		return derp.NewInternalError("schema.Schema.Clean", "Schema is nil")
	}

	if err := schema.Element.Clean(value); err != nil {
		return derp.Wrap(err, "schema.Schema.Clean", "Error cleaning value", value)
	}

	return nil
}

/******************************************
 * Other Data Access Methods
 ******************************************/

func (schema Schema) GetElement(path string) (Element, bool) {

	if isNil(schema.Element) {
		return nil, false
	}

	return schema.Element.GetElement(path)
}

func (schema *Schema) Inherit(parent Schema) {

	if isNil(schema.Element) {
		schema.Element = parent.Element
	} else {
		schema.Element.Inherit(parent.Element)
	}
}

/******************************************
 * Marshaling Methods
 ******************************************/

// MarshalJSON converts a schema into JSON.
func (schema Schema) MarshalJSON() ([]byte, error) {

	if isNil(schema.Element) {
		return []byte("null"), nil
	}

	return json.Marshal(schema.MarshalMap())
}

// MarshalMap converts a schema into a map[string]any
func (schema Schema) MarshalMap() map[string]any {

	if isNil(schema.Element) {
		return map[string]any{}
	}

	result := schema.Element.MarshalMap()

	if schema.ID != "" {
		result["$id"] = schema.ID
	}

	if schema.Comment != "" {
		result["$comment"] = schema.Comment
	}

	return result
}

/***********************************
 * Marshal / Unmarshal Methods
 ***********************************/

// UnmarshalJSON creates a new Schema object using a JSON-serialized byte array.
func (schema *Schema) UnmarshalJSON(data []byte) error {

	unmarshalled := make(map[string]any)

	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return derp.Wrap(err, "schema.UnmarshalJSON", "Invalid JSON", string(data))
	}

	if err := schema.UnmarshalMap(unmarshalled); err != nil {
		return derp.Wrap(err, "schema.UnmarshalJSON", "Unable to unmarshal from Map", unmarshalled)
	}

	return nil
}

// UnmarshalMap updates a Schema using a map[string]any
func (schema *Schema) UnmarshalMap(data map[string]any) error {

	schema.ID = convert.String(data["$id"])
	schema.Comment = convert.String(data["$comment"])

	element, err := UnmarshalMap(data)

	if err != nil {
		return derp.Wrap(err, "schema.Schema.UnmarshalMap", "Error unmarshalling element")
	}

	schema.Element = element

	return nil
}

/******************************************
 * Other Helpers
 ******************************************/

// isNil performs a robust nil check on an interface
// Shout out to: https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
func isNil(i any) bool {

	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
