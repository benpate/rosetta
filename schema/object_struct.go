package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (element Object) getFromStruct(object reflect.Value, path list.List) (reflect.Value, Element, error) {

	const location = "schema.Object.getFromStruct"

	// RULE: if the path is empty, then return the entire struct
	if path.IsEmpty() {
		return object, element, nil
	}

	// Split the path into head and tail
	head, tail := path.Split()

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), element, derp.NewInternalError(location, "Sub-element does not exist for this path", path, object)
	}

	// Retrieve and return the existing value from the struct
	field, err := findFieldByTag(object, head)

	if err != nil {
		return reflect.ValueOf(nil), element, err
	}

	return property.Get(field, tail)
}

func (element Object) setToStruct(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	const location = "schema.Object.setToStruct"

	// Allow direct setting of struct values.
	// TODO: This should probably be validated against the schema.
	if path.IsEmpty() {
		return object, nil
	}

	// Try to find the matching property in this schema
	head, tail := path.Split()
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", path, value)
	}

	// Try to find the matching struct field in the object
	field, err := findFieldByTag(object, head)

	if err != nil {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Cannot find struct field with this path", path)
	}

	// Try to put the value into the object
	result, err := property.Set(field, tail, value)

	if err != nil {
		return reflect.ValueOf(nil), derp.Wrap(err, location, "Error creating value of sub-element", path, value)
	}

	field.Set(result)

	// Done
	return object, nil
}

func (element Object) removeFromStruct(object reflect.Value, path list.List) (reflect.Value, error) {

	const location = "schema.Object.removeFromStruct"

	// Try to find the matching property in this schema
	head, tail := path.Split()
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", element, path)
	}

	// Try to find the field in this struct
	field, err := findFieldByTag(object, head)

	if err != nil {
		return reflect.ValueOf(nil), derp.Wrap(err, location, "Cannot find struct field with this path", path)
	}

	// If we're removing a sub-value, then pass this call to the sub-element.
	if !tail.IsEmpty() {
		subResult, err := property.Remove(field, tail)
		if err != nil {
			return reflect.ValueOf(nil), derp.Wrap(err, location, "Error removing sub-element", path, object.Interface())
		}
		field.Set(subResult)
		return object, nil
	}

	// Otherwise, set the whole property to the default value
	field.Set(reflect.ValueOf(property.DefaultValue()))
	return object, nil
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

	return reflect.Value{}, derp.NewInternalError(location, "Tag does not exist", tag, value.Interface())
}
