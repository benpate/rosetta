package sliceof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type String []string

func NewString() String {
	return make(String, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

func (x String) Length() int {
	return len(x)
}

func (x String) IsLength(length int) bool {
	return len(x) == length
}

func (x String) IsEmpty() bool {
	return len(x) == 0
}

func (x String) First() string {
	if len(x) > 0 {
		return x[0]
	}
	return ""
}

func (x String) Last() string {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return ""
}

func (x String) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

func (x String) Contains(value string) bool {
	return slice.Contains(x, value)
}

func (x String) Equal(value []string) bool {
	return slice.Equal(x, value)
}

/****************************************
 * Getter Interfaces/Setters
 ****************************************/

func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetStringOK(key string) (string, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return "", false
}

func (s *String) SetString(key string, value string) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *String) SetValue(value any) error {
	*s = convert.SliceOfString(value)
	return nil
}

func (s *String) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
