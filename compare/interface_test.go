package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Numeric type comparisons (the full signed/unsigned/float matrix), signed-vs-unsigned
// edge cases, boundary values, and numeric-vs-non-numeric pairings are exercised in
// interface_matrix_test.go. The tests below cover the remaining non-numeric paths.

func TestInterface_String(t *testing.T) {

	result, err := Interface("apple", "banana")
	require.Nil(t, err)
	require.Equal(t, -1, result)

	result, err = Interface("apple", "apple")
	require.Nil(t, err)
	require.Equal(t, 0, result)

	result, err = Interface("banana", "apple")
	require.Nil(t, err)
	require.Equal(t, 1, result)
}

func TestInterface_Stringer(t *testing.T) {

	result, err := Interface(zeroStringer("apple"), zeroStringer("banana"))
	require.Nil(t, err)
	require.Equal(t, -1, result)
}

func TestInterface_Bool(t *testing.T) {

	result, err := Interface(true, true)
	require.Nil(t, err)
	require.Equal(t, 0, result)
}

func TestInterface_IncompatibleTypes(t *testing.T) {

	// completely unknown type
	_, err := Interface(struct{}{}, struct{}{})
	require.NotNil(t, err)
}
