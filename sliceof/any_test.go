package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny_EmptyConstructor(t *testing.T) {

	result := NewAny()
	require.NotNil(t, result)
	require.Zero(t, result.Length())
}

func TestAny_Append(t *testing.T) {

	x := Any{1, "hello", true}

	x.Append(42.0, "world", false)

	require.Equal(t, Any{1, "hello", true, 42.0, "world", false}, x)
}

func TestAny_LengthGetter(t *testing.T) {

	x := Any{1, "hello", true}

	require.Equal(t, 3, x.Length())
	require.Equal(t, 3, (&x).Length())
}

func TestAny_NewConstructor(t *testing.T) {
	slice := NewAny("zero", "one", "two", 3)
	require.Equal(t, 4, slice.Length())
	require.Equal(t, "zero", slice[0])
	require.Equal(t, "one", slice[1])
	require.Equal(t, "two", slice[2])
	require.Equal(t, 3, slice[3])
}
