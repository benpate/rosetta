package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat_EmptyConstructor(t *testing.T) {

	result := NewFloat()
	require.NotNil(t, result)
	require.Zero(t, result.Length())
}
