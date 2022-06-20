package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Thes tests still dont account for overflow errors.

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

	{
		result, natural := StringOk(int(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}

	{
		result, natural := StringOk(int8(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}

	{
		result, natural := StringOk(int16(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}

	{
		result, natural := StringOk(int32(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}

	{
		result, natural := StringOk(int64(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}
}

func TestFloatToString(t *testing.T) {

	{
		result, natural := StringOk(float32(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}

	{
		result, natural := StringOk(float64(10), "default")

		assert.False(t, natural)
		assert.Equal(t, result, "10")
	}
}

func TestStringerToString(t *testing.T) {

	s := getTestStringer()

	{
		result, natural := StringOk(s, "")

		assert.True(t, natural)
		assert.Equal(t, result, "hello")
	}
}

func TestInvalidToString(t *testing.T) {
	result, natural := StringOk(map[string]interface{}{}, "default")

	assert.False(t, natural)
	assert.Equal(t, result, "default")
}
