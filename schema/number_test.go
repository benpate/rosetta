package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestNumber(t *testing.T) {

	s := Number{
		Minimum: null.NewFloat(1),
		Maximum: null.NewFloat(10),
	}

	require.NotNil(t, s.Validate(-1))
	require.NotNil(t, s.Validate(0))
	require.Nil(t, s.Validate(2))
	require.Nil(t, s.Validate(4))
	require.Nil(t, s.Validate(6))
	require.Nil(t, s.Validate(8))
}

func TestNumberEnum(t *testing.T) {

	s := Number{
		Enum: []float64{1, 2, 3},
	}

	require.Nil(t, s.Validate(1))
	require.Nil(t, s.Validate(2))
	require.Nil(t, s.Validate(3))
	require.NotNil(t, s.Validate(4))
	require.NotNil(t, s.Validate("hamburger"))
}

func TestNumberRequired(t *testing.T) {

	j := []byte(`{"type":"number", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(Number).Required)

	require.Nil(t, s.Validate(10.1))
	require.Nil(t, s.Validate(20.0))

	require.NotNil(t, s.Validate(0.0))
}
