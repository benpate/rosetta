package mapof

import (
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"

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

// Returns a SORTED slice of all keys in the map
func (x Object[T]) Keys() []string {
	return maps.KeysSorted(x)
}

// Returns TRUE if the map is empty
func (x Object[T]) IsEmpty() bool {
	return len(x) == 0
}

// Returns TRUE if the map is NOT empty
func (x Object[T]) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (object Object[T]) GetPointer(name string) (any, bool) {
	value, ok := object[name]
	return value, ok
}

func (object *Object[T]) SetObject(element schema.Element, path list.List, value any) error {

	if path.IsEmpty() {
		return derp.InternalError("mapof.Object.SetObject", "Cannot set values on empty path")
	}

	object.makeNotNil()

	head, tail := path.Split()

	if tail.IsEmpty() {
		if typed, ok := value.(T); ok {
			(*object)[head] = typed
			return nil
		}
		return derp.InternalError("mapof.Object.SetObject", "Invalid type", head, value)
	}

	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.InternalError("mapof.Object.SetObject", "Unknown property", head)
	}

	tempValue := (*object)[head]

	if err := schema.SetElement(&tempValue, subElement, tail, value); err != nil {
		return derp.Wrap(err, "mapof.Object.SetObject", "Error setting value", path)
	}

	// Reapply the updated value to the map
	(*object)[head] = tempValue

	return nil
}

func (object *Object[T]) Remove(key string) bool {
	object.makeNotNil()
	delete(*object, key)
	return true
}

func (object *Object[T]) makeNotNil() {
	if *object == nil {
		*object = make(Object[T])
	}
}

/******************************************
 * Other Methods
 ******************************************/

func (object Object[T]) IsZeroValue(name string) bool {
	return compare.IsZero(object[name])
}
