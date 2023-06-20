package convert

import "time"

func Time(value any) time.Time {
	result, _ := TimeOk(value, time.Time{})
	return result
}

func TimeDefault(value any, defaultValue time.Time) time.Time {
	result, _ := TimeOk(value, defaultValue)
	return result
}

func TimeOk(value any, defaultValue time.Time) (time.Time, bool) {

	switch typed := value.(type) {

	case time.Time:
		return typed, true

	case string:

		if parsed, err := time.Parse(time.RFC3339, typed); err == nil {
			return parsed, true
		}

	case int:
		return TimeOk(int64(typed), defaultValue)

	case int64:

		// Assume Seconds
		if typed < 10000000000 {
			return time.Unix(typed, 0), true
		}

		// Assume Miliseconds
		return time.UnixMilli(typed), true
	}

	return defaultValue, false
}
