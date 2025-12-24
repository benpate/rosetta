// Package first provides utility functions to return the first non-empty value from a list of values
package first

// String returns the first non-empty value in the argument list
func String(values ...string) string {
	for index := range values {
		if values[index] != "" {
			return values[index]
		}
	}
	return ""
}

// Int returns the first non-zero value in the argument list
func Int(values ...int) int {
	for index := range values {
		if values[index] != 0 {
			return values[index]
		}
	}
	return 0
}

// Int64 returns the first non-zero value in the argument list
func Int64(values ...int64) int64 {
	for index := range values {
		if values[index] != 0 {
			return values[index]
		}
	}
	return 0
}
