package schema

// validate_Any checks that the provided value meets the requirements of the schema element.
// For "Any" types, there are no requirements, so this function always returns the original value.
func validate_Any[T any](element Any, value T) (T, bool, error) {
	return value, false, nil
}
