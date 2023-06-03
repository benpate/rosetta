package convert

import "time"

func EpochDate(value any) int64 {
	result, _ := EpochDateOk(value)
	return result
}

func EpochDateOk(value any) (int64, bool) {

	switch typed := value.(type) {

	case int:
		return int64(typed), true

	case int64:
		return typed, true

	case time.Time:
		return typed.Unix(), true

	case string:
		if result, err := time.Parse(time.RFC3339, typed); err == nil {
			return result.Unix(), true
		}
	}

	return 0, false
}
