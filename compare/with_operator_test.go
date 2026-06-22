package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithOperator_Numeric(t *testing.T) {

	// assertOperator checks a single operator against a known result.
	assertOperator := func(value1 any, operator string, value2 any, expected bool) {
		result, err := WithOperator(value1, operator, value2)
		require.Nil(t, err)
		require.Equal(t, expected, result, "%v %s %v", value1, operator, value2)
	}

	assertOperator(2, OperatorGreaterThan, 1, true)
	assertOperator(1, OperatorGreaterThan, 2, false)
	assertOperator(1, OperatorGreaterThan, 1, false)

	assertOperator(2, OperatorGreaterOrEqual, 1, true)
	assertOperator(1, OperatorGreaterOrEqual, 1, true)
	assertOperator(1, OperatorGreaterOrEqual, 2, false)

	assertOperator(1, OperatorEqual, 1, true)
	assertOperator(1, OperatorEqual, 2, false)

	assertOperator(1, OperatorNotEqual, 2, true)
	assertOperator(1, OperatorNotEqual, 1, false)

	assertOperator(1, OperatorLessOrEqual, 1, true)
	assertOperator(1, OperatorLessOrEqual, 2, true)
	assertOperator(2, OperatorLessOrEqual, 1, false)

	assertOperator(1, OperatorLessThan, 2, true)
	assertOperator(2, OperatorLessThan, 1, false)
	assertOperator(1, OperatorLessThan, 1, false)
}

func TestWithOperator_Strings(t *testing.T) {

	{
		result, err := WithOperator("hello world", OperatorBeginsWith, "hello")
		require.Nil(t, err)
		require.True(t, result)
	}
	{
		result, err := WithOperator("hello world", OperatorEndsWith, "world")
		require.Nil(t, err)
		require.True(t, result)
	}
	{
		result, err := WithOperator("hello world", OperatorContains, "lo wo")
		require.Nil(t, err)
		require.True(t, result)
	}
	{
		result, err := WithOperator("one", OperatorContainedBy, "one two three")
		require.Nil(t, err)
		require.True(t, result)
	}
}

func TestWithOperator_UnrecognizedOperator(t *testing.T) {
	result, err := WithOperator(1, "BOGUS", 2)
	require.False(t, result)
	require.NotNil(t, err)
}

func TestWithOperator_IncompatibleTypes(t *testing.T) {
	// A string compared against an int cannot be coerced, so Interface errors out.
	result, err := WithOperator("hello", OperatorEqual, 1)
	require.False(t, result)
	require.NotNil(t, err)
}

func TestLessThanGreaterThan(t *testing.T) {

	require.True(t, LessThan(1, 2))
	require.False(t, LessThan(2, 1))
	require.False(t, LessThan(1, 1))

	require.True(t, GreaterThan(2, 1))
	require.False(t, GreaterThan(1, 2))
	require.False(t, GreaterThan(1, 1))

	// Incompatible types fall through to FALSE
	require.False(t, LessThan("hello", 1))
	require.False(t, GreaterThan("hello", 1))
}

func TestNotContains(t *testing.T) {
	require.True(t, NotContains([]string{"a", "b"}, "c"))
	require.False(t, NotContains([]string{"a", "b"}, "a"))
}

func TestNotNil(t *testing.T) {
	require.True(t, NotNil(0))
	require.True(t, NotNil("hello"))
	require.False(t, NotNil(nil))
	require.False(t, NotNil((*int)(nil)))
}
