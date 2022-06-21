package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet_Slice(t *testing.T) {
	value := []string{"zero", "one", "two", "three"}

	v0 := Get(value, "")
	require.Equal(t, []string{"zero", "one", "two", "three"}, v0)

	p0, ok := GetOK(value, "0")
	require.True(t, ok)
	require.Equal(t, "zero", p0)

	p1, ok := GetOK(value, "1")
	require.True(t, ok)
	require.Equal(t, "one", p1)

	p2, ok := GetOK(value, "2")
	require.True(t, ok)
	require.Equal(t, "two", p2)

	p3, ok := GetOK(value, "3")
	require.True(t, ok)
	require.Equal(t, "three", p3)

	p4, ok := GetOK(value, "4")
	require.False(t, ok)
	require.Nil(t, p4)

	p_1, ok := GetOK(value, "-1")
	require.False(t, ok)
	require.Nil(t, p_1)
}

func TestProperties(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "name")
		require.Equal(t, "John Connor", value)
	}

	{
		value := Get(d, "email")
		require.Equal(t, "john@connor.mil", value)
	}

	{
		value := Get(d, "missing property")
		require.Nil(t, value)
	}
}

func TestSubProperties(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "relatives.mom")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "relatives.dad")
		require.Equal(t, "Kyle Reese", value)
	}

	{
		value := Get(d, "relatives.sister")
		require.Nil(t, value)
	}
}

func TestArrays(t *testing.T) {

	d := getTestData()

	{
		value := Get(d, "enemies.0")
		require.Equal(t, "T-1000", value)
	}

	{
		value := Get(d, "enemies.1")
		require.Equal(t, "T-3000", value)
	}

	{
		value := Get(d, "enemies.2")
		require.Equal(t, "T-5000", value)
	}

	{
		value := Get(d, "enemies.-1")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.3")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.100000")
		require.Nil(t, value)
	}

	{
		value := Get(d, "enemies.fred")
		require.Nil(t, value)
	}
}

func TestError(t *testing.T) {

	{
		value := Get("unsupported data", "property")
		require.Nil(t, value)
	}

	{
		value := Get("string at the end of a path", "")
		require.Equal(t, "string at the end of a path", value)
	}
}

func TestGetter(t *testing.T) {

	d := getTestStruct()

	{
		value := Get(d, "name")
		require.Equal(t, "John Connor", value)
	}

	{
		value := Get(d, "email")
		require.Equal(t, "john@connor.mil", value)
	}

	{
		value := Get(d, "relatives.0.name")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "relatives.1.relatives.1.name")
		require.Equal(t, "Sarah Connor", value)
	}

	{
		value := Get(d, "missing-property")
		require.Nil(t, value)
	}
}
