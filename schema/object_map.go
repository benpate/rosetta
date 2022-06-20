package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (element Object) setMap(object reflect.Value, path string, value any) error {

	const location = "schema.Object.setMap"

	var err error

	defer func() {
		if r := recover(); r != nil {
			err = derp.NewInternalError(location, "Panic in reflection", r)
		}
	}()

	if path == "" {
		return derp.NewInternalError(location, "Cannot set map value directly.  Set sub-items instead.", value)
	}

	// If the map is nil then initialize it as a new map[string]interface{}
	if object.IsNil() {
		var key string
		var value interface{}

		keyType := reflect.TypeOf(key)
		valueType := reflect.TypeOf(&value).Elem()

		mapType := reflect.MapOf(keyType, valueType) // TODO: can we be more specific than map[string]interface{}?
		object.Set(reflect.MakeMap(mapType))
	}

	head, tail := list.Split(path, ".")

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return derp.NewInternalError(location, "Sub-element does not exist for this path", path, value)
	}

	// Try to index the map
	keyValue := reflect.ValueOf(head)
	subValue := object.MapIndex(keyValue)

	// If the value already exists, then try to update it
	if subValue.CanSet() {
		if err = property.Set(subValue, tail, value); err != nil {
			return derp.Wrap(err, location, "Error setting sub-element", path, value)
		}
	}

	// Fall through means we're adding a new value ot the map
	newValue := reflect.New(property.Type()).Elem()

	if err := property.Set(newValue, tail, value); err != nil {
		return derp.Wrap(err, location, "Error setting sub-element", path, value)
	}

	object.SetMapIndex(keyValue, newValue)

	// Done
	return err
}
