package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (element Object) getFromStruct(object reflect.Value, path string) (reflect.Value, Element, error) {

	// RULE: if the path is empty, then return the entire struct
	if path == "" {
		return object, element, nil
	}

	// Split the path into head and tail
	head, tail := list.Dot(path).Split()

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), element, derp.NewInternalError("schema.Object.getFromStruct", "Sub-element does not exist for this path", path, object)
	}

	// Retrieve and return the existing value from the struct
	field, err := findFieldByTag(object, head)

	if err != nil {
		return reflect.ValueOf(nil), element, err
	}

	return property.Get(field, tail.String())
}

func (element Object) setToStruct(object reflect.Value, path string, value any) (reflect.Value, error) {

	const location = "schema.Object.setStruct"

	var err error

	if path == "" {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Cannot set struct value directly.  Set sub-items instead.", value)
	}

	head, tail := list.Dot(path).Split()

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", path, value)
	}

	// Try to find the matching struct field in the object
	field, err := findFieldByTag(object, head)

	if err != nil {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Invalid struct tag", path)
	}

	// Try to put the value into the object
	result, err := property.Set(field, tail.String(), value)

	if err != nil {
		return reflect.ValueOf(nil), derp.Wrap(err, location, "Error creating value of sub-element", path, value)
	}

	field.Set(result)

	// Done
	return object, err
}

// findFieldByTag returns the field whose "path" tag matches the provided value.
func findFieldByTag(value reflect.Value, tag string) (reflect.Value, error) {

	const location = "schema.findFieldByTag"

	if value.Kind() != reflect.Struct {
		return reflect.Value{}, derp.NewInternalError(location, "Value must be a struct")
	}

	count := value.NumField()
	typeOf := value.Type()

	for index := 0; index < count; index++ {
		if typeOf.Field(index).Tag.Get("path") == tag {
			return value.Field(index), nil
		}
	}

	return reflect.Value{}, derp.NewInternalError(location, "Tag does not exist", tag)
}
