package schema

import (
	"reflect"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func (element Object) getFromMap(object reflect.Value, path list.List) (reflect.Value, error) {

	const location = "schema.Object.getFromMap"

	// RULE: if the path is empty, then return the entire map
	if path.IsEmpty() {
		return object, nil
	}

	// Split the path into head and tail
	head, tail := path.Split()

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", path)
	}

	// Retrieve and return the existing value from the map
	mapKey := reflect.ValueOf(head)
	return property.Get(object.MapIndex(mapKey), tail)
}

func (element Object) setToMap(object reflect.Value, path list.List, value any) (reflect.Value, error) {

	const location = "schema.Object.setMap"

	// RULE: if the path is empty, then set and return the entire map
	if path.IsEmpty() {
		object.Set(reflect.ValueOf(value))
		return object, nil
	}

	// If the map is nil then initialize it as a new default value
	if object.IsNil() {
		object.Set(reflect.ValueOf(element.DefaultValue()))
	}

	head, tail := path.Split()

	// Try to find the matching property in this schema
	subElement, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", path, value)
	}

	// Retrieve the existing value from the map
	mapKey := reflect.ValueOf(head)     // map key
	mapValue := object.MapIndex(mapKey) // existing map value

	// Try to use the next schema subElement to get the correct, new value for the map
	mapValue, err := subElement.Set(mapValue, tail, value)

	if err != nil {
		return reflect.ValueOf(nil), derp.Wrap(err, location, "Failed to set value", path, value)
	}

	// Apply the new value back into the map.
	object.SetMapIndex(mapKey, mapValue)

	// Done
	return object, nil
}

func (element Object) removeFromMap(object reflect.Value, path list.List) (reflect.Value, error) {

	const location = "schema.Object.removeFromMap"

	// Split the path into head and tail
	head, tail := path.Split()

	// Try to find the matching property in this schema
	property, ok := element.Properties[head]

	if !ok {
		return reflect.ValueOf(nil), derp.NewInternalError(location, "Sub-element does not exist for this path", path)
	}

	// Retrieve the existing value from the map
	mapKey := reflect.ValueOf(head)     // map key
	mapValue := object.MapIndex(mapKey) // existing map value

	// If we're removing a sub-value, then pass this call to the sub-element.
	if !tail.IsEmpty() {
		mapValue, err := property.Remove(mapValue, tail)

		if err != nil {
			return reflect.ValueOf(nil), derp.Wrap(err, location, "Failed to remove value", path)
		}

		// Apply the new value back into the map.
		object.SetMapIndex(mapKey, mapValue)
		return object, nil

	}

	// Otherwise, just remove the map key.
	object.SetMapIndex(mapKey, reflect.Value{})
	return object, nil
}
