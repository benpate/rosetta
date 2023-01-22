package sliceof

import (
	"strconv"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/schema"
)

type Bool []bool

/****************************************
 * Accessors
 ****************************************/

func (x Bool) Length() int {
	return len(x)
}

func (x Bool) IsLength(length int) bool {
	return len(x) == length
}

func (x Bool) IsEmpty() bool {
	return len(x) == 0
}

func (x Bool) First() bool {
	if len(x) > 0 {
		return x[0]
	}
	return false
}

func (x Bool) Last() bool {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return false
}

func (x Bool) Reverse() {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}

/****************************************
 * Getter Interfaces/Setters
 ****************************************/

func (x Bool) GetBool(key string) (bool, bool) {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index], true
		}
	}

	return false, false
}

func (s *Bool) SetBool(key string, value bool) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}

func (s *Bool) SetValue(value any) error {
	*s = convert.SliceOfBool(value)
	return nil
}

func (s *Bool) Remove(key string) bool {

	if index, ok := schema.Index(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
