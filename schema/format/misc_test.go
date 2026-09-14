package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchRegex(t *testing.T) {

	validate := MatchRegex(`^[a-z]+$`)

	{
		result, err := validate("hello")
		require.NoError(t, err)
		require.Equal(t, "hello", result)
	}
	{
		_, err := validate("Hello123")
		require.Error(t, err)
	}
}

func TestMatchRegex_InvalidPattern(t *testing.T) {

	// An un-compilable pattern produces an error when the validator runs
	validate := MatchRegex(`[`)
	_, err := validate("anything")
	require.Error(t, err)
}

func TestIn(t *testing.T) {

	validate := In("red,green,blue")

	{
		result, err := validate("green")
		require.NoError(t, err)
		require.Equal(t, "green", result)
	}
	{
		_, err := validate("purple")
		require.Error(t, err)
	}
}

func TestNotIn(t *testing.T) {

	validate := NotIn("red,green,blue")

	// RULE: An accepted value is returned as itself. NotIn used to return `arg` here, which
	// silently overwrote every field it validated with its own option list.
	{
		result, err := validate("purple")
		require.NoError(t, err)
		require.Equal(t, "purple", result)
	}
	{
		_, err := validate("green")
		require.Error(t, err)
	}
}

func TestUnsafeAny(t *testing.T) {

	result, err := UnsafeAny("")("<script>anything</script>")
	require.NoError(t, err)
	require.Equal(t, "<script>anything</script>", result)
}

func TestHasNumbers(t *testing.T) {

	// arg "2" requires at least 2 numeric characters
	validate := HasNumbers("2")

	{
		result, err := validate("ab12")
		require.NoError(t, err)
		require.Equal(t, "ab12", result)
	}
	{
		_, err := validate("ab1")
		require.Error(t, err)
	}
}
