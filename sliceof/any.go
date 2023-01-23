package sliceof

import (
	"github.com/benpate/rosetta/convert"
)

type Any []any

func NewAny() Any {
	return make(Any, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

func (x Any) Length() int {
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

/******************************************
 * Getter Interfaces
 ******************************************/

func (x Any) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Any) GetAnyOK(key string) (any, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}
	return nil, false
}

func (x Any) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Any) GetBoolOK(key string) (bool, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.BoolOk(value, false)
	}
	return false, false
}

func (x Any) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Any) GetIntOK(key string) (int, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.IntOk(value, 0)
	}

	return 0, false
}

func (x Any) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Any) GetInt64OK(key string) (int64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.Int64Ok(value, 0)
	}
	return 0, false
}

func (x Any) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Any) GetFloatOK(key string) (float64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.FloatOk(value, 0)
	}
	return 0, false
}

func (x *Any) GetObject(key string) (any, bool) {
	if index, ok := sliceIndex(key); ok {
		growSlice(x, index)
		return &(*x)[index], true
	}

	return nil, false
}

func (x Any) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x Any) GetStringOK(key string) (string, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.StringOk(value, "")
	}
	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (x *Any) SetAny(key string, value any) bool {

	if index, ok := sliceIndex(key); ok {
		growSlice(x, index)
		(*x)[index] = value
		return true
	}

	return false
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

func (x *Any) SetValue(value any) error {
	*x = convert.SliceOfInterface(value)
	return nil
}

func (x *Any) Remove(key string) bool {

	if index, ok := sliceIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
