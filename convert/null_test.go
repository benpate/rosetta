package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBool(t *testing.T) {

	{
		v := NullBool(true)
		assert.True(t, v.Bool())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullBool(false)
		assert.False(t, v.Bool())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullBool("true")
		assert.True(t, v.Bool())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullBool("")
		assert.Equal(t, false, v.Bool())
		assert.False(t, v.IsPresent())
		assert.True(t, v.IsNull())
	}
}

func TestNullInt(t *testing.T) {

	{
		v := NullInt(1)
		assert.Equal(t, 1, v.Int())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullInt("1234567890")
		assert.Equal(t, 1234567890, v.Int())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullInt("")
		assert.Equal(t, 0, v.Int())
		assert.False(t, v.IsPresent())
		assert.True(t, v.IsNull())
	}
}

func TestNullInt64(t *testing.T) {

	{
		v := NullInt64(1)
		assert.Equal(t, int64(1), v.Int64())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullInt64("9223372036854775800")
		assert.Equal(t, int64(9223372036854775800), v.Int64())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullInt64("")
		assert.Equal(t, int64(0), v.Int64())
		assert.False(t, v.IsPresent())
		assert.True(t, v.IsNull())
	}
}

func TestNullFloat(t *testing.T) {

	{
		v := NullFloat(1.1)
		assert.Equal(t, 1.1, v.Float())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullFloat("1.1")
		assert.Equal(t, 1.1, v.Float())
		assert.True(t, v.IsPresent())
		assert.False(t, v.IsNull())
	}

	{
		v := NullFloat("")
		assert.Equal(t, 0.0, v.Float())
		assert.False(t, v.IsPresent())
		assert.True(t, v.IsNull())
	}
}
