package format

import (
	"regexp"
	"time"

	"github.com/benpate/derp"
)

var iso8601 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(:\d{2}(Z[+-]\d{2}:\d{2})?)?$`)

// ISO8601 returns a StringFormat that validates a value as an ISO-8601 date string.
func ISO8601(arg string) StringFormat {
	return func(value string) (string, error) {
		if iso8601.Match([]byte(value)) {
			return value, nil
		}

		return "", derp.Validation("Value must be a valid date string", value)
	}
}

// Date returns a StringFormat that validates a value as an RFC-3339 full-date (e.g. "2026-03-04").
func Date(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		if _, err := time.Parse("2006-01-02", value); err != nil {
			return "", derp.Validation("Value is not a valid date", value)
		}

		return value, nil
	}
}

// DateTime returns a StringFormat that validates a value as an RFC-3339 date-time (e.g. "2026-03-04T13:02:00Z").
func DateTime(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// Accept both second- and fractional-second precision.
		if _, err := time.Parse(time.RFC3339, value); err == nil {
			return value, nil
		}
		if _, err := time.Parse(time.RFC3339Nano, value); err == nil {
			return value, nil
		}

		return "", derp.Validation("Value is not a valid date-time", value)
	}
}

// Time returns a StringFormat that validates a value as an RFC-3339 full-time (e.g. "13:02:00" or "13:02:00Z").
func Time(arg string) StringFormat {

	return func(value string) (string, error) {

		// Allow empty values
		if value == "" {
			return "", nil
		}

		// Accept time with seconds, with an optional zone, or to the minute.
		for _, layout := range []string{"15:04:05", "15:04:05Z07:00", "15:04"} {
			if _, err := time.Parse(layout, value); err == nil {
				return value, nil
			}
		}

		return "", derp.Validation("Value is not a valid time", value)
	}
}
