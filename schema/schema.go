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
func New(properties ElementMap) Schema {
	return Schema{
		Element: Object{
			Properties: properties,
		},
	}
}

// Get retrieves a generic value from the object, along with
// the schema element that defines it.  If the object is nil, Get still
// tries to return a default value if provided by the schema
func (schema Schema) Get(object any, path string) (any, Element, error) {

	if schema.Element == nil {
		return nil, nil, derp.NewInternalError("schema.Schema.Get", "Invalid schema.  Element is nil")
	}

	return schema.Element.Get(reflect.ValueOf(object), path)
}

// GetBool retrieves a bool value from this object.  If the value
// is not defined in the object/schema, then the zero value (false) is returned
func (schema Schema) GetBool(object any, path string) bool {
	result, _, _ := schema.Get(object, path)
	return convert.Bool(result)
}

// GetFloat retrieves a float64 value from this object.  If the value
// is not defined in the object/schema, then the zero value (0.0) is returned
func (schema Schema) GetFloat(object any, path string) float64 {
	result, _, _ := schema.Get(object, path)
	return convert.Float(result)
}

// GetInt retrieves a int value from this object.  If the value
// is not defined in the object/schema, then the zero value (0) is returned
func (schema Schema) GetInt(object any, path string) int {
	result, _, _ := schema.Get(object, path)
	return convert.Int(result)
}

// GetInt64 retrieves an int64 value from this object.  If the value
// is not defined in the object/schema, then the zero value (0) is returned
func (schema Schema) GetInt64(object any, path string) int64 {
	result, _, _ := schema.Get(object, path)
	return convert.Int64(result)
}

// GetString retrieves a string value from this object.  If the value
// is not defined in the object/schema, then the zero value ("") is returned
func (schema Schema) GetString(object any, path string) string {
	result, _, _ := schema.Get(object, path)
	return convert.String(result)
}

// SetAll iterates over Set to apply all of the values to the object one at a time, stopping
// at the first error it encounters.  If all values are addedd successfully, then SetAll
// also uses Validate() to confirm that the object is still correct.
func (schema Schema) SetAll(object any, values map[string]any) error {

	const location = "schema.Schema.SetAll"

	// Set each value in the schema
	for path, value := range values {

		// Errors are intentionally ignored here.
		// Unallowed data does not make it through the schema filter
		schema.Set(object, path, value)
	}

	// Validate the whole schema once all the values are set
	if err := schema.Validate(object); err != nil {
		return derp.Wrap(err, location, "Validation Error")
	}

	// Success!!
	return nil
}

// Schema applies a value to the object at the given path.  If the path is invalid
// then it returns an error
func (schema Schema) Set(object any, path string, value any) error {

	const location = "schema.Schema.Set"

	// Shortcut if the object is a PathSetter.  Just call the SetPath function and we're good.
	if setter, ok := object.(PathSetter); ok {
		return setter.SetPath(path, value)
	}

	if schema.Element == nil {
		return derp.NewInternalError("schema.Schema.Set", "Invalid schema.  Element is nil.")
	}

	valueOf := reflect.ValueOf(object)

	// Guarantee that we've been passed a pointer
	if valueOf.Kind() != reflect.Pointer {
		return derp.NewInternalError(location, "Must pass a pointer (not a value) to this function.", object, path, value)
	}

	// Now that we KNOW it's a pointer, dereference
	addressable := valueOf.Elem()

	// Verify that it's still addressable (this should never fail)
	if !addressable.CanSet() {
		return derp.NewInternalError(location, "Cannot set value")
	}

	// Try to set the value in the variable
	if err := schema.Element.SetReflect(addressable, path, value); err != nil {
		return derp.Wrap(err, location, "Error setting value")
	}

	return nil

}

// Validate checks a particular value against this schema.  If the
// provided value is not valid, then an error is returned.
func (schema Schema) Validate(value any) error {

	if schema.Element == nil {
		return derp.NewInternalError("schema.Schema.Validate", "Schema is nil")
	}

	return schema.Element.Validate(value)
}

// MarshalJSON converts a schema into JSON.
func (schema Schema) MarshalJSON() ([]byte, error) {

	if schema.Element == nil {
		return []byte("null"), nil
	}

	return json.Marshal(schema.MarshalMap())
}

// MarshalMap converts a schema into a map[string]any
func (schema Schema) MarshalMap() map[string]any {

	if schema.Element == nil {
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

	var err error

	schema.ID = convert.String(data["$id"])
	schema.Comment = convert.String(data["$comment"])
	schema.Element, err = UnmarshalMap(data)

	if err != nil {
		return derp.Wrap(err, "schema.Schema.UnmarshalMap", "Error unmarshalling element")
	}

	return nil
}
