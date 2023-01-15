package sliceof

import "strconv"

type Int64 []int64

/****************************************
 * Accessors
 ****************************************/

func (x Int64) Len() int {
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
 * Getters/Setters
 ****************************************/

func (x Int64) GetInt64(key string) int64 {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}

	return 0
}

func (s *Int64) SetInt64(key string, value int64) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}
