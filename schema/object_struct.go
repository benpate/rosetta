package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (element Object) setStruct(object reflect.Value, path string, value any) error {

	const location = "schema.Object.setStruct"

	var err error

	defer func() {
		if r := recover(); r != nil {
			err = derp.NewInternalError(location, "Panic in reflection", r)
		}
	}()

	if path == "" {
		return derp.NewInternalError(location, "Cannot set struct value directly.  Set sub-items instead.", value)
	}

	head, tail := list.Split(path, ".")

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return derp.NewInternalError(location, "Sub-element does not exist for this path", path, value)
	}

	// Try to find the matching struct field in the object
	field, err := findFieldByTag(object, head)

	if err != nil {
		return derp.NewInternalError(location, "Invalid struct tag", path)
	}

	// Try to put the value into the object
	if err := property.Set(field, tail, value); err != nil {
		return derp.Wrap(err, location, "Error creating value of sub-element", path, value)
	}

	// Done
	return err
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
