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

func TestStringFormatURL(t *testing.T) {

	uriSchema, err := UnmarshalJSON([]byte(`{"type":"string", "format":"uri"}`))
	require.Nil(t, err)

	urlSchema, err := UnmarshalJSON([]byte(`{"type":"string", "format":"url"}`))
	require.Nil(t, err)

	// Both formats accept an absolute web URL
	_, _, err = validate(uriSchema, "https://example.com/path")
	require.Nil(t, err)

	_, _, err = validate(urlSchema, "https://example.com/path")
	require.Nil(t, err)

	// "url" is registered as the stricter validator: an opaque mailto: address
	// has a scheme but no host, so "uri" accepts it and "url" rejects it
	_, _, err = validate(uriSchema, "mailto:user@example.com")
	require.Nil(t, err)

	_, _, err = validate(urlSchema, "mailto:user@example.com")
	require.NotNil(t, err)

	// Neither format accepts a relative reference
	_, _, err = validate(uriSchema, "example.com/path")
	require.NotNil(t, err)

	_, _, err = validate(urlSchema, "example.com/path")
	require.NotNil(t, err)
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
