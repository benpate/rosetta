package slice

// At returns a bounds-safe value from a slice.  If the index is out of
// bounds for the slice, then `At` returns the zero value for that type.
func At[T any](slice []T, index int) T {
	result, _ := AtOK(slice, index)
	return result
}

// AtOK returns a bounds-safe value from a slice.  If the index is out of
// bounds for the slice, then `At` returns the zero value for that type.
func AtOK[T any](slice []T, index int) (T, bool) {

	// Underflow? Return empty value
	if index < 0 {
		var empty T
		return empty, false
	}

	// Overflow? Return empty value
	if index >= len(slice) {
		var empty T
		return empty, false
	}

	// Yay! Return the correct value
	return slice[index], true
}
