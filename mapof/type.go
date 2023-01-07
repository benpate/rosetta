package mapof

import "github.com/benpate/rosetta/convert"

type Type[T any] map[string]T

/****************************************
 * Path Getters
 ****************************************/

func (x Type[T]) GetAny(key string) any {
	return convert.Interface(x[key])
}

func (x Type[T]) GetBool(key string) bool {
	if value, ok := any(x[key]).(bool); ok {
		return value
	}
	return false
}

func (x Type[T]) GetFloat(key string) float64 {
	if value, ok := any(x[key]).(float64); ok {
		return value
	}
	return 0
}

func (x Type[T]) GetInt(key string) int {
	if value, ok := any(x[key]).(int); ok {
		return value
	}
	return 0
}

func (x Type[T]) GetInt64(key string) int64 {
	if value, ok := any(x[key]).(int64); ok {
		return value
	}
	return 0
}

func (x Type[T]) GetString(key string) string {
	if value, ok := any(x[key]).(string); ok {
		return value
	}
	return ""
}

/****************************************
 * Path Setters
 ****************************************/

func (x *Type[T]) SetAny(key string, value any) bool {

	if typedValue, ok := value.(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

func (x *Type[T]) SetBool(key string, value bool) bool {
	if typedValue, ok := any(value).(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

func (x *Type[T]) SetFloat(key string, value float64) bool {
	if typedValue, ok := any(value).(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

func (x *Type[T]) SetInt(key string, value int) bool {
	if typedValue, ok := any(value).(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

func (x *Type[T]) SetInt64(key string, value Int64) bool {
	if typedValue, ok := any(value).(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

func (x *Type[T]) SetString(key string, value string) bool {
	if typedValue, ok := any(value).(T); ok {
		(*x)[key] = typedValue
		return true
	}

	return false
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Type[T]) GetChild(key string) (any, bool) {
	result, ok := (*x)[key]
	return result, ok
}
