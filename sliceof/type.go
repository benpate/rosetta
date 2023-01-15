package sliceof

import (
	"strconv"
)

type Type[T any] []T

/****************************************
 * Accessors
 ****************************************/

func (x Type[T]) Length() int {
	return len(x)
}

func (x Type[T]) IsLength(length int) bool {
	return len(x) == length
}

func (x Type[T]) IsEmpty() bool {
	return len(x) == 0
}

func (x Type[T]) First() T {
	if len(x) > 0 {
		return x[0]
	}
	var result T
	return result
}

func (x Type[T]) Last() T {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	var result T
	return result
}

func (x Type[T]) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

/****************************************
 * Path Getters
 ****************************************/

func (x Type[T]) GetAny(key string) any {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}
	return nil
}

func (x Type[T]) GetBool(key string) bool {
	if value, ok := x.GetAny(key).(bool); ok {
		return value
	}
	return false
}

func (x Type[T]) GetInt(key string) int {
	if value, ok := x.GetAny(key).(int); ok {
		return value
	}

	return 0
}

func (x Type[T]) GetInt64(key string) int64 {
	if value, ok := x.GetAny(key).(int64); ok {
		return value
	}
	return 0
}

func (x Type[T]) GetFloat(key string) float64 {
	if value, ok := x.GetAny(key).(float64); ok {
		return value
	}
	return 0
}

func (x Type[T]) GetString(key string) string {
	if value, ok := x.GetAny(key).(string); ok {
		return value
	}
	return ""
}

/****************************************
 * Path Getters
 ****************************************/

func (x *Type[T]) SetAny(key string, value any) bool {
	if typedValue, ok := value.(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) SetBool(key string, value bool) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) SetInt(key string, value int) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) SetInt64(key string, value int64) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) SetFloat(key string, value float64) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) SetString(key string, value string) bool {
	if typedValue, ok := any(value).(T); ok {
		return x.set(key, typedValue)
	}
	return false
}

func (x *Type[T]) set(key string, value T) bool {
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

func (x *Type[T]) GetChild(key string) (any, bool) {

	index, err := strconv.Atoi(key)

	if err != nil {
		return nil, false
	}

	if index < 0 {
		return nil, false
	}

	for index >= len(*x) {
		var newValue T
		*x = append(*x, newValue)
	}

	return &(*x)[index], true
}
