package funcmap

import (
	"time"

	"github.com/benpate/rosetta/convert"
)

func addDateFuncs(target map[string]interface{}) {

	target["now"] = time.Now

	target["today"] = func() time.Time {
		return time.Now().Truncate(24 * time.Hour)
	}

	target["dateTime"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("January 2, 2006 3:04 PM")
	}

	target["day"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}

		return valueTime.Format("2")
	}

	target["epochDate"] = func(value any) int64 {
		return convert.Time(value).Unix()
	}

	target["isoDate"] = func(value any) string {

		if valueTime, ok := convert.TimeOk(value, time.Time{}); ok {
			return valueTime.Format(time.RFC3339)
		}

		return ""
	}

	target["longDate"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("Monday, January 2, 2006")
	}

	target["longMonth"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("January")
	}

	target["shortDate"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("Jan 2, 2006")
	}

	target["shortMonth"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}

		return valueTime.Format("Jan")
	}

	target["shortTime"] = func(value any) string {
		valueTime := convert.Time(value)

		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("3:04 PM")
	}

	target["year"] = func(value any) string {

		if value == "" {
			return ""
		}

		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("2006")
	}

	target["yesterday"] = func() time.Time {
		return time.Now().Truncate(24*time.Hour).AddDate(0, 0, -1)
	}

}
