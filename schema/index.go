package schema

import (
	"strconv"
)

// Index converts a string into an array index, returning TRUE if it is a non-negative integer below
// every maximum provided.  Otherwise it returns FALSE, with 0 or one less than the first maximum reached.
func Index(token string, maxLengths ...int) (int, bool) {

	result, err := strconv.Atoi(token)

	if err != nil {
		return 0, false
	}

	if result < 0 {
		return 0, false
	}

	for _, maxLength := range maxLengths {
		if result >= maxLength {
			return maxLength - 1, false
		}
	}

	return result, true
}
