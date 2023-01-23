package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirst(t *testing.T) {

	{
		value := []string{"one", "two", "three"}
		keep := func(value string) bool {
			return value == "two"
		}

		require.Equal(t, "two", First(value, keep))
	}

	{
		value := []string{"one", "two", "three"}
		keep := func(value string) bool {
			return value == "four"
		}

		require.Equal(t, "", First(value, keep))
	}
}
