package sliceof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Float []float64

func NewFloat() Float {
	return make(Float, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

func (x Float) Length() int {
	return len(x)
}

func (x Float) IsLength(length int) bool {
	return len(x) == length
}

func (x Float) IsEmpty() bool {
	return len(x) == 0
}

func (x Float) First() float64 {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

func (x Float) Last() float64 {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

func (x Float) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

func (x Float) Contains(value float64) bool {
	return slice.Contains(x, value)
}

func (x Float) Equal(value []float64) bool {
	return slice.Equal(x, value)
}

func (x *Float) Append(values ...float64) {
	*x = append(*x, values...)
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetFloatOK(key string) (float64, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Float) SetFloat(key string, value float64) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Float) SetValue(value any) error {
	*s = convert.SliceOfFloat(value)
	return nil
}

func (s *Float) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
