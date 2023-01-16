package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringType(t *testing.T) {

	s := String{}

	assert.Nil(t, s.Validate("I'm a string"))
	assert.NotNil(t, s.Validate(0))
	assert.NotNil(t, s.Validate([]string{}))
	assert.NotNil(t, s.Validate(map[string]any{}))
}

func TestStringUnmarshal(t *testing.T) {

	j := []byte(`{"type":"string", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(String).Required)

	require.Nil(t, s.Validate("this should work"))
	require.NotNil(t, s.Validate(""))
}

func TestStringRequired(t *testing.T) {

	s := String{Required: true}

	require.NotNil(t, s.Validate(""))
	require.Nil(t, s.Validate("present"))
}

func TestStringLength(t *testing.T) {

	// No Min/Max Length Defined
	{
		s := String{}
		assert.Nil(t, s.Validate(""))
		assert.Nil(t, s.Validate("ok."))
		assert.Nil(t, s.Validate("this is a really long string but it should be ok."))
	}

	// Mininum Defined
	{
		s := String{MinLength: 10}
		assert.Nil(t, s.Validate("this is ok, becuase it's more than the minimum."))
		assert.NotNil(t, s.Validate("error"))
	}

	// Maxinum Defined
	{
		s := String{MaxLength: 10}
		assert.Nil(t, s.Validate("this is ok"))
		assert.NotNil(t, s.Validate("this is a really long string and it should fail becuase it's too long."))
	}
}

func TestStringEnum(t *testing.T) {

	s := String{
		Enum: []string{"Joseph", "Simon", "Sara", "Mary"},
	}

	require.Nil(t, s.Validate("Joseph"))
	require.Nil(t, s.Validate("Simon"))
	require.Nil(t, s.Validate("Sara"))
	require.Nil(t, s.Validate("Mary"))
	require.NotNil(t, s.Validate("Mr. Black"))
	require.NotNil(t, s.Validate(3.14159265358979323846))
}
