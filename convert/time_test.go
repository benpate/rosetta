package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeInt(t *testing.T) {

	zeroTime := time.Time{}
	oneTime := zeroTime.Add(time.Second)
	baseTime := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)

	// Test time (echo)
	{
		result, ok := TimeOk(baseTime, oneTime)
		require.True(t, ok)
		require.Equal(t, baseTime, result)
	}

	// Test seconds
	{
		result, ok := TimeOk(baseTime.Unix(), oneTime)
		require.True(t, ok)
		require.Equal(t, baseTime, result)
	}

	// Test seconds (int)
	{
		result, ok := TimeOk(int(baseTime.Unix()), oneTime)
		require.True(t, ok)
		require.Equal(t, baseTime, result)
	}

	// Test Miliseconds
	{
		result, ok := TimeOk(baseTime.UnixMilli(), oneTime)
		require.True(t, ok)
		require.Equal(t, baseTime, result)
	}

	// Test String
	{
		result, ok := TimeOk(baseTime.Format(time.RFC3339), oneTime)
		require.True(t, ok)
		require.Equal(t, baseTime, result)
	}

	// Test String (bad)
	{
		result, ok := TimeOk("bad", oneTime)
		require.False(t, ok)
		require.Equal(t, oneTime, result)
	}
}

func TestTimeString(t *testing.T) {
	result := Time("2022-09-25T14:50:32.000Z")
	require.Equal(t, time.Date(2022, 9, 25, 14, 50, 32, 0, time.UTC), result)
}

func TestTimeString2(t *testing.T) {
	result := Time("2022-09-25")
	require.Equal(t, time.Date(2022, 9, 25, 0, 0, 0, 0, time.UTC), result)
}

func TestTimeString3(t *testing.T) {
	result := Time("2023-06-20T00:00:00")
	require.Equal(t, time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC), result)
}
