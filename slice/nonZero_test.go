package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonZero_String(t *testing.T) {

	original := []string{"a", "b", "a", "b", "", "f", "g", "", "i", "j"}
	expected := []string{"a", "b", "a", "b", "f", "g", "i", "j"}

	require.Equal(t, expected, NonZero(original))
}

func TestNonZero_Int(t *testing.T) {

	original := []int{0, 1, 2, 3, 0, 2, 3, 4, 5, 6, 7, 0}
	expected := []int{1, 2, 3, 2, 3, 4, 5, 6, 7}

	require.Equal(t, expected, NonZero(original))
}
