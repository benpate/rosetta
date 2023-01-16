package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringUnmarshalSimple(t *testing.T) {

	value := []byte(`{"type":"string", "minLength":10, "maxLength":100}`)

	st, err := UnmarshalJSON(value)
	require.Nil(t, err)

	str := st.(String)
	require.Equal(t, str.MinLength, 10)
	require.Equal(t, str.MaxLength, 100)
}

func TestStringUnmarshalComplete(t *testing.T) {

	value := []byte(`{"type":"string", "format":"date", "pattern":"abc123", "minLength":10, "maxLength":100, "required":true}`)

	st, err := UnmarshalJSON(value)

	require.Nil(t, err)

	str := st.(String)
	require.Equal(t, str.MinLength, 10)
	require.Equal(t, str.MaxLength, 100)
	require.Equal(t, str.Required, true)
	require.Equal(t, str.Format, "date")
	require.Equal(t, str.Pattern, "abc123")
}
