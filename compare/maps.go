package compare

func Maps[T comparable](a, b map[string]T) bool {

	if len(a) != len(b) {
		return false
	}

	for key, value := range a {
		otherValue, exists := b[key]

		if !exists {
			return false
		}

		if otherValue != value {
			return false
		}
	}

	return true
}
