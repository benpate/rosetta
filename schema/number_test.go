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

	require.NotNil(t, s.Validate(-1.0))
	require.NotNil(t, s.Validate(0.0))
	require.Nil(t, s.Validate(2.0))
	require.Nil(t, s.Validate(4.0))
	require.Nil(t, s.Validate(6.0))
	require.Nil(t, s.Validate(8.0))
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
	s := Number{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Required)

	require.Nil(t, s.Validate(10.1))
	require.Nil(t, s.Validate(20.0))

	require.NotNil(t, s.Validate(0.0))
}

func TestNumberMinimum(t *testing.T) {

	s := Number{Minimum: null.NewFloat(10)}

	require.NotNil(t, s.Validate(1.0))
	require.Nil(t, s.Validate(10.0))
	require.Nil(t, s.Validate(10000.0))
}

func TestNumberMaximum(t *testing.T) {

	s := Number{Maximum: null.NewFloat(10)}

	require.Nil(t, s.Validate(1.0))
	require.Nil(t, s.Validate(5.0))
	require.Nil(t, s.Validate(10.0))
	require.NotNil(t, s.Validate(10000.0))
}
