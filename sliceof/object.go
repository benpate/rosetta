package sliceof

import "math/rand"

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

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Object[T]) GetPointer(name string) (any, bool) {

	// Get a valid index for the slice
	if index, ok := sliceIndex(name); ok {
		growSlice(x, index)

		// Return result
		return &(*x)[index], true
	}

	// Failure!!
	return nil, false
}

func (x *Object[T]) Remove(key string) bool {

	if index, ok := sliceIndex(key, x.Length()); ok {

		// Remove the item
		*x = append((*x)[:index], (*x)[index+1:]...)

		// Success!
		return true
	}

	// Failure!!
	return false
}
