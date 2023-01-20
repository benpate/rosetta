package sliceof

import (
	"strconv"

	"github.com/benpate/rosetta/convert"
)

type String []string

/****************************************
 * Accessors
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

/****************************************
 * Getters/Setters
 ****************************************/

func (x String) GetStringOK(key string) (string, bool) {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index], true
		}
	}

	return "", false
}

func (s *String) SetStringOK(key string, value string) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}

func (s *String) SetValue(value any) error {
	*s = convert.SliceOfString(value)
	return nil
}
