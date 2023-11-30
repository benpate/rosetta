package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnique_String(t *testing.T) {

	original := []string{"a", "b", "c", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	expected := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	require.Equal(t, expected, Unique(original))
}

func TestUnique_Int(t *testing.T) {

	original := []int{1, 2, 3, 1, 2, 3, 4, 5, 6, 7}
	expected := []int{1, 2, 3, 4, 5, 6, 7}

	require.Equal(t, expected, Unique(original))
}
