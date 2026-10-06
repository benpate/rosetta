package mapof

import (
	"fmt"

	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/maps"

	"github.com/benpate/derp"
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

// IsMap returns TRUE, declaring this type a map for schema traversal. Implements schema.MapTyper.
func (x Object[T]) IsMap() bool {
	return true
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

// SetKey stores the value under the key when it is a T, and returns an error otherwise
// (implements the schema KeySetter interface).
func (object *Object[T]) SetKey(key string, value any) error {

	typed, ok := value.(T)

	if !ok {
		return derp.Internal("mapof.Object.SetKey", "Invalid type", key, fmt.Sprintf("%T", value))
	}

	object.makeNotNil()
	(*object)[key] = typed
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
