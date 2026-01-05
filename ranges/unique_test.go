package ranges

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnique(t *testing.T) {
	inputs := []string{"apple", "banana", "apple", "orange", "banana", "kiwi", "lemon", "kiwi"}
	expected := []string{"apple", "banana", "orange", "kiwi", "lemon"}

	values := Unique(Values(inputs...))

	for actual := range values {
		require.Equal(t, expected[0], actual)
		expected = expected[1:]
	}

	require.Zero(t, len(expected))
}
