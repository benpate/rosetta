package funcmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDateFuncs(t *testing.T) {

	f := All()

	// A fixed reference time: March 4, 2026 1:02 PM UTC
	when := time.Date(2026, time.March, 4, 13, 2, 0, 0, time.UTC)

	require.IsType(t, time.Time{}, f["now"].(func() time.Time)())
	require.IsType(t, time.Time{}, f["today"].(func() time.Time)())
	require.IsType(t, time.Time{}, f["yesterday"].(func() time.Time)())

	require.Equal(t, "March 4, 2026 1:02 PM", f["dateTime"].(func(any) string)(when))
	require.Equal(t, "4", f["day"].(func(any) string)(when))
	require.Equal(t, when.Unix(), f["epochDate"].(func(any) int64)(when))
	require.Equal(t, when.Format(time.RFC3339), f["isoDate"].(func(any) string)(when))
	require.Equal(t, "Wednesday, March 4, 2026", f["longDate"].(func(any) string)(when))
	require.Equal(t, "March", f["longMonth"].(func(any) string)(when))
	require.Equal(t, "Mar 4, 2026", f["shortDate"].(func(any) string)(when))
	require.Equal(t, "Mar", f["shortMonth"].(func(any) string)(when))
	require.Equal(t, "1:02 PM", f["shortTime"].(func(any) string)(when))
	require.Equal(t, "2026", f["year"].(func(any) string)(when))
}

func TestDateFuncs_ZeroAndEmpty(t *testing.T) {

	f := All()

	// A zero time produces empty strings for the formatting helpers
	require.Equal(t, "", f["dateTime"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["day"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["isoDate"].(func(any) string)(""))
	require.Equal(t, "", f["longDate"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["longMonth"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["shortDate"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["shortMonth"].(func(any) string)(time.Time{}))
	require.Equal(t, "", f["shortTime"].(func(any) string)(time.Time{}))

	// year short-circuits on an empty string input
	require.Equal(t, "", f["year"].(func(any) string)(""))
	require.Equal(t, "", f["year"].(func(any) string)(time.Time{}))
}
