package sliceof

type Object[T any] []T

func NewObject[T any]() Object[T] {
	return make(Object[T], 0)
}

/******************************************
 * Slice Manipulations
 ******************************************/

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

func (x Object[T]) Reverse() Object[T] {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

func (x *Object[T]) Append(values ...T) {
	*x = append(*x, values...)
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
