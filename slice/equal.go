package slice

// Equal returns true if the two slices are identical, having the same items in the same order, with no alterations.
func Equal[T comparable](value1 []T, value2 []T) bool {

	// Lengths must be identical
	if len(value1) != len(value2) {
		return false
	}

	// Items at each index must be identical
	for index := range value1 {
		if value1[index] != value2[index] {
			return false
		}
	}

	return true
}

// NotEqual returns TRUE if the two slices are NOT identical
func NotEqual[T comparable](value1 []T, value2 []T) bool {
	return !Equal(value1, value2)
}
