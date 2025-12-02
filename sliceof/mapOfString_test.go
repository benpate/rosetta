package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfString_EmptyConstructor(t *testing.T) {

	result := NewMapOfString()
	require.NotNil(t, result)
	require.Zero(t, result.Length())
}
