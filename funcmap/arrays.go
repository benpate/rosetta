package funcmap

import "github.com/benpate/rosetta/compare"

func addArraysFuncs(target map[string]any) {

	// array creates an array out of the provided items
	target["array"] = func(values ...any) []any {
		return values
	}

	// seq creates a sequence of integers from 0 to count-1
	target["seq"] = func(count int) []int {
		result := make([]int, count)
		for i := 0; i < count; i++ {
			result[i] = i
		}
		return result
	}

	// first returns the first non-zero value in the array
	target["first"] = func(values ...any) any {
		for _, value := range values {
			if compare.NotZero(value) {
				return value
			}
		}
		return nil
	}

	// in returns true if the first value is found in the remaining values
	target["in"] = func(value any, values ...any) bool {
		for _, test := range values {
			if value == test {
				return true
			}
		}
		return false
	}

}
