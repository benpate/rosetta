package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringType(t *testing.T) {

	s := String{}

	assert.Nil(t, s.Validate("I'm a string"))
	assert.Nil(t, s.Validate(0))
	assert.Nil(t, s.Validate([]string{}))
	assert.Nil(t, s.Validate(map[string]interface{}{}))
}

func TestStringUnmarshal(t *testing.T) {

	j := []byte(`{"type":"string", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(*String).Required)

	require.Nil(t, s.Validate("this should work"))
	require.NotNil(t, s.Validate(""))
}

func TestStringRequired(t *testing.T) {

	// Required schema
	{
		s := String{Required: true}

		assert.Nil(t, s.Validate("present"))
		assert.NotNil(t, s.Validate(""))
	}
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
		s := String{MinLength: null.NewInt(10)}
		assert.Nil(t, s.Validate("this is ok, becuase it's more than the minimum."))
		assert.NotNil(t, s.Validate("error"))
	}

	// Maxinum Defined
	{
		s := String{MaxLength: null.NewInt(10)}
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

func TestStringEnumUnmarshal(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "enum":["John", "Sarah", "Kyle"]}`))

	require.Nil(t, err)

	require.Nil(t, s.Validate("John"))
	require.Nil(t, s.Validate("Sarah"))
	require.Nil(t, s.Validate("Kyle"))
	require.NotNil(t, s.Validate("Anyone Else"))
}
