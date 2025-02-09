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

	case Timer:
		return typed.Time(), true

	case ToTimer:
		return typed.ToTime(), true

	case string:

		if result, ok := TimeWithLocale(typed); ok {
			return result, true
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

// TimeWithLocale parses a string into a time.Time using the provided locale(s).
// If no locale is provided, it will use a list of common layouts, including RFE3339, RFC3339 (nano), HTTP timestamps, and others.
func TimeWithLocale(value string, layouts ...string) (time.Time, bool) {

	if len(layouts) == 0 {
		layouts = []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02", time.RFC1123, time.RFC1123Z, time.RubyDate, time.UnixDate, time.RFC822, time.RFC822Z}
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}
