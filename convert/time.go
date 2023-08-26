package convert

import (
	"time"
)

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

		timeFormats := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02", time.RFC1123, time.RFC1123Z, time.RubyDate, time.UnixDate, time.RFC822, time.RFC822Z}

		for _, timeFormat := range timeFormats {
			if parsed, err := time.Parse(timeFormat, typed); err == nil {
				return parsed, true
			}
		}

	case int:
		return TimeOk(int64(typed), defaultValue)

	case int64:

		// Assume Seconds
		if typed < 10000000000 {
			return time.Unix(typed, 0).In(time.UTC), true
		}

		// Assume Miliseconds
		return time.UnixMilli(typed).In(time.UTC), true
	}

	return defaultValue, false
}
