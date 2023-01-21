package sliceof

import (
	"strconv"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/schema"
)

type Int64 []int64

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

func (x Int64) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

/****************************************
 * Getter Interfaces/Setters
 ****************************************/

func (x Int64) GetInt64OK(key string) (int64, bool) {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index], true
		}
	}

	return 0, false
}

func (s *Int64) SetInt64OK(key string, value int64) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}

func (s *Int64) SetValue(value any) error {
	*s = convert.SliceOfInt64(value)
	return nil
}

func (s *Int64) Remove(key string) bool {

	if index, ok := schema.Index(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
