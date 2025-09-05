package convert

import "reflect"

// SliceOfInt converts the value into a slice of ints.
// It works with any, []any, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt(value any) []int {
	result, _ := SliceOfIntOk(value)
	return result
}

// SliceOfIntOk converts the value into a slice of ints.
// It works with float, int, int, string, and []any, []float, []int, []int, and []string values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfIntOk(value any) ([]int, bool) {

	// Nil check
	if value == nil {
		return make([]int, 0), false
	}

	// Known types
	switch value := value.(type) {

	case float64:
		item, ok := IntOk(value, 0)
		return []int{item}, ok

	case int:
		return []int{value}, true

	case int64:
		return []int{Int(value)}, true

	case string:
		item, ok := IntOk(value, 0)
		return []int{item}, ok

	case reflect.Value:
		return SliceOfInt(Interface(value)), true

	case Floater:
		return SliceOfInt(value.Float()), true

	case Inter:
		return SliceOfInt(value.Int()), true

	case Int64er:
		return SliceOfInt(value.Int64()), true

	case Stringer:
		item, ok := IntOk(value.String(), 0)
		return []int{item}, ok

	case []any:
		return sliceOfIntOk(value)

	case []float64:
		return sliceOfIntOk(value)

	case []int:
		return value, true

	case []int64:
		return sliceOfIntOk(value)

	case []string:
		return sliceOfIntOk(value)

	case []Floater:
		return sliceOfIntOk(value)

	case []Inter:
		return sliceOfIntOk(value)

	case []Int64er:
		return sliceOfIntOk(value)

	case []Stringer:
		return sliceOfIntOk(value)
	}

	// Use reflection to see if this is even an array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
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
