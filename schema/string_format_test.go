package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringUnmarshalSimple1(t *testing.T) {

	value := []byte(`{"type":"string", "minLength":10, "maxLength":100}`)

	st, err := UnmarshalJSON(value)
	assert.Nil(t, err)

	if err != nil {
		return
	}

	str := st.(String)
	assert.Equal(t, str.MinLength, 10)
	assert.Equal(t, str.MaxLength, 100)
}

func TestStringUnmarshalComplete1(t *testing.T) {

	value := []byte(`{"type":"string", "format":"date", "pattern":"abc123", "minLength":10, "maxLength":100, "required":true}`)

	st, err := UnmarshalJSON(value)

	assert.Nil(t, err)

	if err != nil {
		return
	}

	str := st.(String)
	assert.Equal(t, str.MinLength, 10)
	assert.Equal(t, str.MaxLength, 100)
	assert.Equal(t, str.Required, true)
	assert.Equal(t, str.Format, "date")
}

func TestStringFormatLowercase1(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "format":"lower=2"}`))

	require.Nil(t, err)

	if err != nil {
		return
	}

	_, _, err = validate(s, "NOT-ENOUGH-LOWERCASE")
	require.NotNil(t, err)

	_, _, err = validate(s, "NOT-ENOUGH-LOWERCASE-a")
	require.NotNil(t, err)

	_, _, err = validate(s, "ENOUGH-LOWERCASE-ab")
	require.Nil(t, err)
}

func TestStringFormatUppercase1(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "format":"upper=2"}`))

	require.Nil(t, err)

	if err != nil {
		return
	}

	_, _, err = validate(s, "not-enough-uppercase")
	require.NotNil(t, err)

	_, _, err = validate(s, "not-enough-uppercase-A")
	require.NotNil(t, err)

	_, _, err = validate(s, "enough-uppercase-AB")
	require.Nil(t, err)
}
