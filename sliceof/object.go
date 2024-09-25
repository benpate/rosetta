package sliceof

import (
	"math/rand"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
)

type Object[T any] []T

func NewObject[T any]() Object[T] {
	return make(Object[T], 0)
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

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Object[T]) At(index int) T {
	return slice.At(x, index)
}

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

// Append adds one or more elements to the end of the slice
func (x *Object[T]) Append(values ...T) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Object[T]) Shuffle() Object[T] {
	rand.Shuffle(len(x), func(i, j int) {
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

func (x Object[T]) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Object[T]) GetAnyOK(key string) (any, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
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

func (x *Object[T]) GetPointer(name string) (any, bool) {

	// Get a valid index for the slice
	if index, ok := sliceIndex(name); ok {
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

func (s *Object[T]) SetValue(value any) error {

	if typed, ok := value.(Object[T]); ok {
		*s = typed
		return nil
	}

	return derp.NewInternalError("sliceof.Object[T].SetValue", "Unable to convert value to Object[T]", value)
}

func (x *Object[T]) Remove(key string) bool {

	if index, ok := sliceIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
