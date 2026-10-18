package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {

	assert.Equal(t, "hello there", StringDefault("hello there", ""))
	assert.Equal(t, "10", StringDefault(10, ""))
}

func TestNilToString(t *testing.T) {
	assert.Equal(t, "default", StringDefault(nil, "default"))
}

func TestStringToString(t *testing.T) {

	result, natural := StringOk("Hello there", "default")

	assert.True(t, natural)
	assert.Equal(t, result, "Hello there")
}

func TestByteArrayToString(t *testing.T) {
	result, natural := StringOk([]byte("...from a certain point of view."), "default")

	assert.True(t, natural)
	assert.Equal(t, result, "...from a certain point of view.")
}

func TestIntToString(t *testing.T) {

	// Integers format to decimal strings losslessly (ok=true).
	{
		result, lossless := StringOk(int(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	{
		result, lossless := StringOk(int8(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	{
		result, lossless := StringOk(int16(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	{
		result, lossless := StringOk(int32(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	{
		result, lossless := StringOk(int64(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}
}

func TestFloatToString(t *testing.T) {

	// Whole floats format as integers, and others to two decimal places.
	{
		result, lossless := StringOk(float32(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	{
		result, lossless := StringOk(float64(10), "default")

		assert.True(t, lossless)
		assert.Equal(t, result, "10")
	}

	// A value needing more than two decimals rounds and is reported lossy.
	{
		result, lossless := StringOk(float64(10.129), "default")

		assert.False(t, lossless)
		assert.Equal(t, result, "10.13")
	}

	// A fraction that fits in two decimal places keeps both.
	{
		result, lossless := StringOk(float64(2.5), "default")

		assert.True(t, lossless)
		assert.Equal(t, "2.50", result)
	}

	// A whole float, as hjson decodes `11`, converts back to an integer losslessly.
	{
		value, lossless := IntOk(String(float64(11)), 0)

		assert.True(t, lossless)
		assert.Equal(t, 11, value)
	}
}

func TestStringerToString(t *testing.T) {

	s := getTestStringer()

	result, natural := StringOk(s, "")

	assert.True(t, natural)
	assert.Equal(t, result, "hello")
}

func TestInvalidToString(t *testing.T) {
	result, natural := StringOk(map[string]any{}, "default")

	assert.False(t, natural)
	assert.Equal(t, result, "default")
}
