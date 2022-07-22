package slice

// Equal returns true if the two slices are identical, having the same items in the same order, with no alterations.
func Equal[T comparable](value1 []T, value2 []T) bool {

	if len(value1) != len(value2) {
		return false
	}

	for index := range value1 {
		if value1[index] != value2[index] {
			return false
		}
	}

	return true
}
