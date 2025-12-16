package convert

import (
	"reflect"
	"strings"
)

// SliceOfInt converts the value into a slice of ints.
// It works with any, []any, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt(value any) []int {
	if result, _ := SliceOfIntOk(value); result != nil {
		return result
	}

	return make([]int, 0)
}

// SliceOfIntOk converts the value into a slice of ints.
// It works with float, int, int, string, and []any, []float, []int, []int, and []string values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfIntOk(value any) ([]int, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make([]int, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case float64:
		item, ok := IntOk(value, 0)
		return []int{item}, ok

	case int:
		return []int{typed}, true

	case int64:
		return []int{Int(typed)}, true

	case string:
		split := strings.Split(typed, ",")
		return sliceOfIntOk(split)

	case reflect.Value:
		return SliceOfInt(Interface(typed)), true

	case Floater:
		return SliceOfInt(typed.Float()), true

	case Inter:
		return SliceOfInt(typed.Int()), true

	case Int64er:
		return SliceOfInt(typed.Int64()), true

	case Stringer:
		item, ok := IntOk(typed.String(), 0)
		return []int{item}, ok

	case []any:
		return sliceOfIntOk(typed)

	case []float64:
		return sliceOfIntOk(typed)

	case []int:
		return typed, true

	case []int64:
		return sliceOfIntOk(typed)

	case []string:
		return sliceOfIntOk(typed)

	case []Floater:
		return sliceOfIntOk(typed)

	case []Inter:
		return sliceOfIntOk(typed)

	case []Int64er:
		return sliceOfIntOk(typed)

	case []Stringer:
		return sliceOfIntOk(typed)
	}

	// Use reflection to see if this is even an array/slice
	switch valueOf := reflect.ValueOf(value); valueOf.Kind() {

	case reflect.Pointer:
		return SliceOfIntOk(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]int, length)
		allOk := true
		for index := 0; index < length; index++ {
			item, ok := IntOk(valueOf.Index(index), 0)
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk
	}

	return make([]int, 0), false
}

// sliceOfIntOk converts a slice of any type into a slice of ints.
func sliceOfIntOk[T any](value []T) ([]int, bool) {
	result := make([]int, len(value))
	allOk := true
	for index, v := range value {
		item, ok := IntOk(v, 0)
		result[index] = item
		allOk = allOk && ok
	}
	return result, allOk
}
