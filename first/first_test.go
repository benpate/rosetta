package first

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	require.Equal(t, "one", String("one", "two", "three"))
	require.Equal(t, "two", String("", "two", "three"))
	require.Equal(t, "three", String("", "", "three"))
	require.Equal(t, "", String("", "", ""))
}

func TestInt(t *testing.T) {
	require.Equal(t, 1, Int(1, 2, 3))
	require.Equal(t, 2, Int(0, 2, 3))
	require.Equal(t, 3, Int(0, 0, 3))
	require.Equal(t, 0, Int(0, 0, 0))
}

func TestInt64(t *testing.T) {
	require.Equal(t, int64(1), Int64(1, 2, 3))
	require.Equal(t, int64(2), Int64(0, 2, 3))
	require.Equal(t, int64(3), Int64(0, 0, 3))
	require.Equal(t, int64(0), Int64(0, 0, 0))
}
