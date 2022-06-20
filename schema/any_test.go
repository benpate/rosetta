package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {

	j := []byte(`{"type":"any"}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.Nil(t, s.Validate(0))
	require.Nil(t, s.Validate(0.1))
	require.Nil(t, s.Validate("hello there"))
	require.Nil(t, s.Validate("general kenobi"))

	{
		serialized, err := json.Marshal(s)

		require.Nil(t, err)
		require.Equal(t, j, serialized)
	}
}

func TestAnyRequired(t *testing.T) {

	j := []byte(`{"type":"any", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(*Any).Required)

	require.Nil(t, s.Validate("any string"))
	require.Nil(t, s.Validate(10))
	require.Nil(t, s.Validate(10.1))
	require.Nil(t, s.Validate(true))

	require.NotNil(t, s.Validate(""))
	require.NotNil(t, s.Validate(nil))
}

func TestAnyUnmarshal(t *testing.T) {

	a := Any{}

	{
		d := map[string]interface{}{
			"type": "any",
		}
		require.Nil(t, a.UnmarshalMap(d))
	}

	{
		d := map[string]interface{}{
			"type": "error",
		}
		require.NotNil(t, a.UnmarshalMap(d))
	}

}
