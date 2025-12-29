package convert

import "reflect"

// SliceOfAny converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfAny(value any) []any {
	if result, _ := SliceOfAnyOk(value); result != nil {
		return result
	}

	return make([]any, 0)
}

// SliceOfAnyOk converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfAnyOk(value any) ([]any, bool) {

	// NILCHECK: value cannot be nil
	if value == nil {
		return make([]any, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case bool:
		return []any{typed}, true

	case float64:
		return []any{typed}, true

	case int:
		return []any{typed}, true

	case int64:
		return []any{typed}, true

	case string:
		return []any{typed}, true

	case reflect.Value:
		return SliceOfAnyOk(Interface(typed))

	case Floater:
		return []any{typed.Float()}, true

	case Hexer:
		return []any{typed.Hex()}, true

	case Inter:
		return []any{typed.Int()}, true

	case Int64er:
		return []any{typed.Int64()}, true

	case Stringer:
		return []any{typed.String()}, true

	case []any:
		return typed, true

	case []bool:
		return makeSliceOfAnyOk(typed)

	case []int:
		return makeSliceOfAnyOk(typed)

	case []int64:
		return makeSliceOfAnyOk(typed)

	case []float64:
		return makeSliceOfAnyOk(typed)

	case []string:
		return makeSliceOfAnyOk(typed)

	case []Floater:
		return makeSliceOfAnyOk(typed)

	case []Hexer:
		return makeSliceOfAnyOk(typed)

	case []Inter:
		return makeSliceOfAnyOk(typed)

	case []Int64er:
		return makeSliceOfAnyOk(typed)

	case []Stringer:
		return makeSliceOfAnyOk(typed)
	}

	// Use reflection to see if this is even aa array/slice
	switch valueOf := reflect.ValueOf(value); valueOf.Kind() {

	case reflect.Pointer:
		return SliceOfAny(valueOf.Elem().Interface()), true

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]any, length)
		for index := 0; index < length; index++ {
			result[index] = valueOf.Index(index).Interface()
		}
		return result, true
	}

	// Fall through means this isn't even an array/slice.  Admit defeat and go home.
	return make([]any, 0), false
}

// makeSliceOfAnyOk converts a slice of any type into a slice of int64s.
func makeSliceOfAnyOk[T any](value []T) ([]any, bool) {
	result := make([]any, len(value))
	for index, v := range value {
		result[index] = v
	}
	return result, true
}
