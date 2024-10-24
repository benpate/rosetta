package convert

import (
	"reflect"
)

// SliceOfString converts the value into a slice of strings.
// It works with any, []any, []string, and string values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfString(value any) []string {
	result, _ := SliceOfStringOk(value)
	return result
}

// SliceOfString converts the value into a slice of strings.
// It works with any, []any, []string, and string values.
// It returns TRUE if the value was converted successfullt, and FALSE otherwise.
func SliceOfStringOk(value any) ([]string, bool) {

	// Nil check
	if value == nil {
		return make([]string, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case bool:
		result, ok := StringOk(typed, "false")
		return []string{result}, ok

	case float64:
		result, ok := StringOk(typed, "")
		return []string{result}, ok

	case int:
		result, ok := StringOk(typed, "")
		return []string{result}, ok

	case int64:
		result, ok := StringOk(typed, "")
		return []string{result}, ok

	case string:
		return []string{typed}, true

	case reflect.Value:
		return SliceOfStringOk(Interface(typed))

	case Floater:
		return SliceOfStringOk(typed.Float())

	case Inter:
		return SliceOfStringOk(typed.Int())

	case Int64er:
		return SliceOfStringOk(typed.Int64())

	case Hexer:
		return []string{typed.Hex()}, true

	case Stringer:
		return []string{typed.String()}, true

	case []any:
		return sliceOfStringOk(typed)

	case []bool:
		return sliceOfStringOk(typed)

	case []float64:
		return sliceOfStringOk(typed)

	case []int:
		return sliceOfStringOk(typed)

	case []int64:
		return sliceOfStringOk(typed)

	case []string:
		return typed, true

	case []Floater:
		return sliceOfStringOk(typed)

	case []Inter:
		return sliceOfStringOk(typed)

	case []Int64er:
		return sliceOfStringOk(typed)

	case []Hexer:
		return sliceOfStringOk(typed)

	case []Stringer:
		return sliceOfStringOk(typed)

	}

	// Use reflection to see if this is even an array/slice
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceOfStringOk(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		length := valueOf.Len()
		result := make([]string, length)
		allOk := true
		for index := 0; index < length; index++ {
			item, ok := StringOk(valueOf.Index(index), "")
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk
	}

	// Fall through is failure.  This is a nothing
	return make([]string, 0), false
}

// sliceOfStringOk is a generic helper to convert known slices into a slice of strings.
func sliceOfStringOk[T any](value []T) ([]string, bool) {
	result := make([]string, len(value))
	allOk := true
	for index, item := range value {
		itemString, ok := StringOk(item, "")
		result[index] = itemString
		allOk = allOk && ok
	}
	return result, allOk
}
