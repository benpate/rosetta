package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendToNil(t *testing.T) {

	// Fake nil
	{
		var value []string

		result, index, err := AnyAppend(value, "howdy")
		require.Nil(t, err)
		require.Equal(t, 1, len(result.([]string)))
		require.Equal(t, 0, index)
		require.Equal(t, []string{"howdy"}, result)
	}

	// Real nil
	{
		value := map[string]interface{}{}

		result, index, err := AnyAppend(value["empty"], "full")
		require.Nil(t, err)
		require.Equal(t, 1, len(result.([]string)))
		require.Equal(t, 0, index)
		require.Equal(t, []string{"full"}, result)
	}
}

func TestAppendToArray(t *testing.T) {

	value := [0]string{}

	result1, index, err := AnyAppend(value, "hello")
	require.Nil(t, err)
	require.Equal(t, 0, index)
	require.Equal(t, 1, len(result1.([1]string)))
	require.Equal(t, "hello", result1.([1]string)[0])

	result2, index, err := AnyAppend(result1, "there")
	require.Nil(t, err)
	require.Equal(t, 1, index)
	require.Equal(t, 2, len(result2.([2]string)))
	require.Equal(t, "there", result2.([2]string)[1])

	result3, index, err := AnyAppend(result2, "general")
	require.Nil(t, err)
	require.Equal(t, 2, index)
	require.Equal(t, 3, len(result3.([3]string)))
	require.Equal(t, "general", result3.([3]string)[2])

	result4, index, err := AnyAppend(result3, "kenobi")
	require.Nil(t, err)
	require.Equal(t, 3, index)
	require.Equal(t, 4, len(result4.([4]string)))
	require.Equal(t, "kenobi", result4.([4]string)[3])

	require.Equal(t, [4]string{"hello", "there", "general", "kenobi"}, result4)
}

func TestAppendToSlice(t *testing.T) {

	value := []string{}

	result, index, err := AnyAppend(value, "alpha")
	require.Equal(t, 1, len(result.([]string)))
	require.Equal(t, 0, index)
	require.Nil(t, err)

	result, index, err = AnyAppend(result, "bravo")
	require.Equal(t, 2, len(result.([]string)))
	require.Equal(t, 1, index)
	require.Nil(t, err)

	result, index, err = AnyAppend(result, "charlie")
	require.Equal(t, 3, len(result.([]string)))
	require.Equal(t, 2, index)
	require.Nil(t, err)

	result, index, err = AnyAppend(result, "delta")
	require.Equal(t, 4, len(result.([]string)))
	require.Equal(t, 3, index)
	require.Nil(t, err)
}
