package mapof

import (
	"github.com/benpate/rosetta/list"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

type Object[T any] map[string]T

func NewObject[T any]() Object[T] {
	return make(Object[T])
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Object[T]) Keys() []string {
	keys := make([]string, 0, len(x))
	for key := range x {
		keys = append(keys, key)
	}
	return keys
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (object Object[T]) GetObject(name string) (any, bool) {
	value, ok := object[name]
	return value, ok
}

func (object *Object[T]) SetObject(element schema.Element, path list.List, value any) error {

	if path.IsEmpty() {
		return derp.NewInternalError("mapof.Object.SetObject", "Cannot set values on empty path")
	}

	object.makeNotNil()

	head, tail := path.Split()

	if tail.IsEmpty() {
		if typed, ok := value.(T); ok {
			(*object)[head] = typed
			return nil
		}
		return derp.NewInternalError("mapof.Object.SetObject", "Invalid type", head, value)
	}

	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.NewInternalError("mapof.Object.SetObject", "Unknown property", head)
	}

	tempValue := (*object)[head]

	if err := schema.SetElement(&tempValue, subElement, tail, value); err != nil {
		return derp.Wrap(err, "mapof.Object.SetObject", "Error setting value", path)
	}

	// Reapply the updated value to the map
	(*object)[head] = tempValue

	return nil
}

func (x *Object[T]) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Object[T]) makeNotNil() {
	if *x == nil {
		*x = make(Object[T])
	}
}
