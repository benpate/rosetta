package convert

// SliceOfMap converts the value into a slice of map[string]any.
// It works with []any, []map[string]any.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfMap(value any) []map[string]any {
	result, _ := SliceOfMapOk(value)
	return result
}

// SliceOfMapOk converts the value into a slice of map[string]any.
// It works with []any, []map[string]any, []map[string]string, and []MapOfAnyGetter.
// It returns TRUE if the conversion was successful, and FALSE otherwise.
func SliceOfMapOk(value any) ([]map[string]any, bool) {

	// Nil check
	if value == nil {
		return make([]map[string]any, 0), false
	}

	// Known types
	switch typed := value.(type) {

	case []any:
		result := make([]map[string]any, len(typed))
		allOk := true
		for index, v := range typed {
			item, ok := MapOfAnyOk(v)
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk

	case []map[string]any:
		return typed, true

	case []map[string]string:
		result := make([]map[string]any, len(typed))
		allOk := true
		for index, v := range typed {
			item, ok := MapOfAnyOk(v)
			result[index] = item
			allOk = allOk && ok
		}
		return result, allOk

	case []MapOfAnyGetter:
		result := make([]map[string]any, len(typed))
		for index, v := range typed {
			result[index] = v.MapOfAny()
		}
		return result, true
	}

	// Fall through means the conversion was unsuccessful
	return make([]map[string]any, 0), false
}
