package mapof

import "github.com/benpate/rosetta/maps"

// Int is a map of string keys to int values with typed accessors.
type Int map[string]int

// NewInt returns a new, initialized Int map.
func NewInt() Int {
	return make(Int)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x Int) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Int) Keys() []string {
	return maps.KeysSorted(x)
}

// IsMap returns TRUE, declaring this type a map for schema traversal. Implements schema.MapTyper.
func (x Int) IsMap() bool {
	return true
}

// Equal returns TRUE if this map has the same keys and values as the provided map.
func (x Int) Equal(value map[string]int) bool {
	return maps.Equal(x, value)
}

// NotEqual returns TRUE if this map differs from the provided map.
func (x Int) NotEqual(value map[string]int) bool {
	return maps.NotEqual(x, value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Int) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Int) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetInt returns the value for the key, or 0 if the key is not present.
func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

// GetIntOK returns the value for the key and TRUE if it is present, or (0, false) if not.
func (x Int) GetIntOK(key string) (int, bool) {
	result, ok := x[key]
	return result, ok
}

// SetInt stores a non-zero value at the key (deleting the key when the value is zero to keep the map sparse),
// allocating the map first if it is nil.
func (x *Int) SetInt(key string, value int) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// Remove deletes the key from the map.
func (x *Int) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *Int) makeNotNil() {
	if *x == nil {
		*x = make(Int)
	}
}
