package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
