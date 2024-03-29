package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceOfString(t *testing.T) {

	{
		value := []any{"first", "second", "third"}
		require.Equal(t, []string{"first", "second", "third"}, SliceOfString(value))
	}

	{
		value := []any{1, 2, 3}
		require.Equal(t, []string{"1", "2", "3"}, SliceOfString(value))
	}

	{
		value := []string{"1", "2", "3"}
		require.Equal(t, []string{"1", "2", "3"}, SliceOfString(value))
	}

	{
		value := "hello"
		require.Equal(t, []string{"hello"}, SliceOfString(value))
	}
}
