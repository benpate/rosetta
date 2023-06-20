package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeInt(t *testing.T) {

	zeroTime := time.Time{}
	oneTime := zeroTime.Add(time.Second)
	baseTime := time.Date(2023, 1, 2, 3, 4, 5, 0, time.Local)

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
