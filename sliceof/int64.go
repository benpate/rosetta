package sliceof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Int64 []int64

func NewInt64() Int64 {
	return make(Int64, 0)
}

/****************************************
 * Accessors
 ****************************************/

func (x Int64) Length() int {
	return len(x)
}

func (x Int64) IsLength(length int) bool {
	return len(x) == length
}

func (x Int64) IsEmpty() bool {
	return len(x) == 0
}

func (x Int64) NotEmpty() bool {
	return len(x) > 0
}

func (x Int64) First() int64 {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

func (x Int64) Last() int64 {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

func (x Int64) Reverse() Int64 {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

func (x Int64) Contains(value int64) bool {
	return slice.Contains(x, value)
}

func (x Int64) ContainsAny(values ...int64) bool {
	return slice.ContainsAny(x, values...)
}

func (x Int64) ContainsAll(values ...int64) bool {
	return slice.ContainsAll(x, values...)
}

func (x Int64) Equal(value []int64) bool {
	return slice.Equal(x, value)
}

func (x *Int64) Append(values ...int64) {
	*x = append(*x, values...)
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Int64) GetInt64OK(key string) (int64, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Int64) SetInt64(key string, value int64) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Int64) SetValue(value any) error {
	*s = convert.SliceOfInt64(value)
	return nil
}

func (s *Int64) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
