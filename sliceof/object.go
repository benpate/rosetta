package sliceof

import (
	"github.com/benpate/rosetta/schema"
)

type Object[T any] []T

/****************************************
 * Accessors
 ****************************************/

func (x Object[T]) Length() int {
	return len(x)
}

func (x Object[T]) IsLength(length int) bool {
	return len(x) == length
}

func (x Object[T]) IsEmpty() bool {
	return len(x) == 0
}

func (x Object[T]) First() T {
	if len(x) > 0 {
		return x[0]
	}
	var result T
	return result
}

func (x Object[T]) Last() T {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	var result T
	return result
}

func (x Object[T]) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Object[T]) GetObjectOK(name string) (any, bool) {

	// Get a valid index for the slice
	if index, ok := schema.Index(name); ok {

		// Grow if necessary
		for index >= len(*x) {
			var newValue T
			*x = append(*x, newValue)
		}

		// Return result
		return &(*x)[index], true
	}

	// Failure!!
	return nil, false
}

func (x *Object[T]) Remove(key string) bool {

	if index, ok := schema.Index(key, x.Length()); ok {

		// Remove the item
		*x = append((*x)[:index], (*x)[index+1:]...)

		// Success!
		return true
	}

	// Failure!!
	return false
}
