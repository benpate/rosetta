package compare

// Int compares two int values.
// It returns -1 if value1 is LESS THAN value2.
// It returns 0 if value1 is EQUAL TO value2.
// It returns 1 if value1 is GREATER THAN value2.
func Bool(value1 bool, value2 bool) int {

	if value1 {
		if value2 {
			return 0
		}
		return -1
	}

	if value2 {
		return 1
	}

	return 0
}
