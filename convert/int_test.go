package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: LOW: These tests still don't account for overflow errors.

func TestInt(t *testing.T) {
	require.Equal(t, IntDefault(10, 10), int(10))
}

func TestNilToInt(t *testing.T) {

	result, natural := IntOk(nil, int(-1))
	require.Equal(t, result, int(-1))
	require.False(t, natural)
}

func TestFloat32ToInt(t *testing.T) {

	result, natural := IntOk(float32(10), -1)

	require.True(t, natural)
	require.Equal(t, result, int(10))
}

func TestFloat64ToInt(t *testing.T) {

	result, natural := IntOk(float64(10), -1)

	require.True(t, natural)
	require.Equal(t, result, int(10))
}

func TestIntToInt(t *testing.T) {

	{
		result, natural := IntOk(int(10), -1)

		require.True(t, natural)
		require.Equal(t, result, int(10))
	}

	{
		result, natural := IntOk(int8(10), -1)

		require.True(t, natural)
		require.Equal(t, result, int(10))
	}

	{
		result, natural := IntOk(int16(10), -1)

		require.True(t, natural)
		require.Equal(t, result, int(10))
	}

	{
		result, natural := IntOk(int32(10), -1)

		require.True(t, natural)
		require.Equal(t, result, int(10))
	}

	{
		result, natural := IntOk(int64(10), -1)

		require.True(t, natural)
		require.Equal(t, result, int(10))
	}
}

func TestStringToInt(t *testing.T) {

	{
		result, ok := IntOk("0", -1)
		require.True(t, ok)
		require.Zero(t, result)
	}

	{
		result, ok := IntOk("1", -1)
		require.True(t, ok)
		require.Equal(t, 1, result)
	}

	{
		result, natural := IntOk("10", -1)
		require.True(t, natural)
		require.Equal(t, result, int(10))
	}

	{
		result, natural := IntOk("invalid", -1)
		require.False(t, natural)
		require.Equal(t, result, int(-1))
	}
}

func TestStringArrayInt(t *testing.T) {
	s := []string{"100", "200", "300"}

	result, natural := IntOk(s, -1)

	require.False(t, natural)
	require.Equal(t, result, int(100))
}

func TestStringerToInt(t *testing.T) {

	s := getTestStringer()

	result, natural := IntOk(s, -1)

	require.False(t, natural)
	require.Equal(t, result, int(-1))
}

func TestInvalidToInt(t *testing.T) {
	result, natural := IntOk(map[string]any{}, -1)

	require.False(t, natural)
	require.Equal(t, result, int(-1))
}
