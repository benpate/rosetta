package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: LOW: Thes tests still dont account for overflow errors.

func TestNilToFloat(t *testing.T) {

	result, natural := FloatOk(nil, float64(-1))
	assert.Equal(t, result, float64(-1))
	assert.False(t, natural)
}

func TestFloat32ToFloat(t *testing.T) {

	result, natural := FloatOk(float32(10), -1)

	assert.True(t, natural)
	assert.Equal(t, result, float64(10))
}

func TestFloatToFloat(t *testing.T) {

	assert.Equal(t, float64(10), FloatDefault(float64(10), -1))
	assert.Equal(t, float64(-1), FloatDefault("hello", -1))
}

func TestIntToFloat(t *testing.T) {

	{
		result, natural := FloatOk(int(10), -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}

	{
		result, natural := FloatOk(int8(10), -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}

	{
		result, natural := FloatOk(int16(10), -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}

	{
		result, natural := FloatOk(int32(10), -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}

	{
		result, natural := FloatOk(int64(10), -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}
}

func TestStringToFloat(t *testing.T) {

	{
		result, natural := FloatOk("10", -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(10))
	}

	{
		result, natural := FloatOk("invalid", -1)

		assert.False(t, natural)
		assert.Equal(t, result, float64(-1))
	}
}

func TestStringerToFloat(t *testing.T) {

	s := getTestStringer()

	{
		result, natural := FloatOk(s, -1)

		assert.False(t, natural)
		assert.Equal(t, result, float64(-1))
	}

	s[0] = "100"
	{
		result, natural := FloatOk(s, -1)

		assert.True(t, natural)
		assert.Equal(t, result, float64(100))
	}
}

func TestInvalidToFloat(t *testing.T) {
	result, natural := FloatOk(map[string]interface{}{}, -1)

	assert.False(t, natural)
	assert.Equal(t, result, float64(-1))
}
