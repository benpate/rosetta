package compare

import (
	"strings"

	"github.com/benpate/rosetta/convert"
)

// Contains is a simple "generic-safe" function for string comparison.  It returns TRUE if value1 contains value2
func Contains(value1 any, value2 any) bool {

	switch value1 := value1.(type) {

	case string:

		if value2, ok := convert.StringOk(value2, ""); ok {
			return strings.Contains(value1, value2)
		}

	case []string:

		if value2, ok := convert.StringOk(value2, ""); ok {

			for index := range value1 {
				if value1[index] == value2 {
					return true
				}
			}
		}

	case []int:

		if value2, ok := convert.IntOk(value2, 0); ok {

			for index := range value1 {
				if value1[index] == value2 {
					return true
				}
			}
		}

	case []float64:

		if value2, ok := convert.FloatOk(value2, 0); ok {

			for index := range value1 {
				if value1[index] == value2 {
					return true
				}
			}
		}
	}

	return false
}
