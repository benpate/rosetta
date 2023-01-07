package sliceof

import "strconv"

type Int64 []int64

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
