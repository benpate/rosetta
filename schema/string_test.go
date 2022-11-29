package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/assert"
)

func TestStringUnmarshalSimple(t *testing.T) {

	value := []byte(`{"type":"string", "minLength":10, "maxLength":100}`)

	st, err := UnmarshalJSON(value)
	assert.Nil(t, err)

	str := st.(String)
	assert.Equal(t, str.MinLength, null.NewInt(10))
	assert.Equal(t, str.MaxLength, null.NewInt(100))
}

func TestStringUnmarshalComplete(t *testing.T) {

	value := []byte(`{"type":"string", "format":"date", "pattern":"abc123", "minLength":10, "maxLength":100, "required":true}`)

	st, err := UnmarshalJSON(value)

	assert.Nil(t, err)

	str := st.(String)
	assert.Equal(t, str.MinLength, null.NewInt(10))
	assert.Equal(t, str.MaxLength, null.NewInt(100))
	assert.Equal(t, str.Required, true)
	assert.Equal(t, str.Format, "date")
	assert.Equal(t, str.Pattern, "abc123")
}

/*
func TestStringFormatLowercase(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "format":"lowercase=2"}`))

	require.Nil(t, err)

	require.NotNil(t, s.Validate("NOT-ENOUGH-LOWERCASE"))
	require.NotNil(t, s.Validate("NOT-ENOUGH-LOWERCASE-a"))
	require.Nil(t, s.Validate("ENOUGH-LOWERCASE-ab"))
}

func TestStringFormatUppercase(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "format":"uppercase=2"}`))

	require.Nil(t, err)

	require.NotNil(t, s.Validate("not-enough-uppercase"))
	require.NotNil(t, s.Validate("not-enough-uppercase-A"))
	require.Nil(t, s.Validate("enough-uppercase-AB"))
}
*/
