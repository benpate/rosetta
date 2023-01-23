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
