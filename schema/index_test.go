package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {

	{
		result, ok := Index("0", 10)
		require.Equal(t, 0, result)
		require.True(t, ok)
	}

	{
		result, ok := Index("1", 10)
		require.Equal(t, 1, result)
		require.True(t, ok)
	}

	{
		result, ok := Index("9", 10)
		require.Equal(t, 9, result)
		require.True(t, ok)
	}

	{
		result, ok := Index("10", 10)
		require.Equal(t, 9, result)
		require.False(t, ok)
	}

	{
		result, ok := Index("11", 10)
		require.Equal(t, 9, result)
		require.False(t, ok)
	}

	{
		result, ok := Index("-1", 10)
		require.Equal(t, 0, result)
		require.False(t, ok)
	}

	{
		result, ok := Index("NOT AN INTEGER", 10)
		require.Equal(t, 0, result)
		require.False(t, ok)
	}
}
