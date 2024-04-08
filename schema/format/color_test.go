package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColor(t *testing.T) {
	{
		result, err := Color("")("#123456")
		require.Nil(t, err)
		require.Equal(t, "#123456", result)
	}

	{
		result, err := Color("")("#abcdef")
		require.Nil(t, err)
		require.Equal(t, "#abcdef", result)
	}

	{
		result, err := Color("")("#ABCDEF")
		require.Nil(t, err)
		require.Equal(t, "#ABCDEF", result)
	}

	{
		result, err := Color("")("#123456A")
		require.Error(t, err)
		require.Equal(t, "", result)
	}

	{
		result, err := Color("")("not hex")
		require.Error(t, err)
		require.Equal(t, "", result)
	}

	{
		result, err := Color("")("#nothex")
		require.Error(t, err)
		require.Equal(t, "", result)
	}
}
