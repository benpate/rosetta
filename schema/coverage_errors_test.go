package schema

import (
	"testing"

	"github.com/benpate/derp"
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

// TestGet_String_UnsupportedProperty demonstrates the improved error messaging:
// when the object DOES implement the required getter interfaces but simply doesn't
// back the requested property, the error must say the *property* is unsupported --
// NOT that the object is the wrong type.  (This is the "title on model.Stream" case.)
func TestGet_String_UnsupportedProperty(t *testing.T) {

	// testStructB implements StringGetter AND PointerGetter, but only backs
	// "name" as a string -- it has no "title" property at all.
	schema := New(Object{Properties: ElementMap{"title": String{}}})

	_, err := schema.Get(newTestStructB(), "title")
	require.Error(t, err)

	// Dig to the leaf error that actually classified the failure.
	root := derp.RootCause(err)

	// The interface IS satisfied, so we must NOT report a type mismatch...
	require.NotContains(t, derp.Message(root), "must be a StringGetter")

	// ...instead we report that this specific property is unsupported...
	require.Equal(t, "Object does not support this string property", derp.Message(root))

	// ...and name the offending property in the details.
	require.Contains(t, derp.Details(root), "title")
}

// TestGet_String_UnsupportedType is the contrasting case: an object that implements
// NONE of the required interfaces still reports the original type-mismatch message.
func TestGet_String_UnsupportedType(t *testing.T) {

	schema := New(Object{Properties: ElementMap{"value": String{}}})

	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)

	// The object is genuinely the wrong type, so the original message stands.
	require.Equal(t, "Object must be a StringGetter or PointerGetter", derp.Message(derp.RootCause(err)))
}
