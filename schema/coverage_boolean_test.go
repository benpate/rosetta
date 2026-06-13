package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestBoolean_DefaultValue(t *testing.T) {
	require.Equal(t, false, Boolean{}.DefaultValue())
	require.Equal(t, true, Boolean{Default: null.NewBool(true)}.DefaultValue())
}

func TestBoolean_IsRequired(t *testing.T) {
	require.False(t, Boolean{}.IsRequired())
	require.True(t, Boolean{Required: true}.IsRequired())
}

func TestBoolean_Validate(t *testing.T) {
	require.NoError(t, Boolean{}.Validate(true))
	require.NoError(t, Boolean{}.Validate(false))
	require.NoError(t, Boolean{Required: true}.Validate(true))
}

func TestBoolean_Validate_NotBoolean(t *testing.T) {
	require.Error(t, Boolean{}.Validate("not-a-bool"))
}

func TestBoolean_Validate_RequiredFalse(t *testing.T) {
	// A required boolean cannot be FALSE
	require.Error(t, Boolean{Required: true}.Validate(false))
}

func TestBoolean_GetElement(t *testing.T) {
	element := Boolean{}

	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)

	found, ok = element.GetElement("nope")
	require.False(t, ok)
	require.Nil(t, found)
}

func TestBoolean_Inherit(t *testing.T) {
	require.NotPanics(t, func() { Boolean{}.Inherit(String{}) })
}

func TestBoolean_AllProperties(t *testing.T) {
	element := Boolean{}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestBoolean_MarshalMap(t *testing.T) {
	result := Boolean{Default: null.NewBool(true), Required: true, RequiredIf: "a is b"}.MarshalMap()

	require.Equal(t, TypeBoolean, result["type"])
	require.Equal(t, true, result["default"])
	require.Equal(t, true, result["required"])
	require.Equal(t, "a is b", result["required-if"])
}

func TestBoolean_MarshalMap_Empty(t *testing.T) {
	// An empty Boolean only reports its type
	result := Boolean{}.MarshalMap()

	require.Equal(t, TypeBoolean, result["type"])
	require.NotContains(t, result, "default")
	require.NotContains(t, result, "required")
	require.NotContains(t, result, "required-if")
}

func TestBoolean_MarshalJSON(t *testing.T) {
	bytes, err := json.Marshal(Boolean{Required: true})

	require.NoError(t, err)
	require.JSONEq(t, `{"type":"boolean","required":true}`, string(bytes))
}

func TestBoolean_UnmarshalMap(t *testing.T) {
	element := Boolean{}
	err := element.UnmarshalMap(map[string]any{"type": "boolean", "default": true, "required": true})

	require.NoError(t, err)
	require.True(t, element.Default.Bool())
	require.True(t, element.Required)
}

func TestBoolean_UnmarshalMap_WrongType(t *testing.T) {
	element := Boolean{}
	require.Error(t, element.UnmarshalMap(map[string]any{"type": "string"}))
}
