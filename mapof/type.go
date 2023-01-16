package mapof

type Type[T any] map[string]T

/****************************************
 * Path Getters
 ****************************************/

func (x Type[T]) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Type[T]) GetAnyOK(key string) (any, bool) {
	result, ok := x[key]
	return result, ok
}

func (x Type[T]) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Type[T]) GetBoolOK(key string) (bool, bool) {
	if value, ok := any(x[key]).(bool); ok {
		return value, true
	}
	return false, false
}

func (x Type[T]) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Type[T]) GetFloatOK(key string) (float64, bool) {
	if value, ok := any(x[key]).(float64); ok {
		return value, true
	}
	return 0, false
}

func (x Type[T]) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Type[T]) GetIntOK(key string) (int, bool) {
	if value, ok := any(x[key]).(int); ok {
		return value, true
	}
	return 0, false
}

func (x Type[T]) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Type[T]) GetInt64OK(key string) (int64, bool) {
	if value, ok := any(x[key]).(int64); ok {
		return value, true
	}
	return 0, true
}

func (x Type[T]) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x Type[T]) GetStringOK(key string) (string, bool) {
	if value, ok := any(x[key]).(string); ok {
		return value, true
	}
	return "", false
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
