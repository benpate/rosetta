package schema

import (
	"testing"

	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

func TestSchema_Wildcard(t *testing.T) {
	schema := Wildcard()
	require.Equal(t, Any{}, schema.Element)
}

func TestSchema_GetElement_NilElement(t *testing.T) {
	schema := Schema{}
	_, ok := schema.GetElement("anything")
	require.False(t, ok)
}

func TestSchema_GetArrayElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"list": Array{Items: String{}}}})

	element, ok := schema.GetArrayElement("list")
	require.True(t, ok)
	require.Equal(t, String{}, element.Items)
}

func TestSchema_GetArrayElement_FromAny(t *testing.T) {
	// An "Any" element resolves to an array of Any
	schema := New(Any{})

	element, ok := schema.GetArrayElement("anything")
	require.True(t, ok)
	require.Equal(t, Any{}, element.Items)
}

func TestSchema_GetArrayElement_WrongType(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"name": String{}}})

	_, ok := schema.GetArrayElement("name")
	require.False(t, ok)
}

func TestSchema_GetBooleanElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"active": Boolean{}}})

	_, ok := schema.GetBooleanElement("active")
	require.True(t, ok)

	_, ok = schema.GetBooleanElement("missing")
	require.False(t, ok)
}

func TestSchema_GetIntegerElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"age": Integer{}}})

	_, ok := schema.GetIntegerElement("age")
	require.True(t, ok)

	_, ok = schema.GetIntegerElement("missing")
	require.False(t, ok)
}

func TestSchema_GetNumberElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"lat": Number{}}})

	_, ok := schema.GetNumberElement("lat")
	require.True(t, ok)

	_, ok = schema.GetNumberElement("missing")
	require.False(t, ok)
}

func TestSchema_GetObjectElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"child": Object{}}})

	_, ok := schema.GetObjectElement("child")
	require.True(t, ok)

	_, ok = schema.GetObjectElement("missing")
	require.False(t, ok)
}

func TestSchema_GetStringElement(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"name": String{}}})

	_, ok := schema.GetStringElement("name")
	require.True(t, ok)

	_, ok = schema.GetStringElement("missing")
	require.False(t, ok)
}

func TestSchema_Inherit_NilElement(t *testing.T) {
	// Inheriting into an empty schema adopts the parent's element
	child := Schema{}
	parent := New(String{})

	child.Inherit(parent)
	require.Equal(t, String{}, child.Element)
}

func TestSchema_Inherit_ExistingElement(t *testing.T) {
	child := New(Object{Properties: ElementMap{"a": String{}}})
	parent := New(Object{Properties: ElementMap{"b": Integer{}}})

	child.Inherit(parent)

	object := child.Element.(Object)
	require.Contains(t, object.Properties, "b")
}

func TestSchema_AllProperties(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"name": String{}}})
	require.Equal(t, String{}, schema.AllProperties()["name"])
}

func TestSchema_Match_Operators(t *testing.T) {
	schema := New(testStructB_Schema())
	value := newTestStructB()

	cases := []struct {
		operator string
		field    string
		value    any
		expected bool
	}{
		{exp.OperatorEqual, "name", value.Name, true},
		{exp.OperatorNotEqual, "name", "different", true},
		{exp.OperatorGreaterThan, "age", value.Age - 1, true},
		{exp.OperatorLessThan, "age", value.Age + 1, true},
		{exp.OperatorGreaterOrEqual, "age", value.Age, true},
		{exp.OperatorLessOrEqual, "age", value.Age, true},
	}

	for _, c := range cases {
		predicate := exp.Predicate{Field: c.field, Operator: c.operator, Value: c.value}
		match, err := schema.Match(value, predicate)
		require.NoError(t, err)
		require.Equal(t, c.expected, match, "operator %s", c.operator)
	}
}

func TestSchema_Match_UnknownField(t *testing.T) {
	schema := New(testStructB_Schema())

	predicate := exp.Predicate{Field: "does-not-exist", Operator: exp.OperatorEqual, Value: "x"}
	match, err := schema.Match(newTestStructB(), predicate)

	require.Error(t, err)
	require.False(t, match)
}
