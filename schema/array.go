package schema

import (
	"reflect"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/null"
)

// Array represents an array data type within a JSON-Schema.
type Array struct {
	Items     Element
	MinLength null.Int
	MaxLength null.Int
	Required  bool
	Delimiter string
}

// Type returns the reflection type of this Element
func (element Array) Type() reflect.Type {
	return reflect.SliceOf(element.Items.Type())
}

// IsRequired returns TRUE if this element is a required field
func (element Array) IsRequired() bool {
	return element.Required
}

// Find locates a child of this element
func (element Array) Get(object reflect.Value, path string) (any, Element, error) {

	// Validate that we have the right type of object
	switch object.Kind() {
	case reflect.Array, reflect.Slice:
		// Move along, these types are good.

	case reflect.Interface, reflect.Pointer:
		// Dereferenced pointers
		return element.Get(object.Elem(), path)

	case reflect.String:
		// Strings can be split into arrays
		object = reflect.ValueOf(convert.SplitSliceOfString(object, element.Delimiter))

	default:
		// All other types are invalid.
		return nil, element, derp.NewBadRequestError("schema.Array.Get", "Value must be an array, slice, or a string that can be split into an array.", object.Kind(), path, object.Interface())
	}

	// If the request is for this object, then convert it from
	if path == "" {
		return convert.Interface(object), element, nil
	}

	// Finf (and validate) the requested array index
	head, tail := list.Dot(path).Split()
	index, err := strconv.Atoi(head)

	if err != nil {
		return nil, element, derp.NewBadRequestError("schema.Array.Get", "Invalid index (not an integer)", path)
	}

	if index < 0 {
		return nil, element, derp.NewBadRequestError("schema.Array.Get", "Invalid index (less than zero)", path)
	}

	if index >= object.Len() {
		return nil, element, derp.NewBadRequestError("schema.Array.Find", "Invalid index (overflow)", path)
	}

	result := object.Index(index)
	return element.Items.Get(result, string(tail))
}

// Set formats/validates a generic value using this schema
func (element Array) Set(object reflect.Value, path string, value any) error {

	const location = "schema.Array.Set"
	var err error

	// Catch any reflection panics
	defer func() {
		if r := recover(); r != nil {
			err = derp.NewInternalError(location, "Error in reflection", r)
		}
	}()

	// Try to set the value directly.  This will barf if the types don't match
	// but it's better to try and fail, than to not try at all.
	if path == "" {
		object.Set(reflect.ValueOf(value))
		return err
	}

	// Try to calculate the array index
	head, tail := list.Dot(path).Split()
	index, err := strconv.Atoi(head)

	if err != nil {
		return derp.NewBadRequestError(location, "Invalid array index", head)
	}

	if index < 0 {
		return derp.NewBadRequestError(location, "Index out of bounds (negative index)", index)
	}

	// Verify that the index has not overflowed the slice/array bounds
	switch object.Kind() {

	case reflect.Array:

		// If the index is too large, then error because we cannot increase the size of the array
		if index >= object.Len() {
			return derp.NewInternalError(location, "Index out of bounds (array overflow)", index)
		}

	case reflect.Slice:

		// If the index is too large, then increase the size of the slice
		minLength := index + 1
		if needed := minLength - object.Len(); needed > 0 {
			emptySlice := reflect.MakeSlice(element.Type(), needed, needed)
			newSlice := reflect.AppendSlice(object, emptySlice)
			object.Set(newSlice)
		}

	// If we have a "nil" object, then make a new slice of the required size
	case reflect.Invalid:
		minLength := index + 1
		object.Set(reflect.MakeSlice(element.Type(), minLength, minLength))

	default:
		return derp.NewInternalError(location, "Value must be an array")
	}

	// Try to set the result into this object
	subObject := object.Index(index)

	if err := element.Items.Set(subObject, string(tail), value); err != nil {
		return derp.Wrap(err, location, "Error setting array index")
	}

	// Should be nil (unless reflection panic)
	return err
}

// Validate validates a value against this schema
func (element Array) Validate(value any) error {

	var errorReport error

	v := reflect.ValueOf(value)

	// Verify that the object is an array or slice
	if kind := v.Kind(); (kind != reflect.Array) && (kind != reflect.Slice) {
		return Invalid("Element must be  an array")
	}

	// Check minimum/maximum lengths
	length := v.Len()

	if element.Required && length == 0 {
		return Invalid("field is required")
	}

	if element.MinLength.IsPresent() && (length < element.MinLength.Int()) {
		return Invalid("minimum length is " + element.MinLength.String())
	}

	if element.MaxLength.IsPresent() && (length > element.MaxLength.Int()) {
		return Invalid("maximum length is " + element.MaxLength.String())
	}

	// Verify that each item in the array/slice is also valid
	for index := 0; index < length; index = index + 1 {

		item := v.Index(index).Interface()
		if err := element.Items.Validate(item); err != nil {
			errorReport = derp.Append(errorReport, addPath(convert.String(index), err))
		}
	}

	return errorReport
}

// MarshalMap populates object data into a map[string]any
func (element Array) MarshalMap() map[string]any {

	return map[string]any{
		"type":      TypeArray,
		"items":     element.Items.MarshalMap(),
		"minLength": element.MinLength,
		"maxLength": element.MaxLength,
	}
}

// UnmarshalMap tries to populate this object using data from a map[string]any
func (element *Array) UnmarshalMap(data map[string]any) error {

	var err error

	if convert.String(data["type"]) != "array" {
		return derp.New(500, "schema.Array.UnmarshalMap", "Data is not type 'array'", data)
	}

	element.Items, err = UnmarshalMap(data["items"])
	element.Required = convert.Bool(data["required"])
	element.MinLength = convert.NullInt(data["minLength"])
	element.MaxLength = convert.NullInt(data["maxLength"])

	return err
}
