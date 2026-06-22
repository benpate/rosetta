package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareString(t *testing.T) {
	assertCompare(t, String, "apple", "banana")
	require.Equal(t, 0, String("", ""))
}

func TestBeginsWith_String(t *testing.T) {
	require.True(t, BeginsWith("hello world", "hello"))
	require.True(t, BeginsWith("hello world", ""))
	require.False(t, BeginsWith("hello world", "world"))

	// value2 is not a string
	require.False(t, BeginsWith("hello world", 42))
}

func TestBeginsWith_StringSlice(t *testing.T) {
	require.True(t, BeginsWith([]string{"first", "second"}, "first"))
	require.False(t, BeginsWith([]string{"first", "second"}, "second"))

	// Empty slice has no first element
	require.False(t, BeginsWith([]string{}, "first"))

	// value2 is not a string
	require.False(t, BeginsWith([]string{"first"}, 1))
}

func TestBeginsWith_OtherType(t *testing.T) {
	require.False(t, BeginsWith(42, "4"))
}

func TestEndsWith_String(t *testing.T) {
	require.True(t, EndsWith("hello world", "world"))
	require.True(t, EndsWith("hello world", ""))
	require.False(t, EndsWith("hello world", "hello"))

	// value2 is not a string
	require.False(t, EndsWith("hello world", 42))
}

func TestEndsWith_StringSlice(t *testing.T) {
	require.True(t, EndsWith([]string{"first", "second"}, "second"))
	require.False(t, EndsWith([]string{"first", "second"}, "first"))

	// Empty slice has no last element
	require.False(t, EndsWith([]string{}, "second"))

	// value2 is not a string
	require.False(t, EndsWith([]string{"first"}, 1))
}

func TestEndsWith_OtherType(t *testing.T) {
	require.False(t, EndsWith(42, "2"))
}
