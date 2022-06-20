package convert

import (
	"reflect"
)

// SliceOfString converts the value into a slice of strings.
// It works with interface{}, []interface{}, []string, and string values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfString(value interface{}) []string {

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

	case []interface{}:
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

// SliceOfInt converts the value into a slice of ints.
// It works with interface{}, []interface{}, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt(value interface{}) []int {

	switch value := value.(type) {

	case []int:
		return value

	case []interface{}:
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

// SliceOfFloat converts the value into a slice of floats.
// It works with interface{}, []interface{}, []float64, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfFloat(value interface{}) []float64 {

	switch value := value.(type) {

	case []interface{}:
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

// SliceOfMap converts the value into a slice of map[string]interface{}.
// It works with []interface{}, []map[string]interface{}.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfMap(value interface{}) []map[string]interface{} {

	switch value := value.(type) {

	case []map[string]interface{}:
		return value

	case []interface{}:
		result := make([]map[string]interface{}, len(value))
		for index, v := range value {
			result[index] = MapOfInterface(v)
		}
		return result
	}

	return make([]map[string]interface{}, 0)
}

// SliceOfBool converts the value into a slice of interface{}.
// It works with interface{}, []interface{}, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfBool(value interface{}) []bool {

	switch value := value.(type) {

	case []bool:
		return value

	case []interface{}:
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

// SliceOfInterface converts the value into a slice of interface{}.
// It works with interface{}, []interface{}, []string, []int, []float64, string, int, and float64 values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInterface(value interface{}) []interface{} {

	switch value := value.(type) {

	case []interface{}:
		return value

	case []string:
		result := make([]interface{}, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []int:
		result := make([]interface{}, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case []float64:
		result := make([]interface{}, len(value))
		for index, v := range value {
			result[index] = v
		}
		return result

	case string, int, float64:
		return []interface{}{value}
	}

	// Use reflection to see if this is even aa array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceOfInterface(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]interface{}, length)
		for index := 0; index < length; index++ {
			result[index] = valueOf.Index(index).Interface()
		}
		return result
	}

	// Fall through means this isn't even an array/slice.  Admit defeat and go home.
	return make([]interface{}, 0)
}
