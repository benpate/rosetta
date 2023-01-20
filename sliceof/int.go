package sliceof

import (
	"strconv"

	"github.com/benpate/rosetta/convert"
)

type Int []int

/****************************************
 * Accessors
 ****************************************/

func (x Int) Len() int {
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

/****************************************
 * Getters/Setters
 ****************************************/

func (x Int) GetIntOK(key string) (int, bool) {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index], true
		}
	}

	return 0, false
}

func (s *Int) SetIntOK(key string, value int) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}

func (s *Int) SetValue(value any) error {
	*s = convert.SliceOfInt(value)
	return nil
}
