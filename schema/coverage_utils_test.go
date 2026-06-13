package schema

import (
	"testing"

	"github.com/benpate/rosetta/schema/format"
	"github.com/stretchr/testify/require"
)

func TestIndirect_Pointer(t *testing.T) {
	value := "hello"
	require.Equal(t, "hello", indirect(&value))
}

func TestIndirect_NotPointer(t *testing.T) {
	require.Equal(t, "hello", indirect("hello"))
}

func TestGetLength(t *testing.T) {
	length, ok := getLength(testArrayA{"one", "two"})
	require.True(t, ok)
	require.Equal(t, 2, length)
}

func TestGetLength_NotLengthGetter(t *testing.T) {
	_, ok := getLength("not-an-array")
	require.False(t, ok)
}

func TestGetIndex(t *testing.T) {
	item, ok := getIndex(testArrayA{"one", "two"}, 0)
	require.True(t, ok)
	require.Equal(t, "one", item)
}

func TestGetIndex_NotArrayGetter(t *testing.T) {
	_, ok := getIndex("not-an-array", 0)
	require.False(t, ok)
}

func TestIsMultipleOf(t *testing.T) {
	require.True(t, isMultipleOf(10, 5))
	require.False(t, isMultipleOf(10, 3))
}

func TestNotMultipleOf(t *testing.T) {
	require.False(t, notMultipleOf(10, 5))
	require.True(t, notMultipleOf(10, 3))
}

func TestType_String(t *testing.T) {
	require.Equal(t, "string", TypeString.String())
}

func TestUseFormat_NilIsIgnored(t *testing.T) {
	// Registering a nil generator is a no-op and must not panic
	require.NotPanics(t, func() { UseFormat("coverage-nil-format", nil) })
	require.NotContains(t, formats, "coverage-nil-format")
}

func TestUseFormat_RegistersGenerator(t *testing.T) {
	UseFormat("coverage-test-format", format.NoHTML)
	require.Contains(t, formats, "coverage-test-format")
	delete(formats, "coverage-test-format")
}
