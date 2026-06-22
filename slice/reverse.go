package slice

// Reverse reverses the order of the elements in the slice in place, and returns it.
func Reverse[T any](x []T) []T {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}
