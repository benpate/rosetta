package convert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilToFloat(t *testing.T) {

	result, lossless := FloatOk(nil, float64(-1))
	assert.Equal(t, result, float64(-1))
	assert.False(t, lossless)
}

func TestNaNToFloat(t *testing.T) {

	result, lossless := FloatOk(math.NaN(), float64(-1))
	assert.True(t, math.IsNaN(result))
	assert.True(t, lossless)
}

func TestFloat32ToFloat(t *testing.T) {

	result, lossless := FloatOk(float32(10), -1)

	assert.True(t, lossless)
	assert.Equal(t, result, float64(10))
}

func TestFloatToFloat(t *testing.T) {

	assert.Equal(t, float64(10), FloatDefault(float64(10), -1))
	assert.Equal(t, float64(-1), FloatDefault("hello", -1))
}

func TestIntToFloat(t *testing.T) {

	{
		result, lossless := FloatOk(int(10), -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}

	{
		result, lossless := FloatOk(int8(10), -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}

	{
		result, lossless := FloatOk(int16(10), -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}

	{
		result, lossless := FloatOk(int32(10), -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}

	{
		result, lossless := FloatOk(int64(10), -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}
}

func TestStringToFloat(t *testing.T) {

	{
		result, lossless := FloatOk("10", -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(10))
	}

	{
		result, lossless := FloatOk("invalid", -1)

		assert.False(t, lossless)
		assert.Equal(t, result, float64(-1))
	}
}

func TestStringerToFloat(t *testing.T) {

	s := getTestStringer()

	{
		result, lossless := FloatOk(s, -1)

		assert.False(t, lossless)
		assert.Equal(t, result, float64(-1))
	}

	s[0] = "100"
	{
		result, lossless := FloatOk(s, -1)

		assert.True(t, lossless)
		assert.Equal(t, result, float64(100))
	}
}

func TestInvalidToFloat(t *testing.T) {
	result, lossless := FloatOk(map[string]any{}, -1)

	assert.False(t, lossless)
	assert.Equal(t, result, float64(-1))
}
