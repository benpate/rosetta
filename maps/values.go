package maps

import "cmp"

// Values returns a slice of all values of a map
func Values[T cmp.Ordered, U any](value map[T]U) []U {
	values := make([]U, 0, len(value))
	for _, v := range value {
		values = append(values, v)
	}

	return values
}
