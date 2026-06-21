package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests verify that narrowing integer conversions never silently truncate:
// a value outside the target range must clamp to the boundary AND return Ok=false.

func TestInt32Ok_Overflow(t *testing.T) {

	// int64 above/below int32 range
	{
		result, ok := Int32Ok(int64(math.MaxInt32)+1, -1)
		require.False(t, ok)
		require.Equal(t, int32(math.MaxInt32), result)
	}
	{
		result, ok := Int32Ok(int64(math.MinInt32)-1, -1)
		require.False(t, ok)
		require.Equal(t, int32(math.MinInt32), result)
	}

	// int above int32 range (64-bit platforms) — the previously-unguarded path
	{
		result, ok := Int32Ok(int(math.MaxInt32)+1, -1)
		require.False(t, ok)
		require.Equal(t, int32(math.MaxInt32), result)
	}

	// In-range int is still natural
	{
		result, ok := Int32Ok(int(123), -1)
		require.True(t, ok)
		require.Equal(t, int32(123), result)
	}

	// float64 at exactly 2^31 (the boundary float32(MaxInt32) rounds up to) must NOT slip through
	{
		result, ok := Int32Ok(float64(1<<31), -1)
		require.False(t, ok)
		require.Equal(t, int32(math.MaxInt32), result)
	}
	// float32 at exactly 2^31 likewise
	{
		result, ok := Int32Ok(float32(1<<31), -1)
		require.False(t, ok)
		require.Equal(t, int32(math.MaxInt32), result)
	}
}

func TestInt64Ok_Overflow(t *testing.T) {

	// float64 at exactly 2^63 (which math.MaxInt64 rounds up to) must NOT overflow through the guard
	{
		result, ok := Int64Ok(float64(1<<63), -1)
		require.False(t, ok)
		require.Equal(t, int64(math.MaxInt64), result)
	}
	// float64 below int64 range
	{
		result, ok := Int64Ok(-float64(1<<63)*2, -1)
		require.False(t, ok)
		require.Equal(t, int64(math.MinInt64), result)
	}
	// float32 at 2^63 likewise
	{
		result, ok := Int64Ok(float32(1<<63), -1)
		require.False(t, ok)
		require.Equal(t, int64(math.MaxInt64), result)
	}

	// In-range float is still natural
	{
		result, ok := Int64Ok(float64(1000), -1)
		require.True(t, ok)
		require.Equal(t, int64(1000), result)
	}
}

func TestIntOk_Overflow(t *testing.T) {

	// On 64-bit platforms, float64 at 2^63 must not overflow through the guard.
	// (On 32-bit, math.MaxInt is smaller and this value is still out of range — either way Ok=false.)
	{
		result, ok := IntOk(float64(1<<63), -1)
		require.False(t, ok)
		require.Equal(t, math.MaxInt, result)
	}

	// In-range float is still natural
	{
		result, ok := IntOk(float64(1000), -1)
		require.True(t, ok)
		require.Equal(t, 1000, result)
	}
}
