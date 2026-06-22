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

	{
		result, err := validate("purple")
		require.NoError(t, err)
		require.Equal(t, "red,green,blue", result) // NotIn returns the arg on success
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

func TestWebFinger(t *testing.T) {

	validate := WebFinger("")

	{
		// Valid handle, with the leading @ stripped from the result
		result, err := validate("@sara@sky.net")
		require.NoError(t, err)
		require.Equal(t, "sara@sky.net", result)
	}
	{
		// Empty string is allowed
		result, err := validate("")
		require.NoError(t, err)
		require.Equal(t, "", result)
	}
	{
		// Missing leading @
		_, err := validate("sara@sky.net")
		require.Error(t, err)
	}
	{
		// Not a valid email after the @
		_, err := validate("@not an email")
		require.Error(t, err)
	}
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
