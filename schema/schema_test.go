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

	require.True(t, schema.Match(value, exp.Predicate{
		Field:    "active",
		Operator: exp.OperatorEqual,
		Value:    true,
	}))

	require.True(t, schema.Match(value, exp.Predicate{
		Field:    "name",
		Operator: exp.OperatorEqual,
		Value:    "John Connor",
	}))

	require.True(t, schema.Match(value, exp.Predicate{
		Field:    "name",
		Operator: exp.OperatorNotEqual,
		Value:    "Sarah Connor",
	}))
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
