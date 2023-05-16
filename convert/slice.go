package convert

import (
	"reflect"
	"strings"
)

// SliceOfString converts the value into a slice of strings.
// It works with any, []any, []string, and string values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfString(value any) []string {

	switch value := value.(type) {

	case string:
		return []string{value}

	case []string:
		return value

	case []int:
		result := make([]string, len(value))
		for index, v := range value {
			result[index] = String(v)
		}
		return result

	case []float64:
		result := make([]string, len(value))
		for index, v := range value {
			result[index] = String(v)
		}
		return result

	case []Stringer:
		result := make([]string, len(value))
		for index, v := range value {
			result[index] = v.String()
		}
		return result

	case []any:
		result := make([]string, len(value))
		for index, v := range value {
			result[index] = String(v)
		}
		return result

	case reflect.Value:
		return SliceOfString(Interface(value))
	}

	// Use reflection to see if this is even aa array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceOfString(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]string, length)
		for index := 0; index < length; index++ {
			result[index] = String(valueOf.Index(index))
		}
		return result
	}

	// Fall through is failure.  This is a nothing
	return make([]string, 0)
}

// SplitSliceOfString splits is a special case of SliceOfString.  If it receives a string, Stringer, or reflect.String,
// it will make a slice of strings by splitting the value into a slice.  All other values are passed to SliceOfString
// to be processed normally.
func SplitSliceOfString(value any, sep string) []string {

	if sep != "" {

		switch value := value.(type) {
		case string:
			return strings.Split(value, sep)
		case reflect.Value:
			return SplitSliceOfString(value.Interface(), sep)
		case Stringer:
			return strings.Split(value.String(), sep)
		}
	}

	return SliceOfString(value)
}

// SliceOfInt converts the value into a slice of ints.
// It works with any, []any, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt(value any) []int {

	switch value := value.(type) {

	case []int:
		return value

	case []any:
		result := make([]int, len(value))
		for index, v := range value {
			result[index] = Int(v)
		}
		return result

	case []string:
		result := make([]int, len(value))
		for index, v := range value {
			result[index] = Int(v)
		}
		return result

	case int:
		return []int{value}
	}

	return make([]int, 0)
}

// SliceOfInt64 converts the value into a slice of int64s.
// It works with any, []any, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt64(value any) []int64 {

	switch value := value.(type) {

	case []int64:
		return value

	case []int:
		result := make([]int64, len(value))
		for index, v := range value {
			result[index] = Int64(v)
		}
		return result

	case []any:
		result := make([]int64, len(value))
		for index, v := range value {
			result[index] = Int64(v)
		}
		return result

	case []string:
		result := make([]int64, len(value))
		for index, v := range value {
			result[index] = Int64(v)
		}
		return result

	case int:
		return []int64{int64(value)}

	case int64:
		return []int64{value}
	}

	return make([]int64, 0)
}

// SliceOfFloat converts the value into a slice of floats.
// It works with any, []any, []float64, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfFloat(value any) []float64 {

	switch value := value.(type) {

	case []any:
		result := make([]float64, len(value))
		for index, v := range value {
			result[index] = Float(v)
		}
		return result

	case []float64:
		return value

	case float64:
		return []float64{value}
	}

	return make([]float64, 0)
}

// SliceOfMap converts the value into a slice of map[string]any.
// It works with []any, []map[string]any.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfMap(value any) []map[string]any {

	switch value := value.(type) {

	case []map[string]any:
		return value

	case []any:
		result := make([]map[string]any, len(value))
		for index, v := range value {
			result[index] = MapOfAny(v)
		}
		return result
	}

	return make([]map[string]any, 0)
}

// SliceOfBool converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfBool(value any) []bool {

	switch value := value.(type) {

	case []bool:
		return value

	case []any:
		result := make([]bool, len(value))
		for index, v := range value {
			result[index] = Bool(v)
		}
		return result

	case []string:
		result := make([]bool, len(value))
		for index, v := range value {
			result[index] = Bool(v)
		}
		return result

	case []int:
		result := make([]bool, len(value))
		for index, v := range value {
			result[index] = Bool(v)
		}
		return result

	case []float64:
		result := make([]bool, len(value))
		for index, v := range value {
			result[index] = Bool(v)
		}
		return result

	case string, int, float64:
		return []bool{Bool(value)}
	}

	return make([]bool, 0)
}

// SliceOfAny converts the value into a slice of any.
// It works with any, []any, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfAny(value any) []any {

	if value == nil {
		return make([]any, 0)
	}

	switch value := value.(type) {

	case []any:
		return value

	case []bool:
		result := make([]any, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []int:
		result := make([]any, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []int64:
		result := make([]any, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []float64:
		result := make([]any, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []string:
		result := make([]any, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case bool, int, int64, float64, string:
		return []any{value}
	}

	// Use reflection to see if this is even aa array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceOfAny(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]any, length)
		for index := 0; index < length; index++ {
			result[index] = valueOf.Index(index).Interface()
		}
		return result
	}

	// Fall through means this isn't even an array/slice.  Admit defeat and go home.
	return make([]any, 0)
}

func SliceLength(value any) int {

	if value == nil {
		return 0
	}

	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceLength(valueOf.Elem().Interface())
	case reflect.Array, reflect.Slice:
		return valueOf.Len()
	}

	return 0
}
