package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapOfAny_EmptyConstructor(t *testing.T) {

	result := NewMapOfAny()

	require.NotNil(t, result)
	require.Zero(t, result.Length())
}
