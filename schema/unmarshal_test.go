package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal_Success(t *testing.T) {

	s := Unmarshal(`{"$id":"TEST-SCHEMA", "$comment":"Test Unmarshal", "type":"string", "maxLength":10}`)

	require.Equal(t, "TEST-SCHEMA", s.ID)
	require.Equal(t, "Test Unmarshal", s.Comment)

	element := s.Element.(String)
	require.True(t, element.MaxLength.IsPresent())
	require.Equal(t, element.MaxLength.Int(), 10)

	require.False(t, element.MinLength.IsPresent())
}

func TestUnmarshal_Error(t *testing.T) {

	s := Unmarshal("this will not work")

	require.Empty(t, s.ID)
	require.Empty(t, s.Comment)
	require.Nil(t, s.Element)
}

func TestManualUnmarshal(t *testing.T) {
	schema := getTestSchema()
	assert.NotNil(t, schema)
}
