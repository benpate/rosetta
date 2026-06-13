package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObject_GetElement_Self(t *testing.T) {
	element := Object{}
	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)
}

func TestObject_GetElement_Property(t *testing.T) {
	element := Object{Properties: ElementMap{"name": String{}}}

	found, ok := element.GetElement("name")
	require.True(t, ok)
	require.Equal(t, String{}, found)
}

func TestObject_GetElement_Nested(t *testing.T) {
	element := Object{Properties: ElementMap{
		"child": Object{Properties: ElementMap{"name": String{}}},
	}}

	found, ok := element.GetElement("child.name")
	require.True(t, ok)
	require.Equal(t, String{}, found)
}

func TestObject_GetElement_Wildcard(t *testing.T) {
	element := Object{Wildcard: String{}}

	found, ok := element.GetElement("anything")
	require.True(t, ok)
	require.Equal(t, String{}, found)
}

func TestObject_GetElement_NotFound(t *testing.T) {
	element := Object{Properties: ElementMap{"name": String{}}}

	_, ok := element.GetElement("missing")
	require.False(t, ok)
}

func TestObject_DefaultValue(t *testing.T) {
	element := Object{Properties: ElementMap{
		"name":   String{Default: "x"},
		"active": Boolean{},
	}}

	result, ok := element.DefaultValue().(map[string]any)
	require.True(t, ok)
	require.Equal(t, "x", result["name"])
	require.Equal(t, false, result["active"])
}

func TestObject_IsRequired(t *testing.T) {
	require.False(t, Object{}.IsRequired())
	require.True(t, Object{Required: true}.IsRequired())
}

func TestObject_Inherit_AddsParentProperties(t *testing.T) {
	child := Object{Properties: ElementMap{"a": String{}}}
	parent := Object{Properties: ElementMap{"a": String{}, "b": Integer{}}}

	child.Inherit(parent)

	// The "b" property is inherited from the parent
	require.Contains(t, child.Properties, "b")
	require.Equal(t, Integer{}, child.Properties["b"])
}

func TestObject_Inherit_NilProperties(t *testing.T) {
	// Inheriting into an object with no properties must not panic
	require.NotPanics(t, func() {
		Object{}.Inherit(Object{Properties: ElementMap{"a": String{}}})
	})
}

func TestObject_AllProperties_Flattens(t *testing.T) {
	element := Object{Properties: ElementMap{
		"name":  String{},
		"child": Object{Properties: ElementMap{"age": Integer{}}},
	}}

	result := element.AllProperties()

	require.Equal(t, String{}, result["name"])
	require.Equal(t, Integer{}, result["child.age"])
}

func TestObject_MarshalMap(t *testing.T) {
	element := Object{
		Properties: ElementMap{"name": String{}},
		Wildcard:   String{},
		Required:   true,
		RequiredIF: "a is b",
	}

	result := element.MarshalMap()

	require.Equal(t, TypeObject, result["type"])
	require.Equal(t, true, result["required"])
	require.Equal(t, "a is b", result["required-if"])
	require.Contains(t, result, "properties")
	require.Contains(t, result, "wildcard")
}

func TestObject_UnmarshalMap(t *testing.T) {
	element := Object{}
	err := element.UnmarshalMap(map[string]any{
		"type":     "object",
		"required": true,
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
		"wildcard": map[string]any{"type": "string"},
	})

	require.NoError(t, err)
	require.True(t, element.Required)
	require.Contains(t, element.Properties, "name")
	require.NotNil(t, element.Wildcard)
}

func TestObject_UnmarshalMap_WrongType(t *testing.T) {
	element := Object{}
	require.Error(t, element.UnmarshalMap(map[string]any{"type": "string"}))
}
