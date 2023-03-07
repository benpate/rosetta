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
