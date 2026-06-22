package mapof

import (
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

// Object is a map of string keys to values of a single type T, with schema-traversal support.
type Object[T any] map[string]T

// NewObject returns a new, initialized Object map.
func NewObject[T any]() Object[T] {
	return make(Object[T])
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x Object[T]) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Object[T]) Keys() []string {
	return maps.KeysSorted(x)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Object[T]) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Object[T]) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetPointer returns the value for the key (implements the schema PointerGetter interface).
func (object Object[T]) GetPointer(name string) (any, bool) {
	value, ok := object[name]
	return value, ok
}

// SetObject descends the path (creating child entries as needed) and sets the value (implements the schema ObjectSetter interface).
func (object *Object[T]) SetObject(element schema.Element, path list.List, value any) error {

	if path.IsEmpty() {
		return derp.Internal("mapof.Object.SetObject", "Cannot set values on empty path")
	}

	object.makeNotNil()

	head, tail := path.Split()

	if tail.IsEmpty() {
		if typed, ok := value.(T); ok {
			(*object)[head] = typed
			return nil
		}
		return derp.Internal("mapof.Object.SetObject", "Invalid type", head, value)
	}

	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.Internal("mapof.Object.SetObject", "Unknown property", head)
	}

	tempValue := (*object)[head]

	if err := schema.SetProperty(subElement, &tempValue, tail.String(), value); err != nil {
		return derp.Wrap(err, "mapof.Object.SetObject", "Unable to set value", path)
	}

	// Reapply the updated value to the map
	(*object)[head] = tempValue

	return nil
}

// Remove deletes the key from the map.
func (object *Object[T]) Remove(key string) bool {
	object.makeNotNil()
	delete(*object, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (object *Object[T]) makeNotNil() {
	if *object == nil {
		*object = make(Object[T])
	}
}

/******************************************
 * Other Methods
 ******************************************/

// IsZeroValue returns TRUE if the named property is absent or holds a zero value.
func (object Object[T]) IsZeroValue(name string) bool {
	return compare.IsZero(object[name])
}
