package slice

// Contains scans a slice for a matching value, and
// returns TRUE if the value is found.
func Contains[T comparable](slice []T, value T) bool {

	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// NotContains returns TRUE if the value DOES NOT exist in the slice
func NotContains[T comparable](slice []T, value T) bool {
	return !Contains(slice, value)
}

// ContainsAny returns TRUE if the slice contains ANY of the provided values
func ContainsAny[T comparable](slice []T, values ...T) bool {
	for _, value := range values {
		if Contains(slice, value) {
			return true
		}
	}
	return false
}

// ContainsAll returns TRUE if the slice contains ALL of the provided values
func ContainsAll[T comparable](slice []T, values ...T) bool {
	for _, value := range values {
		if !Contains(slice, value) {
			return false
		}
	}
	return true
}
