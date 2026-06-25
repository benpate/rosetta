package sliceof

import (
	"iter"
	"math/rand/v2"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
)

// Object is a slice of values of a single type T, with typed accessors and schema-traversal support.
type Object[T any] []T

// NewObject returns an Object slice containing the provided values (or an empty slice if none are given).
func NewObject[T any](values ...T) Object[T] {

	if len(values) == 0 {
		return make(Object[T], 0)
	}

	return values
}

/******************************************
 * Slice Manipulations
 ******************************************/

// Length returns the number of elements in the slice
func (x Object[T]) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Object[T]) IsLength(length int) bool {
	return len(x) == length
}

// IsZero returns TRUE if the slice contains no elements.
// This is an alias for IsEmpty, and implements the `Zeroer`
// interface used by many packages (including go/json)
func (x Object[T]) IsZero() bool {
	return len(x) == 0
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Object[T]) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Object[T]) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Object[T]) First() T {
	if len(x) > 0 {
		return x[0]
	}
	var result T
	return result
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Object[T]) FirstN(n int) Object[T] {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Object[T]) Last() T {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	var result T
	return result
}

// Find returns the first element in the slice that satisfies the provided function.
func (x Object[T]) Find(fn func(T) bool) (T, bool) {
	return slice.Find(x, fn)
}

// Filter returns all elements in the slice that satisfies the provided function.
func (x Object[T]) Filter(fn func(T) bool) Object[T] {
	return slice.Filter(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Object[T]) At(index int) T {
	return slice.At(x, index)
}

// AtOK returns the element at the index and TRUE, or the zero value and FALSE if out of range.
func (x Object[T]) AtOK(index int) (T, bool) {
	return slice.AtOK(x, index)
}

// Reverse returns a new slice with the elements in reverse order
func (x Object[T]) Reverse() Object[T] {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the "match" function returns TRUE for any element in the slice
func (x Object[T]) Contains(match func(T) bool) bool {
	for _, value := range x {
		if match(value) {
			return true
		}
	}

	return false
}

// Range returns an iterator that yields each value in this slice.
func (x Object[T]) Range() iter.Seq2[int, T] {
	return slice.Range(x)
}

// Append adds one or more elements to the end of the slice
func (x *Object[T]) Append(values ...T) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Object[T]) Shuffle() Object[T] {
	rand.Shuffle(len(x), func(i, j int) { // NOSONAR: does not need to be cryptographically secure.
		x[i], x[j] = x[j], x[i]
	})
	return x
}

// Keys returns a slice of strings representing the indexes of this slice
func (x Object[T]) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/******************************************
 * Getter Interfaces
 ******************************************/

// GetAny returns the value for the key index, or nil if absent.
func (x Object[T]) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

// GetAnyOK returns the value for the key index and TRUE if present.
func (x Object[T]) GetAnyOK(key string) (any, bool) {
	if index, ok := sliceStringIndex(key, x.Length()); ok {
		return x[index], true
	}

	switch key {

	case "last":
		return x.AtOK(x.Length() - 1)
	case "next":
		return x.AtOK(x.Length())
	}

	return nil, false
}

// GetPointer returns a pointer to the element at the key index, growing the slice as needed (implements schema PointerGetter).
func (x *Object[T]) GetPointer(name string) (any, bool) {

	// Get a valid index for the slice
	if index, ok := sliceStringIndex(name); ok {
		growSlice(x, index)

		// Return result
		return &(*x)[index], true
	}

	switch name {

	case "last":
		return x.GetPointer(strconv.Itoa(x.Length() - 1))
	case "next":
		return x.GetPointer(strconv.Itoa(x.Length()))
	}

	// Failure!!
	return nil, false
}

/******************************************
 * Setter Interfaces
 ******************************************/

// SetIndex stores the value at the index, growing the slice to fit if necessary.
func (s *Object[T]) SetIndex(index int, value any) bool {

	typed, ok := value.(T)

	if !ok {
		return false
	}

	growSlice(s, index)
	(*s)[index] = typed
	return true
}

// GetIndex returns the value at the specified index, and a boolean indicating success
func (x Object[T]) GetIndex(index int) (any, bool) {
	return slice.AtOK(x, index)
}

// SetValue replaces the entire slice with the provided value, if it can be converted.
func (s *Object[T]) SetValue(value any) error {

	switch typed := value.(type) {

	case Object[T]:
		*s = typed
		return nil

	case *Object[T]:
		*s = *typed
		return nil
	}

	return derp.Internal("sliceof.Object[T].SetValue", "Unable to convert value to Object[T]", value)
}

// Remove deletes the element identified by the key index.
func (x *Object[T]) Remove(key string) bool {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}

// RemoveAt deletes the element at the given index.
func (x *Object[T]) RemoveAt(index int) bool {

	if index, ok := sliceIndex(index, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
