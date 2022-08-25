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
	Delimiter string // DEPRECATED
}

/***********************************
 * ELEMENT META-DATA
 ***********************************/

// Type returns the reflection type of this Element
func (element Array) Type() reflect.Type {
	return reflect.SliceOf(element.Items.Type())
}

// DefaultValue returns the default value for this element type
func (element Array) DefaultValue() any {
	sliceType := reflect.SliceOf(element.Items.Type())
	return reflect.MakeSlice(sliceType, 0, 0).Interface()
}

// IsRequired returns TRUE if this element is a required field
func (element Array) IsRequired() bool {
	return element.Required
}

/***********************************
 * PRIMARY INTERFACE METHODS
 ***********************************/

func (element Array) Get(object reflect.Value, path list.List) (reflect.Value, Element, error) {

	const location = "schema.Array.Get"

	// Validate that we have the right type of object
	switch object.Kind() {

	// If the value is invalid (nil) then return nil
	case reflect.Invalid:
		return reflect.ValueOf(element.DefaultValue()), element, nil

	// Dereference interfaces
	case reflect.Interface:
		return element.Get(object.Elem(), path)

	// Dereference pointers
	case reflect.Pointer:
		return element.Get(object.Elem(), path)

	// Move along, these types are good.
	case reflect.Array, reflect.Slice:

	// All other types are invalid.
	default:
		return reflect.ValueOf(nil), element, derp.NewBadRequestError(location, "Value must be an array, slice.", object.Kind(), path, object.Interface())
	}

	// If the request is for this object, then convert it from
	if path.IsEmpty() {
		return object, element, nil
	}

	// Get (and bounds-check) the array index
	head, tail := path.Split()
	index, err := strconv.Atoi(head)

	if err != nil {
		return reflect.ValueOf(nil), element, derp.NewBadRequestError("schema.Array.Get", "Invalid index (not an integer)", path)
	}

	if index < 0 {
		return reflect.ValueOf(nil), element, derp.NewBadRequestError("schema.Array.Get", "Invalid index (less than zero)", path)
	}

	if index >= object.Len() {
		return reflect.ValueOf(nil), element, derp.NewBadRequestError("schema.Array.Find", "Invalid index (overflow)", path)
	}

	//
	subValue := object.Index(index)
	return element.Items.Get(subValue, tail)
}

// Set formats/validates a generic value using this schema
func (element Array) Set(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	const location = "schema.Array.Set"

	// Validate that we have the right type of object
	switch object.Kind() {

	// If the value is invalid (nil) then return nil
	case reflect.Invalid:
		return element.Set(reflect.ValueOf(element.DefaultValue()), path, value)

	// Dereference interfaces
	case reflect.Interface:
		return element.Set(object.Elem(), path, value)

	// Dereference pointers
	case reflect.Pointer:
		return element.Set(object.Elem(), path, value)
	}

	// Try to set the value directly.  This will barf if the types don't match
	// but it's better to try and fail, than to not try at all.
	if path.IsEmpty() {
		return object, nil
	}

	// Get (and bounds-check) the array index
	head, tail := path.Split()
	index, err := strconv.Atoi(head)

	if err != nil {
		return reflect.ValueOf(nil), derp.NewBadRequestError(location, "Invalid array index", head)
	}

	if index < 0 {
		return reflect.ValueOf(nil), derp.NewBadRequestError(location, "Index out of bounds (negative index)", index)
	}

	// Bounds check index overflow...
	switch object.Kind() {

	// Bounds check for Arrays: return an error if the index is too big.
	case reflect.Array:

		// If the index is too large, then error because we cannot increase the size of the array
		if index >= object.Len() {
			return reflect.ValueOf(nil), derp.NewInternalError(location, "Index out of bounds (array overflow)", index)
		}

	// Bounds check for slices: grow the slice if the index is too big.
	case reflect.Slice:

		minLength := index + 1
		if needed := minLength - object.Len(); needed > 0 {
			elementType := object.Type().Elem()
			sliceType := reflect.SliceOf(elementType)
			emptySlice := reflect.MakeSlice(sliceType, needed, needed)
			object = reflect.AppendSlice(object, emptySlice)
		}

	default:
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Value must be an array")
	}

	// Get the value in the array (should be addressable)
	subObject := object.Index(index)

	// Try to set the value of the indexed sub-object
	subResult, err := element.Items.Set(subObject, tail, value)

	if err != nil {
		return reflect.ValueOf(nil), derp.Wrap(err, location, "Error setting array index")
	}

	// Put the value back into the array/slice
	subObject.Set(subResult)

	// Done
	return object, nil
}

func (element Array) Remove(object reflect.Value, path list.List) (reflect.Value, error) {

	const location = "schema.Array.Remove"

	// Validate that we have the right type of object
	switch object.Kind() {

	// If the value is invalid (nil) then use the default value
	case reflect.Invalid:
		return element.Remove(reflect.ValueOf(element.DefaultValue()), path)

	// Dereference interfaces
	case reflect.Interface:
		return element.Remove(object.Elem(), path)

	// Dereference pointers
	case reflect.Pointer:
		return element.Remove(object.Elem(), path)

	// Allow processing of arrays and slices to continue
	case reflect.Array, reflect.Slice:

	// All other types are an error.
	default:
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Cannot remove from this type", object.Kind().String(), object.Interface(), path)
	}

	// Cannot remove arrays directly.  This should never happen
	// because this operation SHOULD HAVE been handled by the upstream element.
	if path.IsEmpty() {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Cannot remove array directly.  This should never happen")
	}

	// Get (and bounds-check) the array index
	head, tail := path.Split()
	index, err := strconv.Atoi(head)

	if err != nil {
		return reflect.ValueOf(nil), derp.NewBadRequestError(location, "Invalid array index", head)
	}

	if index < 0 {
		return reflect.ValueOf(nil), derp.NewBadRequestError(location, "Index out of bounds (negative index)", index)
	}

	length := object.Len()
	if index > length {
		return object, nil
	}

	// If we're removing a sub-value, then pass this call to the sub-element.
	if !tail.IsEmpty() {
		subValue := object.Index(index)
		subResult, err := element.Items.Remove(subValue, tail)

		if err != nil {
			return reflect.ValueOf(nil), derp.Wrap(err, location, "Error removing elements in array index", index)
		}

		subValue.Set(subResult)
		return object, nil
	}

	// Otherwise, we're removing a whole element in this array

	// Create a new value (array or slice) that is one item shorter
	var result reflect.Value
	elementType := object.Type().Elem()

	switch object.Kind() {

	case reflect.Array:
		resultType := reflect.ArrayOf(length-1, elementType)
		result = reflect.Zero(resultType)

	case reflect.Slice:
		resultType := reflect.SliceOf(elementType)
		result = reflect.MakeSlice(resultType, length-1, length-1)

	default:
		return reflect.ValueOf(nil), derp.NewInternalError("slice.RemoveFromInterface", "Value must be an array or slice", object.Kind().String(), object.Interface())
	}

	// Copy items from the old value into the new value, skipping the item at the specified index
	newIndex := 0
	for oldIndex := 0; oldIndex < length; oldIndex++ {

		if oldIndex == index {
			continue
		}

		result.Index(newIndex).Set(object.Index(oldIndex))
		newIndex += 1
	}

	// Success!!
	return result, nil
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

/***********************************
 * MARSHAL / UNMARSHAL METHODS
 ***********************************/

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
