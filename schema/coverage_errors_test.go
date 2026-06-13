package schema

import (
	"testing"

	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

func TestSchema_Match_UnknownOperator(t *testing.T) {
	// An unrecognized operator never matches, but is not an error
	schema := New(testStructB_Schema())

	predicate := exp.Predicate{Field: "name", Operator: "not-a-real-operator", Value: "x"}
	match, err := schema.Match(newTestStructB(), predicate)

	require.NoError(t, err)
	require.False(t, match)
}

func TestSet_Boolean_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Boolean{}}})
	require.Error(t, schema.Set(unsupportedObject{}, "value", true))
}

func TestSet_Integer32_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{}}})
	require.Error(t, schema.Set(unsupportedObject{}, "value", 5))
}

func TestSet_Integer64_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{BitSize: 64}}})
	require.Error(t, schema.Set(unsupportedObject{}, "value", 5))
}

func TestSet_Number_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Number{}}})
	require.Error(t, schema.Set(unsupportedObject{}, "value", 1.5))
}

func TestGet_Integer32_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{}}})
	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)
}
