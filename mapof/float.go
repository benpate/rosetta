package mapof

import "github.com/benpate/rosetta/maps"

// Float is a map of string keys to float64 values with typed accessors.
type Float map[string]float64

// NewFloat returns a new, initialized Float map.
func NewFloat() Float {
	return make(Float)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x Float) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Float) Keys() []string {
	return maps.KeysSorted(x)
}

// Equal returns TRUE if this map has the same keys and values as the provided map.
func (x Float) Equal(value map[string]float64) bool {
	return maps.Equal(x, value)
}

// NotEqual returns TRUE if this map differs from the provided map.
func (x Float) NotEqual(value map[string]float64) bool {
	return maps.NotEqual(x, value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Float) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Float) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

// GetFloat returns the value for the key, or 0 if the key is not present.
func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

// GetFloatOK returns the value for the key and TRUE if it is present, or (0, false) if not.
func (x Float) GetFloatOK(key string) (float64, bool) {
	result, ok := x[key]
	return result, ok
}

// SetFloat stores a non-zero value at the key (deleting the key when the value is zero to keep the map sparse),
// allocating the map first if it is nil.
func (x *Float) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// Remove deletes the key from the map.
func (x *Float) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *Float) makeNotNil() {
	if *x == nil {
		*x = make(Float)
	}
}
