package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt_EmptyConstructor(t *testing.T) {

	result := NewInt()
	require.NotNil(t, result)
	require.Zero(t, result.Length())
}
