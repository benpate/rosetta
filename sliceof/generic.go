package sliceof

import "strconv"

type Generic[T any] []T

/****************************************
 * Path Getters
 ****************************************/

func (x Generic[T]) GetAny(key string) any {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}
	return nil
}

func (x Generic[T]) GetBool(key string) bool {
	if value, ok := x.GetAny(key).(bool); ok {
		return value
	}
	return false
}

func (x Generic[T]) GetInt(key string) int {
	if value, ok := x.GetAny(key).(int); ok {
		return value
	}

	return 0
}

func (x Generic[T]) GetInt64(key string) int64 {
	if value, ok := x.GetAny(key).(int64); ok {
		return value
	}
	return 0
}

func (x Generic[T]) GetFloat(key string) float64 {
	if value, ok := x.GetAny(key).(float64); ok {
		return value
	}
	return 0
}

func (x Generic[T]) GetString(key string) string {
	if value, ok := x.GetAny(key).(string); ok {
		return value
	}
	return ""
}

/****************************************
 * Path Getters
 ****************************************/

func (x *Generic[T]) SetAny(key string, value any) bool {
	if typedValue, ok := value.(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) SetBool(key string, value bool) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) SetInt(key string, value int) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) SetInt64(key string, value int64) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) SetFloat(key string, value float64) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) SetString(key string, value string) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Generic[T]) set(key string, value T) bool {
	index, err := strconv.Atoi(key)

	if err != nil {
		return false
	}

	if index < 0 {
		return false
	}

	for index >= len(*x) {
		var newValue T
		*x = append(*x, newValue)
	}

	(*x)[index] = value
	return true
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Generic[T]) GetChild(key string) (any, bool) {
	index, err := strconv.Atoi(key)

	if err != nil {
		return nil, false
	}

	if index < 0 {
		return nil, false
	}

	for index >= len(*x) {
		return nil, false
	}

	return &(*x)[index], true
}
