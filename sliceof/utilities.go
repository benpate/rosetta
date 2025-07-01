package sliceof

import (
	"strconv"
)

func sliceIndex(index int, maximums ...int) (int, bool) {

	if index < 0 {
		return 0, false
	}

	for _, max := range maximums {
		if index >= max {
			return 0, false
		}
	}

	return index, true
}

func sliceStringIndex(key string, maximums ...int) (int, bool) {

	index, err := strconv.Atoi(key)

	if err != nil {
		return 0, false
	}

	return sliceIndex(index, maximums...)
}

func growSlice[T any, S ~[]T](value *S, length int) {

	for len(*value) <= length {
		var item T
		*value = append(*value, item)
	}
}
