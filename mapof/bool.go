package mapof

import "github.com/benpate/rosetta/maps"

// Bool is a map of string keys to bool values with typed accessors.
type Bool map[string]bool

// NewBool returns a new, initialized Bool map.
func NewBool() Bool {
	return make(Bool)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x Bool) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Bool) Keys() []string {
	return maps.KeysSorted(x)
}

// Equal returns TRUE if this map has the same keys and values as the provided map.
func (x Bool) Equal(value Bool) bool {
	return maps.Equal(x, value)
}

// NotEqual returns TRUE if this map differs from the provided map.
func (x Bool) NotEqual(value Bool) bool {
	return maps.NotEqual(x, value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Bool) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Bool) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetBool returns the value for the key, or FALSE if the key is not present.
func (x Bool) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

// GetBoolOK returns the value for the key and TRUE if it is present, or (false, false) if not.
func (x Bool) GetBoolOK(key string) (bool, bool) {
	if result, ok := x[key]; ok {
		return result, true
	}
	return false, false
}

// SetBool stores the value at the key, allocating the map first if it is nil.
func (x *Bool) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

// Remove deletes the key from the map.
func (x *Bool) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *Bool) makeNotNil() {
	if *x == nil {
		*x = make(Bool)
	}
}
