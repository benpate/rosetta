package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilToBool(t *testing.T) {

	{
		result, natural := BoolOk(nil, true)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(nil, false)
		assert.False(t, result)
		assert.False(t, natural)
	}
}

func TestBoolToBool(t *testing.T) {

	assert.True(t, BoolDefault(true, true))
	assert.True(t, BoolDefault(true, false))
	assert.False(t, BoolDefault(false, true))
	assert.False(t, BoolDefault(false, false))
}

func TestIntToBool(t *testing.T) {

	{
		result, natural := BoolOk(0, false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(1, false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(0, true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(1, true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestInt8ToBool(t *testing.T) {

	{
		result, natural := BoolOk(int8(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int8(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int8(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int8(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestInt16ToBool(t *testing.T) {

	{
		result, natural := BoolOk(int16(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int16(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int16(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int16(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestInt32ToBool(t *testing.T) {

	{
		result, natural := BoolOk(int32(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int32(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int32(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int32(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestInt64ToBool(t *testing.T) {

	{
		result, natural := BoolOk(int64(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int64(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int64(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(int64(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestFloat32ToBool(t *testing.T) {

	{
		result, natural := BoolOk(float32(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float32(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float32(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float32(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestFloat64ToBool(t *testing.T) {

	{
		result, natural := BoolOk(float64(0), false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float64(1), false)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float64(0), true)
		assert.False(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(float64(1), true)
		assert.True(t, result)
		assert.False(t, natural)
	}
}

func TestStringToBool(t *testing.T) {

	{
		result, natural := BoolOk("true", true)
		assert.True(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk("false", true)
		assert.False(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk("true", false)
		assert.True(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk("false", false)
		assert.False(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk("Somethig Else", true)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk("Somethig Else", false)
		assert.False(t, result)
		assert.False(t, natural)
	}
}

func TestStringerToBool(t *testing.T) {

	s := getTestStringer()

	{
		result, natural := BoolOk(s, true)
		assert.True(t, result)
		assert.False(t, natural)
	}

	{
		result, natural := BoolOk(s, false)
		assert.False(t, result)
		assert.False(t, natural)
	}

	s[0] = "true"

	{
		result, natural := BoolOk(s, true)
		assert.True(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk(s, false)
		assert.True(t, result)
		assert.True(t, natural)
	}

	s[0] = "false"

	{
		result, natural := BoolOk(s, true)
		assert.False(t, result)
		assert.True(t, natural)
	}

	{
		result, natural := BoolOk(s, false)
		assert.False(t, result)
		assert.True(t, natural)
	}
}

func TestInvalidToBool(t *testing.T) {

	result, natural := BoolOk(map[string]interface{}{}, true)

	assert.False(t, natural)
	assert.True(t, result)
}
