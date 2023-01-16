package schema

import (
	"strconv"
)

// Index converts a string into an array index that is bounded
// by zero and the maximum value provided.  It returns the
// final index and a boolean that is TRUE if the index was
// converted successfully, and FALSE if it was truncated.
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
