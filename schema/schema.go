// Package schema provides tools for defining, validating, and manipulating JSON-Schema-like structures in Go.
package schema

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
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

// Wildcard returns a Schema that accepts any data type
func Wildcard() Schema {
	return Schema{
		Element: Any{},
	}
}

/******************************************
 * Validation Methods
 ******************************************/

// Validate checks a particular value against this schema, updating values when
// possible so that they pass validation.  If the provided value is not valid
// (and cannot be coerced into being valid) then it returns an error.
func (schema Schema) Validate(value any) error {

	const location = "schema.Schema.Validate"

	// RULE: Schema element cannot be nil
	if isNil(schema.Element) {
		return derp.Internal(location, "Schema must not be nil")
	}

	// Validate all elements in the value
	if err := schema.Element.Validate(value); err != nil {
		return derp.Wrap(err, location, "Value is not valid for this schema", value)
	}

	// Handle special cases for "required-if" fields
	if err := schema.ValidateRequiredIf(value); err != nil {
		return derp.Wrap(err, location, "Unable to validate `required-if` fields", value)
	}

	return nil
}

// Match returns TRUE if the provided value (as accessed via this schema) matches
// the provided expression.  This is useful for server-side data validation.
func (schema Schema) Match(value any, expression exp.Expression) (bool, error) {

	var resultError error

	evaluatePredicate := func(predicate exp.Predicate) bool {

		fieldValue, err := schema.Get(value, predicate.Field)

		if err != nil {
			resultError = derp.Wrap(err, "schema.schema.Match", "schema value not found", predicate.Field)
			return false
		}

		switch predicate.Operator {

		case exp.OperatorEqual:
			return compare.Equal(fieldValue, predicate.Value)

		case exp.OperatorGreaterThan:
			return compare.GreaterThan(fieldValue, predicate.Value)

		case exp.OperatorLessThan:
			return compare.LessThan(fieldValue, predicate.Value)

		case exp.OperatorNotEqual:
			return !compare.Equal(fieldValue, predicate.Value)

		case exp.OperatorGreaterOrEqual:
			return !compare.LessThan(fieldValue, predicate.Value)

		case exp.OperatorLessOrEqual:
			return !compare.GreaterThan(fieldValue, predicate.Value)
		}

		return false
	}

	// Evaluate the predicate
	result := expression.Match(evaluatePredicate)

	return result, resultError
}

// ValidateRequiredIf implements the Element interface
// It returns an error if the conditional expression is true but the value is empty
func (schema Schema) ValidateRequiredIf(value any) error {
	return schema.Element.ValidateRequiredIf(schema, list.ByDot(""), value)
}

/******************************************
 * Other Data Access Methods
 ******************************************/

// GetElement returns the element at the specified path, or FALSE if the element does not exist
func (schema Schema) GetElement(path string) (Element, bool) {

	if isNil(schema.Element) {
		return nil, false
	}

	return schema.Element.GetElement(path)
}

// GetArrayElement returns the array element at the specified path, or FALSE if invalid
func (schema Schema) GetArrayElement(path string) (Array, bool) {

	if element, ok := schema.GetElement(path); ok {
		switch typed := element.(type) {
		case Array:
			return typed, true
		case Any:
			return Array{Items: Any{}}, true
		}
	}

	return Array{}, false
}

// GetBooleanElement returns the boolean element at the specified path, or FALSE if invalide
func (schema Schema) GetBooleanElement(path string) (Boolean, bool) {

	if element, ok := schema.GetElement(path); ok {
		if booleanElement, ok := element.(Boolean); ok {
			return booleanElement, true
		}
	}

	return Boolean{}, false
}

// GetIntegerElement returns the integer element at the specified path, or FALSE if invalid
func (schema Schema) GetIntegerElement(path string) (Integer, bool) {

	if element, ok := schema.GetElement(path); ok {
		if integerElement, ok := element.(Integer); ok {
			return integerElement, true
		}
	}

	return Integer{}, false
}

// GetNumberElement returns the number element at the specified path, or FALSE if invalid
func (schema Schema) GetNumberElement(path string) (Number, bool) {

	if element, ok := schema.GetElement(path); ok {
		if numberElement, ok := element.(Number); ok {
			return numberElement, true
		}
	}

	return Number{}, false
}

// GetObjectElement returns the object element at the specified path, or FALSE if invalid
func (schema Schema) GetObjectElement(path string) (Object, bool) {

	if element, ok := schema.GetElement(path); ok {
		if objectElement, ok := element.(Object); ok {
			return objectElement, true
		}
	}

	return Object{}, false
}

// GetStringElement returns the string element at the specified path, or FALSE if invalid
func (schema Schema) GetStringElement(path string) (String, bool) {

	if element, ok := schema.GetElement(path); ok {
		if stringElement, ok := element.(String); ok {
			return stringElement, true
		}
	}

	return String{}, false
}

// Inherit updates this schema with properties from the parent schema
func (schema *Schema) Inherit(parent Schema) {

	if isNil(schema.Element) {
		schema.Element = parent.Element
	} else {
		schema.Element.Inherit(parent.Element)
	}
}

// AllProperties returns a flat slice of all properties in this schema
func (schema *Schema) AllProperties() ElementMap {
	return schema.Element.AllProperties()
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
