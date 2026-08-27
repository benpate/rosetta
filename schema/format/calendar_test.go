package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestISO8601(t *testing.T) {

	validate := ISO8601("")

	{
		result, err := validate("2026-03-04T13:02")
		require.Nil(t, err)
		require.Equal(t, "2026-03-04T13:02", result)
	}
	{
		result, err := validate("2026-03-04T13:02:00")
		require.Nil(t, err)
		require.Equal(t, "2026-03-04T13:02:00", result)
	}
	// RULE: A canonical zone is EITHER "Z" OR a numeric offset. The pattern used to require
	// "Z" followed by an offset, which rejected every real timestamp and accepted only the
	// impossible combination of both.
	for _, valid := range []string{"2026-03-04T13:02:00Z", "2026-03-04T13:02:00+01:00", "2026-03-04T13:02:00-05:00"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	for _, invalid := range []string{"not a date", "2026-03-04T13:02:00Z+01:00", "2026-03-04T13:02Z"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestDate(t *testing.T) {

	validate := Date("")

	// Empty and valid full-date values pass through unchanged
	for _, valid := range []string{"", "2026-03-04", "1999-12-31"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	// Times, date-times, and garbage are rejected
	for _, invalid := range []string{"2026-3-4", "2026-03-04T13:02:00Z", "13:02:00", "not a date"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestDateTime(t *testing.T) {

	validate := DateTime("")

	for _, valid := range []string{"", "2026-03-04T13:02:00Z", "2026-03-04T13:02:00-05:00", "2026-03-04T13:02:00.5Z"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	for _, invalid := range []string{"2026-03-04", "13:02:00", "2026-03-04 13:02:00", "not a date-time"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}

func TestTime(t *testing.T) {

	validate := Time("")

	for _, valid := range []string{"", "13:02:00", "13:02:00Z", "13:02:00-05:00", "13:02"} {
		result, err := validate(valid)
		require.Nil(t, err, valid)
		require.Equal(t, valid, result, valid)
	}

	for _, invalid := range []string{"2026-03-04", "25:00:00", "not a time"} {
		result, err := validate(invalid)
		require.NotNil(t, err, invalid)
		require.Equal(t, "", result, invalid)
	}
}
