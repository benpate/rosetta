package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {

	s := Boolean{}

	assert.Nil(t, s.Validate(true))
	assert.Nil(t, s.Validate(false))

	assert.NotNil(t, s.Validate(1))
	assert.NotNil(t, s.Validate("string-bad"))

}

func TestBool_Required(t *testing.T) {

	s := Boolean{Required: true}

	require.Nil(t, s.Validate(true))
	require.NotNil(t, s.Validate(false))
}

func TestBool_Marshal(t *testing.T) {

	b := Boolean{}

	result, err := json.Marshal(b)
	require.Nil(t, err)
	require.Equal(t, `{"type":"boolean"}`, string(result))
}

func TestBool_Unmarshal(t *testing.T) {

	j := []byte(`{"type":"object", "wildcard":{"type":"boolean", "required":true}}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(Object).Wildcard.(Boolean).Required)
}
