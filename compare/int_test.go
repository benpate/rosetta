package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// assertCompare exercises a comparison function with values known to be
// less-than, equal-to, and greater-than each other, expecting -1, 0, and 1.
func assertCompare[T any](t *testing.T, fn func(T, T) int, low T, high T) {
	require.Equal(t, -1, fn(low, high), "expected low < high")
	require.Equal(t, 0, fn(low, low), "expected low == low")
	require.Equal(t, 1, fn(high, low), "expected high > low")
}

func TestCompareInt(t *testing.T) {
	assertCompare(t, Int, 1, 2)
	assertCompare(t, Int8, int8(1), int8(2))
	assertCompare(t, Int16, int16(1), int16(2))
	assertCompare(t, Int32, int32(1), int32(2))
	assertCompare(t, Int64, int64(1), int64(2))
}

func TestCompareUInt(t *testing.T) {
	assertCompare(t, UInt, uint(1), uint(2))
	assertCompare(t, UInt8, uint8(1), uint8(2))
	assertCompare(t, UInt16, uint16(1), uint16(2))
	assertCompare(t, UInt32, uint32(1), uint32(2))
	assertCompare(t, UInt64, uint64(1), uint64(2))
}

func TestCompareInt_Negatives(t *testing.T) {
	require.Equal(t, -1, Int(-5, -1))
	require.Equal(t, 1, Int(-1, -5))
	require.Equal(t, 0, Int(-3, -3))
}
