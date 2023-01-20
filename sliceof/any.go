package sliceof

import (
	"strconv"

	"github.com/benpate/rosetta/convert"
)

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

func (x Any) GetAnyOK(key string) (any, bool) {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index], true
		}
	}
	return nil, false
}

func (x Any) GetBoolOK(key string) (bool, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		if typed, ok := value.(bool); ok {
			return typed, true
		}
	}
	return false, false
}

func (x Any) GetIntOK(key string) (int, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		if typed, ok := value.(int); ok {
			return typed, true
		}
	}

	return 0, false
}

func (x Any) GetInt64OK(key string) (int64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		if typed, ok := value.(int64); ok {
			return typed, true
		}
	}
	return 0, false
}

func (x Any) GetFloatOK(key string) (float64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		if typed, ok := value.(float64); ok {
			return typed, true
		}
	}
	return 0, false
}

func (x Any) GetStringOK(key string) (string, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		if typed, ok := value.(string); ok {
			return typed, true
		}
	}
	return "", false
}

/****************************************
 * Path Getters
 ****************************************/

func (x *Any) SetAnyOK(key string, value any) (bool, bool) {
	index, err := strconv.Atoi(key)

	if err != nil {
		return false, false
	}

	if index < 0 {
		return false, false
	}

	for index >= len(*x) {
		*x = append(*x, nil)
	}

	(*x)[index] = value
	return true, true
}

func (x *Any) SetBoolOK(key string, value bool) (bool, bool) {
	return x.SetAnyOK(key, value)
}

func (x *Any) SetIntOK(key string, value int) (bool, bool) {
	return x.SetAnyOK(key, value)
}

func (x *Any) SetInt64OK(key string, value int64) (bool, bool) {
	return x.SetAnyOK(key, value)
}

func (x *Any) SetFloatOK(key string, value float64) (bool, bool) {
	return x.SetAnyOK(key, value)
}

func (x *Any) GetObjectOK(key string) (any, bool) {
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

func (x *Any) SetStringOK(key string, value string) (bool, bool) {
	return x.SetAnyOK(key, value)
}

func (x *Any) SetValue(value any) error {
	*x = convert.SliceOfInterface(value)
	return nil
}
