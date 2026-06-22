package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny_DefaultValue(t *testing.T) {
	require.Nil(t, Any{}.DefaultValue())
}

func TestAny_IsRequired(t *testing.T) {
	require.False(t, Any{}.IsRequired())
	require.True(t, Any{Required: true}.IsRequired())
}

func TestAny_Validate(t *testing.T) {
	// Any accepts every value, so Validate never returns an error
	require.NoError(t, Any{}.Validate(nil))
	require.NoError(t, Any{}.Validate("anything"))
	require.NoError(t, Any{Required: true}.Validate(42))
}

func TestAny_GetElement(t *testing.T) {
	element := Any{}

	// Any returns itself for every path
	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)

	found, ok = element.GetElement("anything")
	require.True(t, ok)
	require.Equal(t, element, found)
}

func TestAny_Inherit(t *testing.T) {
	// Inherit is a no-op; it must not panic
	require.NotPanics(t, func() { Any{}.Inherit(String{}) })
}

func TestAny_AllProperties(t *testing.T) {
	element := Any{}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestAny_ValidateRequiredIf_Empty(t *testing.T) {
	// With no condition, ValidateRequiredIf always succeeds
	schema := New(Any{})
	require.NoError(t, schema.ValidateRequiredIf(nil))
}

func TestAny_MarshalMap(t *testing.T) {
	result := Any{Required: true, RequiredIf: "x is y"}.MarshalMap()

	require.Equal(t, TypeAny, result["type"])
	require.Equal(t, true, result["required"])
	require.Equal(t, "x is y", result["required-if"])
}

func TestAny_UnmarshalMap(t *testing.T) {
	element := Any{}
	err := element.UnmarshalMap(map[string]any{"type": "any", "required": true})

	require.NoError(t, err)
	require.True(t, element.Required)
}

func TestAny_UnmarshalMap_WrongType(t *testing.T) {
	element := Any{}
	err := element.UnmarshalMap(map[string]any{"type": "string"})

	require.Error(t, err)
}
