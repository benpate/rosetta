package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInteger(t *testing.T) {

	s := Integer{}

	assert.Nil(t, s.Validate(0))
	assert.Nil(t, s.Validate(1))

	assert.NotNil(t, s.Validate(1.1))
	assert.NotNil(t, s.Validate("string-bad"))
}

func TestIntegerEnum(t *testing.T) {

	s := Integer{
		Enum: []int{1, 2, 3},
	}

	require.Nil(t, s.Validate(1))
	require.Nil(t, s.Validate(2))
	require.Nil(t, s.Validate(3))
	require.NotNil(t, s.Validate(4))
	require.NotNil(t, s.Validate("hamburger"))
}

func TestIntegerRequired(t *testing.T) {

	j := []byte(`{"type":"integer", "required":true}`)
	s := Schema{}

	err := json.Unmarshal(j, &s)
	require.Nil(t, err)

	require.True(t, s.Element.(Integer).Required)

	require.Nil(t, s.Validate(10))
	require.Nil(t, s.Validate(20))

	require.NotNil(t, s.Validate(0))
}

func TestIntegerMinimum(t *testing.T) {

	s := Integer{
		Minimum: null.NewInt64(1),
	}

	assert.NotNil(t, s.Validate(-1))
	assert.NotNil(t, s.Validate(0))
	assert.Nil(t, s.Validate(1))
	assert.Nil(t, s.Validate(2))
}

func TestIntegerMaximum(t *testing.T) {

	s := Integer{
		Maximum: null.NewInt64(5),
	}

	assert.Nil(t, s.Validate(1))
	assert.Nil(t, s.Validate(2))
	assert.Nil(t, s.Validate(3))
	assert.Nil(t, s.Validate(4))
	assert.Nil(t, s.Validate(5))
	assert.NotNil(t, s.Validate(6))
	assert.NotNil(t, s.Validate(7))
}

func TestIntegerMultipleOf(t *testing.T) {

	s := Integer{
		MultipleOf: null.NewInt64(3),
	}

	assert.NotNil(t, s.Validate(-1))
	assert.Nil(t, s.Validate(0))
	assert.NotNil(t, s.Validate(1))
	assert.NotNil(t, s.Validate(2))
	assert.Nil(t, s.Validate(3))
	assert.NotNil(t, s.Validate(4))
	assert.NotNil(t, s.Validate(5))
	assert.Nil(t, s.Validate(6))
	assert.NotNil(t, s.Validate(7))
	assert.NotNil(t, s.Validate(8))
	assert.Nil(t, s.Validate(9))
	assert.NotNil(t, s.Validate(10))
}
