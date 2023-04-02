package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterfaceBool(t *testing.T) {

	require.True(t, Equal(true, true))
	require.False(t, Equal(true, false))
	require.False(t, Equal(false, true))
	require.True(t, Equal(false, false))
}

func TestInterfaceInt(t *testing.T) {

	{
		result, err := Interface(1, 1)

		require.Equal(t, 0, result)
		require.Nil(t, err)
	}

	{
		result, err := Interface(int8(1), int8(1))

		require.Equal(t, 0, result)
		require.Nil(t, err)
	}

	{
		result, err := Interface(int16(1), int16(1))

		require.Equal(t, 0, result)
		require.Nil(t, err)
	}

}

func TestWithOperatorStringSlice(t *testing.T) {

	{
		result, err := WithOperator([]string{"one", "two", "three"}, OperatorContains, "one")
		require.True(t, result)
		require.Nil(t, err)
	}
	{
		result, err := WithOperator([]string{"one", "two", "three"}, OperatorContains, "two")
		require.True(t, result)
		require.Nil(t, err)
	}
	{
		result, err := WithOperator([]string{"one", "two", "three"}, OperatorContains, "three")
		require.True(t, result)
		require.Nil(t, err)
	}
	{
		result, err := WithOperator([]string{"one", "two", "three"}, OperatorContains, "four")
		require.False(t, result)
		require.Nil(t, err)
	}
}
