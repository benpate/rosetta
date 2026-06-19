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
	{
		result, err := validate("not a date")
		require.NotNil(t, err)
		require.Equal(t, "", result)
	}
}

// Date, DateTime, and Time are currently pass-through no-ops.
func TestCalendar_Passthrough(t *testing.T) {

	require.NoError(t, callFormat(t, Date(""), "anything"))
	require.NoError(t, callFormat(t, DateTime(""), "anything"))
	require.NoError(t, callFormat(t, Time(""), "anything"))
}

// callFormat runs a StringFormat and confirms the value passes through unchanged,
// returning any validation error.
func callFormat(t *testing.T, format StringFormat, value string) error {
	t.Helper()
	result, err := format(value)
	if err == nil {
		require.Equal(t, value, result)
	}
	return err
}
