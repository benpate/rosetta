package slice

// At returns a bounds-safe value from a slice.  If the index is out of
// bounds for the slice, then `At` returns the zero value for that type.
func At[T any](slice []T, index int) T {

	// Underflow? Return empty value
	if index < 0 {
		var empty T
		return empty
	}

	// Overflow? Return empty value
	if index >= len(slice) {
		var empty T
		return empty
	}

	// Yay! Return the correct value
	return slice[index]
}
