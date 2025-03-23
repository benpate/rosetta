package funcmap

import (
	"time"

	"github.com/benpate/rosetta/convert"
)

func addDateFuncs(target map[string]interface{}) {

	target["now"] = time.Now

	target["isoDate"] = func(value any) string {

		if valueTime, ok := convert.TimeOk(value, time.Time{}); ok {
			return valueTime.Format(time.RFC3339)
		}

		return ""
	}

	target["epochDate"] = func(value any) int64 {
		return convert.Time(value).Unix()
	}

	target["longMonth"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("January")
	}

	target["shortMonth"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}

		return valueTime.Format("Jan")
	}

	target["day"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}

		return valueTime.Format("2")
	}

	target["shortDate"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("Jan 2, 2006")
	}

	target["longDate"] = func(value any) string {
		valueTime := convert.Time(value)
		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("Monday, January 2, 2006")
	}

	target["shortTime"] = func(value any) string {
		valueTime := convert.Time(value)

		if valueTime.IsZero() {
			return ""
		}
		return valueTime.Format("3:04 PM")
	}

}
