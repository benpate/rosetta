package sliceof

import "strconv"

type Any []any

/****************************************
 * Accessors
 ****************************************/

func (x Any) Len() int {
	return len(x)
}

func (x Any) IsLength(length int) bool {
	return len(x) == length
}

func (x Any) IsEmpty() bool {
	return len(x) == 0
}

func (x Any) First() any {
	if len(x) > 0 {
		return x[0]
	}
	return nil
}

func (x Any) Last() any {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return nil
}

func (x Any) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

/****************************************
 * Path Getters
 ****************************************/

func (x Any) GetAny(key string) any {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}
	return nil
}

func (x Any) GetBool(key string) bool {
	if value, ok := x.GetAny(key).(bool); ok {
		return value
	}
	return false
}

func (x Any) GetInt(key string) int {
	if value, ok := x.GetAny(key).(int); ok {
		return value
	}

	return 0
}

func (x Any) GetInt64(key string) int64 {
	if value, ok := x.GetAny(key).(int64); ok {
		return value
	}
	return 0
}

func (x Any) GetFloat(key string) float64 {
	if value, ok := x.GetAny(key).(float64); ok {
		return value
	}
	return 0
}

func (x Any) GetString(key string) string {
	if value, ok := x.GetAny(key).(string); ok {
		return value
	}
	return ""
}

/****************************************
 * Path Getters
 ****************************************/

func (x *Any) SetAny(key string, value any) bool {
	index, err := strconv.Atoi(key)

	if err != nil {
		return false
	}

	if index < 0 {
		return false
	}

	for index >= len(*x) {
		*x = append(*x, nil)
	}

	(*x)[index] = value
	return true
}

func (x *Any) SetBool(key string, value bool) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetInt(key string, value int) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetInt64(key string, value int64) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetFloat(key string, value float64) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetString(key string, value string) bool {
	return x.SetAny(key, value)
}

func (x *Any) GetChild(key string) (any, bool) {
	index, err := strconv.Atoi(key)

	if err != nil {
		return nil, false
	}

	if index < 0 {
		return nil, false
	}

	for index >= len(*x) {
		*x = append(*x, nil)
	}

	return &(*x)[index], true
}
