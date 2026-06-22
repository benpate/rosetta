package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringType(t *testing.T) {

	s := String{}

	require.Nil(t, s.Validate("I'm a string"))
	require.NotNil(t, s.Validate(0))
	require.NotNil(t, s.Validate([]string{}))
	require.NotNil(t, s.Validate(map[string]any{}))
}

func TestStringUnmarshal(t *testing.T) {

	j := []byte(`{"type":"string", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(String).Required)

	_, _, err = Validate(s, "this should work")
	require.Nil(t, err)

	_, _, err = Validate(s, "")

	require.Error(t, err)
}

func TestStringRequired(t *testing.T) {

	s := String{Required: true}

	_, _, err := validate(s, "")
	require.Error(t, err)

	_, _, err = validate(s, "present")
	require.Nil(t, err)
}

func TestStringLength(t *testing.T) {

	// No Min/Max Length Defined
	{
		s := String{}
		_, _, err := validate(s, "")
		require.Nil(t, err)

		_, _, err = validate(s, "ok.")
		require.Nil(t, err)

		_, _, err = validate(s, "this is a really long string but it should be ok.")
		require.Nil(t, err)
	}

	// Mininum Defined
	{
		s := String{MinLength: 10}
		_, _, err := validate(s, "this is ok, because it's more than the minimum.")
		require.Nil(t, err)

		_, _, err = validate(s, "error")
		require.Error(t, err)
	}

	// Within Maximum
	{
		s := String{MaxLength: 10}

		newValue, changed, err := validate(s, "this is ok")
		require.Nil(t, err)
		require.False(t, changed)
		require.Equal(t, "this is ok", newValue)
	}

	// Exceeds Maximum (rewrite value)
	{
		s := String{MaxLength: 10}

		newValue, changed, err := validate(s, "1234567890 - this is a really long string and it should fail because it's too long.")
		require.NoError(t, err)
		require.True(t, changed)
		require.Equal(t, "1234567890", newValue)
	}
}

func TestStringEnum(t *testing.T) {

	s := String{
		Enum: []string{"Joseph", "Simon", "Sara", "Mary"},
	}

	_, changed, err := validate(s, "Joseph")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "Simon")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "Sara")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "Mary")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "Mr. Black")
	require.Error(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, 3.14159265358979323846)
	require.Error(t, err)
	require.False(t, changed)
}

func TestStringEnumRequired(t *testing.T) {

	s := String{
		Required: true,
		Enum:     []string{"Joseph", "Simon", "Sara", "Mary"},
	}

	require.NotNil(t, s.Validate(""))
}

func TestStringMinValue(t *testing.T) {
	s := String{MinValue: "abc"}

	_, _, err := validate(s, "abcd")
	require.Nil(t, err)
	require.NotNil(t, s.Validate("ab"))
	require.NotNil(t, s.Validate("a"))
	require.NotNil(t, s.Validate(""))
}

func TestStringMaxValue(t *testing.T) {
	s := String{MaxValue: "abc"}

	_, changed, err := validate(s, "ab")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "a")
	require.Nil(t, err)
	require.False(t, changed)

	_, changed, err = validate(s, "abcd")
	require.Error(t, err)
	require.False(t, changed)
}
