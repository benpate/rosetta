package sliceof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Int []int

func NewInt() Int {
	return make(Int, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

func (x Int) Length() int {
	return len(x)
}

func (x Int) IsLength(length int) bool {
	return len(x) == length
}

func (x Int) IsEmpty() bool {
	return len(x) == 0
}

func (x Int) First() int {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

func (x Int) Last() int {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

func (x Int) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

func (x Int) Contains(value int) bool {
	return slice.Contains(x, value)
}

func (x Int) Equal(value []int) bool {
	return slice.Equal(x, value)
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetIntOK(key string) (int, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Int) SetInt(key string, value int) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Int) SetValue(value any) error {
	*s = convert.SliceOfInt(value)
	return nil
}

func (s *Int) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
