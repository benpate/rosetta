package compare

// Bool compares two bool values, ordering FALSE before TRUE.
// It returns -1 if value1 is LESS THAN value2 (false < true).
// It returns 0 if value1 is EQUAL TO value2.
// It returns 1 if value1 is GREATER THAN value2 (true > false).
func Bool(value1 bool, value2 bool) int {

	if value1 == value2 {
		return 0
	}

	if value1 {
		return 1
	}

	return -1
}
