package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// assertInterface confirms that Interface reports low < high, low == low,
// and high > low for a pair of (possibly mixed-type) values.
func assertInterface(t *testing.T, low any, high any) {

	t.Helper()

	{
		result, err := Interface(low, high)
		require.Nil(t, err)
		require.Equal(t, -1, result, "expected %v < %v", low, high)
	}
	{
		result, err := Interface(low, low)
		require.Nil(t, err)
		require.Equal(t, 0, result, "expected %v == %v", low, low)
	}
	{
		result, err := Interface(high, low)
		require.Nil(t, err)
		require.Equal(t, 1, result, "expected %v > %v", high, low)
	}
}

func TestInterface_IntCombinations(t *testing.T) {
	assertInterface(t, 1, 2)
	assertInterface(t, 1, int8(2))
	assertInterface(t, 1, int16(2))
	assertInterface(t, 1, int32(2))
	assertInterface(t, 1, int64(2))
	assertInterface(t, 1, float32(2))
	assertInterface(t, 1, float64(2))
}

func TestInterface_Int8Combinations(t *testing.T) {
	assertInterface(t, int8(1), 2)
	assertInterface(t, int8(1), int8(2))
	assertInterface(t, int8(1), int16(2))
	assertInterface(t, int8(1), int32(2))
	assertInterface(t, int8(1), int64(2))
	assertInterface(t, int8(1), float32(2))
	assertInterface(t, int8(1), float64(2))
}

func TestInterface_Int16Combinations(t *testing.T) {
	assertInterface(t, int16(1), 2)
	assertInterface(t, int16(1), int8(2))
	assertInterface(t, int16(1), int16(2))
	assertInterface(t, int16(1), int32(2))
	assertInterface(t, int16(1), int64(2))
	assertInterface(t, int16(1), float32(2))
	assertInterface(t, int16(1), float64(2))
}

// TestInterface_Int16VsUint covers the one-directional int16-vs-uint path.
// (The reverse, uint-vs-int16, is not supported by Interface, so this case
// is intentionally tested in a single direction only.)
func TestInterface_Int16VsUint(t *testing.T) {
	result, err := Interface(int16(1), uint(2))
	require.Nil(t, err)
	require.Equal(t, -1, result)
}

func TestInterface_Int32Combinations(t *testing.T) {
	assertInterface(t, int32(1), 2)
	assertInterface(t, int32(1), int8(2))
	assertInterface(t, int32(1), int16(2))
	assertInterface(t, int32(1), int32(2))
	assertInterface(t, int32(1), int64(2))
	assertInterface(t, int32(1), float32(2))
	assertInterface(t, int32(1), float64(2))
}

func TestInterface_Int64Combinations(t *testing.T) {
	assertInterface(t, int64(1), 2)
	assertInterface(t, int64(1), int8(2))
	assertInterface(t, int64(1), int16(2))
	assertInterface(t, int64(1), int32(2))
	assertInterface(t, int64(1), int64(2))
	assertInterface(t, int64(1), float32(2))
	assertInterface(t, int64(1), float64(2))
}

func TestInterface_UIntCombinations(t *testing.T) {
	assertInterface(t, uint(1), uint(2))
	assertInterface(t, uint(1), uint8(2))
	assertInterface(t, uint(1), uint16(2))
	assertInterface(t, uint(1), uint32(2))
	assertInterface(t, uint(1), uint64(2))
	assertInterface(t, uint(1), float32(2))
	assertInterface(t, uint(1), float64(2))
}

func TestInterface_UInt8Combinations(t *testing.T) {
	assertInterface(t, uint8(1), uint(2))
	assertInterface(t, uint8(1), uint8(2))
	assertInterface(t, uint8(1), uint16(2))
	assertInterface(t, uint8(1), uint32(2))
	assertInterface(t, uint8(1), uint64(2))
	assertInterface(t, uint8(1), float32(2))
	assertInterface(t, uint8(1), float64(2))
}

func TestInterface_UInt16Combinations(t *testing.T) {
	assertInterface(t, uint16(1), uint(2))
	assertInterface(t, uint16(1), uint8(2))
	assertInterface(t, uint16(1), uint16(2))
	assertInterface(t, uint16(1), uint32(2))
	assertInterface(t, uint16(1), uint64(2))
	assertInterface(t, uint16(1), float32(2))
	assertInterface(t, uint16(1), float64(2))
}

func TestInterface_UInt32Combinations(t *testing.T) {
	assertInterface(t, uint32(1), uint(2))
	assertInterface(t, uint32(1), uint8(2))
	assertInterface(t, uint32(1), uint16(2))
	assertInterface(t, uint32(1), uint32(2))
	assertInterface(t, uint32(1), uint64(2))
	assertInterface(t, uint32(1), float32(2))
	assertInterface(t, uint32(1), float64(2))
}

func TestInterface_UInt64Combinations(t *testing.T) {
	assertInterface(t, uint64(1), uint(2))
	assertInterface(t, uint64(1), uint8(2))
	assertInterface(t, uint64(1), uint16(2))
	assertInterface(t, uint64(1), uint32(2))
	assertInterface(t, uint64(1), uint64(2))
	assertInterface(t, uint64(1), float32(2))
	assertInterface(t, uint64(1), float64(2))
}

func TestInterface_Float32Combinations(t *testing.T) {
	assertInterface(t, float32(1), 2)
	assertInterface(t, float32(1), int8(2))
	assertInterface(t, float32(1), int16(2))
	assertInterface(t, float32(1), int32(2))
	assertInterface(t, float32(1), int64(2))
	assertInterface(t, float32(1), uint(2))
	assertInterface(t, float32(1), uint8(2))
	assertInterface(t, float32(1), uint16(2))
	assertInterface(t, float32(1), uint32(2))
	assertInterface(t, float32(1), uint64(2))
	assertInterface(t, float32(1), float32(2))
	assertInterface(t, float32(1), float64(2))
}

func TestInterface_Float64Combinations(t *testing.T) {
	assertInterface(t, float64(1), 2)
	assertInterface(t, float64(1), int8(2))
	assertInterface(t, float64(1), int16(2))
	assertInterface(t, float64(1), int32(2))
	assertInterface(t, float64(1), int64(2))
	assertInterface(t, float64(1), uint(2))
	assertInterface(t, float64(1), uint8(2))
	assertInterface(t, float64(1), uint16(2))
	assertInterface(t, float64(1), uint32(2))
	assertInterface(t, float64(1), uint64(2))
	assertInterface(t, float64(1), float32(2))
	assertInterface(t, float64(1), float64(2))
}

func TestInterface_String(t *testing.T) {
	assertInterface(t, "apple", "banana")
}

func TestInterface_Stringer(t *testing.T) {
	assertInterface(t, zeroStringer("apple"), zeroStringer("banana"))
}

func TestInterface_Bool(t *testing.T) {
	result, err := Interface(true, true)
	require.Nil(t, err)
	require.Equal(t, 0, result)
}

func TestInterface_IncompatibleTypes(t *testing.T) {

	// int vs. string is not coercible
	{
		_, err := Interface(1, "hello")
		require.NotNil(t, err)
	}

	// string vs. int is not coercible
	{
		_, err := Interface("hello", 1)
		require.NotNil(t, err)
	}

	// completely unknown type
	{
		_, err := Interface(struct{}{}, struct{}{})
		require.NotNil(t, err)
	}
}
