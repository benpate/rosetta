package convert

import (
	"reflect"
	"strings"
)

// SliceOfInt64 converts the value into a slice of int64s.
// It works with any, []any, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt64(value any) []int64 {
	if result, _ := SliceOfInt64Ok(value); result != nil {
		return result
	}

	return make([]int64, 0)
}

// SliceOfInt64Ok converts the value into a slice of int64s.
// It works with float64, int, int64, string, and []any, []float64, []int, []int64, and []string values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfInt64Ok(value any) ([]int64, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make([]int64, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case float64:
		item, ok := Int64Ok(typed, 0)
		return []int64{item}, ok

	case int:
		return []int64{int64(typed)}, true

	case int64:
		return []int64{typed}, true

	case string:
		split := strings.Split(typed, ",")
		return sliceOfInt64Ok(split)

	case reflect.Value:
		return SliceOfInt64(Interface(typed)), true

	case Floater:
		return SliceOfInt64(typed.Float()), true

	case Inter:
		return SliceOfInt64(typed.Int()), true

	case Int64er:
		return SliceOfInt64(typed.Int64()), true

	case Stringer:
		item, ok := Int64Ok(typed.String(), 0)
		return []int64{item}, ok

	case []any:
		return sliceOfInt64Ok(typed)

	case []float64:
		return sliceOfInt64Ok(typed)

	case []int:
		return sliceOfInt64Ok(typed)

	case []int64:
		return typed, true

	case []string:
		return sliceOfInt64Ok(typed)

	case []Floater:
		return sliceOfInt64Ok(typed)

	case []Inter:
		return sliceOfInt64Ok(typed)

	case []Int64er:
		return sliceOfInt64Ok(typed)

	case []Stringer:
		return sliceOfInt64Ok(typed)
	}

	// Use reflection to see if this is even an array/slice
	switch valueOf := reflect.ValueOf(value); valueOf.Kind() {

	case reflect.Pointer:
		return SliceOfInt64Ok(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]int64, length)
		allOk := true
		for index := 0; index < length; index++ {
			item, ok := Int64Ok(valueOf.Index(index), 0)
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk
	}

	return make([]int64, 0), false
}

// sliceOfInt64Ok converts a slice of any type into a slice of int64s.
func sliceOfInt64Ok[T any](value []T) ([]int64, bool) {
	result := make([]int64, len(value))
	allOk := true
	for index, v := range value {
		item, ok := Int64Ok(v, 0)
		result[index] = item
		allOk = allOk && ok
	}
	return result, allOk
}
