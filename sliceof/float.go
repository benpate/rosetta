package sliceof

import "strconv"

type Float []float64

func (x Float) GetFloat(key string) float64 {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(x)) {
			return x[index]
		}
	}

	return 0
}

func (s *Float) SetFloat(key string, value float64) bool {
	if index, err := strconv.Atoi(key); err == nil {
		if (index >= 0) && (index < len(*s)) {
			(*s)[index] = value
			return true
		}
	}

	return false
}
