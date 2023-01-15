package sliceof

import "strconv"

type Bool []bool

/****************************************
 * Accessors
 ****************************************/

func (x Bool) Len() int {
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
 * Getters/Setters
 ****************************************/

func (x Bool) GetBool(key string) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}

	return false
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
