package mapof

import "github.com/benpate/rosetta/maps"

// Int64 is a map of string keys to int64 values with typed accessors.
type Int64 map[string]int64

// NewInt64 returns a new, initialized Int64 map.
func NewInt64() Int64 {
	return make(Int64)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x Int64) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Int64) Keys() []string {
	return maps.KeysSorted(x)
}

// IsMap returns TRUE, declaring this type a map for schema traversal. Implements schema.MapTyper.
func (x Int64) IsMap() bool {
	return true
}

// Equal returns TRUE if this map has the same keys and values as the provided map.
func (x Int64) Equal(value map[string]int64) bool {
	return maps.Equal(x, value)
}

// NotEqual returns TRUE if this map differs from the provided map.
func (x Int64) NotEqual(value map[string]int64) bool {
	return maps.NotEqual(x, value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Int64) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Int64) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetInt64 returns the value for the key, or 0 if the key is not present.
func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

// GetInt64OK returns the value for the key and TRUE if it is present, or (0, false) if not.
func (x Int64) GetInt64OK(key string) (int64, bool) {
	result, ok := x[key]
	return result, ok
}

// SetInt64 stores a non-zero value at the key (deleting the key when the value is zero to keep the map sparse),
// allocating the map first if it is nil.
func (x *Int64) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// Remove deletes the key from the map.
func (x *Int64) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *Int64) makeNotNil() {
	if *x == nil {
		*x = make(Int64)
	}
}
