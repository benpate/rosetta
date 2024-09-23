package mapof

func MapOfAny(value any) (Any, bool) {

	switch typed := value.(type) {

	case Any:
		return typed, true

	case *Any:
		return *typed, true

	case Bool:
		result := make(Any, len(typed))

		for key, value := range typed {
			result[key] = value
		}

		return result, true

	case Float:
		result := make(Any, len(typed))

		for key, value := range typed {
			result[key] = value
		}

		return result, true

	case Int:
		result := make(Any, len(typed))

		for key, value := range typed {
			result[key] = value
		}

		return result, true

	case Int64:
		result := make(Any, len(typed))

		for key, value := range typed {
			result[key] = value
		}

		return result, true

	case String:
		result := make(Any, len(typed))

		for key, value := range typed {
			result[key] = value
		}

		return result, true

	case map[string]any:
		return Any(typed), true

	case map[string]bool:
		return MapOfAny(Bool(typed))

	case map[string]float64:
		return MapOfAny(Float(typed))

	case map[string]int:
		return MapOfAny(Int(typed))

	case map[string]int64:
		return MapOfAny(Int64(typed))

	case map[string]string:
		return MapOfAny(String(typed))

	case *map[string]any:
		return Any(*typed), true

	case *map[string]bool:
		return MapOfAny(Bool(*typed))

	case *map[string]float64:
		return MapOfAny(Float(*typed))

	case *map[string]int:
		return MapOfAny(Int(*typed))

	case *map[string]int64:
		return MapOfAny(Int64(*typed))

	case *map[string]string:
		return MapOfAny(String(*typed))

	}

	return nil, false
}
