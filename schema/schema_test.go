package schema

import (
	"testing"

	"github.com/benpate/exp"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestSchemaMatch(t *testing.T) {

	value := newTestStructA()
	schema := New(testStructA_Schema())

	match, err := schema.Match(value, exp.Predicate{
		Field:    "active",
		Operator: exp.OperatorEqual,
		Value:    true,
	})

	require.NoError(t, err)
	require.True(t, match)

	match, err = schema.Match(value, exp.Predicate{
		Field:    "name",
		Operator: exp.OperatorEqual,
		Value:    "John Connor",
	})
	require.NoError(t, err)
	require.True(t, match)

	match, err = schema.Match(value, exp.Predicate{
		Field:    "name",
		Operator: exp.OperatorNotEqual,
		Value:    "Sarah Connor",
	})
	require.True(t, match)
	require.NoError(t, err)
}

func TestSchemaMatchError(t *testing.T) {

	value := newTestStructA()
	schema := New(testStructA_Schema())

	match, err := schema.Match(value, exp.Predicate{
		Field:    "missing_property",
		Operator: exp.OperatorEqual,
		Value:    "Sarah Connor",
	})
	require.Error(t, err)
	require.False(t, match)
}

func TestSchemaValidateRequiredIf(t *testing.T) {

	spew.Config.DisableMethods = true

	value := newTestStructA()
	schema := New(testStructA_Schema())

	{
		err := schema.Validate(&value)
		require.NoError(t, err)
	}

	{
		value.Name = "Aethelflad"
		err := schema.Validate(&value)
		require.Error(t, err)
	}
}
