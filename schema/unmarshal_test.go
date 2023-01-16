package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal_Success(t *testing.T) {

	var s Schema
	err := json.Unmarshal([]byte(`{"$id":"TEST-SCHEMA", "$comment":"Test Unmarshal", "type":"array", "items":{"type":"string"}, "maxLength":10}`), &s)

	require.Nil(t, err)
	require.Equal(t, "TEST-SCHEMA", s.ID)
	require.Equal(t, "Test Unmarshal", s.Comment)

	element := s.Element.(Array)
	require.Equal(t, element.MaxLength, 10)

	require.Zero(t, element.MinLength)
}

func TestUnmarshal_Error(t *testing.T) {

	var s Schema
	err := json.Unmarshal([]byte("this will not work"), &s)

	require.NotNil(t, err)
	require.Empty(t, s.ID)
	require.Empty(t, s.Comment)
	require.Nil(t, s.Element)
}

func TestStringEnumUnmarshal(t *testing.T) {

	s, err := UnmarshalJSON([]byte(`{"type":"string", "enum":["John", "Sarah", "Kyle"]}`))

	require.Nil(t, err)

	require.Nil(t, s.Validate("John"))
	require.Nil(t, s.Validate("Sarah"))
	require.Nil(t, s.Validate("Kyle"))
	require.NotNil(t, s.Validate("Anyone Else"))
}
