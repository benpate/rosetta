package convert

import "reflect"

// SliceOfFloat converts the value into a slice of floats.
// It works with any, []any, []float64, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfFloat(value any) []float64 {
	result, _ := SliceOfFloatOk(value)
	return result
}

// SliceOfFloat converts the value into a slice of floats.
// It works with any, []any, []float64, and float64 values.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfFloatOk(value any) ([]float64, bool) {

	// Nil check
	if value == nil {
		return make([]float64, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case float64:
		return []float64{typed}, true

	case int:
		return []float64{float64(typed)}, true

	case int64:
		return []float64{float64(typed)}, true

	case string:
		item, ok := FloatOk(typed, 0.0)
		return []float64{item}, ok

	case reflect.Value:
		return SliceOfFloat(Interface(typed)), true

	case Floater:
		return SliceOfFloat(typed.Float()), true

	case Inter:
		return SliceOfFloat(typed.Int()), true

	case Int64er:
		return SliceOfFloat(typed.Int64()), true

	case Stringer:
		item, ok := FloatOk(typed.String(), 0.0)
		return []float64{item}, ok

	case []any:
		return sliceOfFloatOk(typed)

	case []float64:
		return typed, true

	case []int:
		return sliceOfFloatOk(typed)

	case []int64:
		return sliceOfFloatOk(typed)

	case []string:
		return sliceOfFloatOk(typed)

	case []Floater:
		return sliceOfFloatOk(typed)

	case []Inter:
		return sliceOfFloatOk(typed)

	case []Int64er:
		return sliceOfFloatOk(typed)

	case []Stringer:
		return sliceOfFloatOk(typed)
	}

	// Use reflection to see if this is even an array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceOfFloatOk(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]float64, length)
		allOk := true
		for index := 0; index < length; index++ {
			item, ok := FloatOk(valueOf.Index(index), 0)
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk
	}

	return make([]float64, 0), false
}

// sliceOfFloatOk converts a slice of any type into a slice of float64s.
func sliceOfFloatOk[T any](value []T) ([]float64, bool) {
	result := make([]float64, len(value))
	allOk := true
	for index, v := range value {
		item, ok := FloatOk(v, 0)
		result[index] = item
		allOk = allOk && ok
	}
	return result, allOk
}
