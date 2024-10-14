package convert

import "reflect"

// SliceOfAny converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfAny(value any) []any {
	result, _ := SliceOfAnyOk(value)
	return result
}

// SliceOfAny converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfAnyOk(value any) ([]any, bool) {

	// Nil check
	if value == nil {
		return make([]any, 0), false
	}

	// Known types
	switch typed := value.(type) {

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

	case []int:
		return sliceOfAnyOk(typed)

	case []int64:
		return sliceOfAnyOk(typed)

	case []float64:
		return sliceOfAnyOk(typed)

	case []string:
		return sliceOfAnyOk(typed)

	case []Floater:
		return sliceOfAnyOk(typed)

	case []Hexer:
		return sliceOfAnyOk(typed)

	case []Inter:
		return sliceOfAnyOk(typed)

	case []Int64er:
		return sliceOfAnyOk(typed)

	case []Stringer:
		return sliceOfAnyOk(typed)
	}

	// Use reflection to see if this is even aa array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
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

// sliceOfAnyOk converts a slice of any type into a slice of int64s.
func sliceOfAnyOk[T any](value []T) ([]any, bool) {
	result := make([]any, len(value))
	for index, v := range value {
		result[index] = v
	}
	return result, true
}
