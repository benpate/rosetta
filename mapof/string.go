package mapof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/maps"
)

// String is a map of string keys to string values with typed accessors.
type String map[string]string

// NewString returns a new, initialized String map.
func NewString() String {
	return make(String)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x String) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x String) Keys() []string {
	return maps.KeysSorted(x)
}

// Equal returns TRUE if this map has the same keys and values as the provided map.
func (x String) Equal(value map[string]string) bool {
	// Lengths must be identical
	if len(x) != len(value) {
		return false
	}

	// Items at each index must be identical
	for key := range x {
		if x[key] != value[key] {
			return false
		}
	}

	return true
}

// NotEqual returns TRUE if this map differs from the provided map.
func (x String) NotEqual(value map[string]string) bool {
	return !x.Equal(value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x String) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x String) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetString returns the value for the key, or "" if the key is not present.
func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

// GetStringOK returns the value for the key and TRUE if it is present, or ("", false) if not.
func (x String) GetStringOK(key string) (string, bool) {
	result, ok := x[key]
	return result, ok
}

// SetString stores a non-empty value at the key (deleting the key when the value is empty to keep the map sparse),
// allocating the map first if it is nil.
func (x *String) SetString(key string, value string) bool {
	x.makeNotNil()

	if value == "" {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// Remove deletes the key from the map.
func (x *String) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *String) makeNotNil() {
	if *x == nil {
		*x = make(String)
	}
}

// MapOfAny returns the map's contents as a map[string]any.
func (x String) MapOfAny() map[string]any {
	return convert.MapOfAny(x.MapOfString())
}

// MapOfString returns the underlying map[string]string.
func (x String) MapOfString() map[string]string {
	return x
}
