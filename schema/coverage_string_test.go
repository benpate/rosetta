package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringElement_DefaultValue(t *testing.T) {
	require.Equal(t, "hello", String{Default: "hello"}.DefaultValue())
}

func TestStringElement_IsRequired(t *testing.T) {
	require.False(t, String{}.IsRequired())
	require.True(t, String{Required: true}.IsRequired())
}

func TestStringElement_Validate_OK(t *testing.T) {
	require.NoError(t, String{}.Validate("hello"))
	require.NoError(t, String{MinLength: 2, MaxLength: 10}.Validate("hello"))
	require.NoError(t, String{Enum: []string{"a", "b"}}.Validate("a"))
	require.NoError(t, String{MinValue: "a", MaxValue: "z"}.Validate("m"))
}

func TestStringElement_Validate_NotString(t *testing.T) {
	require.Error(t, String{}.Validate(42))
}

func TestStringElement_Validate_Required(t *testing.T) {
	require.Error(t, String{Required: true}.Validate(""))
}

func TestStringElement_Validate_MinValue(t *testing.T) {
	require.Error(t, String{MinValue: "m"}.Validate("a"))
}

func TestStringElement_Validate_MaxValue(t *testing.T) {
	require.Error(t, String{MaxValue: "m"}.Validate("z"))
}

func TestStringElement_Validate_MinLength(t *testing.T) {
	require.Error(t, String{MinLength: 10}.Validate("short"))
}

func TestStringElement_Validate_MaxLength(t *testing.T) {
	require.Error(t, String{MaxLength: 3}.Validate("too-long"))
}

func TestStringElement_Validate_Enum(t *testing.T) {
	require.Error(t, String{Enum: []string{"a", "b"}}.Validate("z"))
}

func TestStringElement_Validate_FormatError(t *testing.T) {
	// The "email" format rejects values that are not valid email addresses
	require.Error(t, String{Format: "email"}.Validate("not-an-email"))
}

func TestStringElement_GetElement(t *testing.T) {
	element := String{}

	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)

	_, ok = element.GetElement("nope")
	require.False(t, ok)
}

func TestStringElement_Inherit(t *testing.T) {
	require.NotPanics(t, func() { String{}.Inherit(Integer{}) })
}

func TestStringElement_AllProperties(t *testing.T) {
	element := String{}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestStringElement_Enumerate(t *testing.T) {
	require.Equal(t, []string{"a", "b"}, String{Enum: []string{"a", "b"}}.Enumerate())
}

func TestStringElement_MarshalMap(t *testing.T) {
	result := String{
		Default:    "x",
		MinLength:  1,
		MaxLength:  2,
		Format:     "email",
		Enum:       []string{"a", "b"},
		RequiredIf: "a is b",
	}.MarshalMap()

	require.Equal(t, TypeString, result["type"])
	require.Equal(t, "x", result["default"])
	require.Equal(t, 1, result["minLength"])
	require.Equal(t, 2, result["maxLength"])
	require.Equal(t, "email", result["format"])
	require.Equal(t, []string{"a", "b"}, result["enum"])
	require.Equal(t, "a is b", result["required-if"])
}

func TestStringElement_MarshalMap_Empty(t *testing.T) {
	result := String{}.MarshalMap()

	require.Equal(t, TypeString, result["type"])
	require.Equal(t, false, result["required"])
	require.NotContains(t, result, "default")
}

func TestStringElement_UnmarshalMap(t *testing.T) {
	element := String{}
	err := element.UnmarshalMap(map[string]any{
		"type":        "string",
		"default":     "x",
		"minLength":   1,
		"maxLength":   2,
		"format":      "email",
		"enum":        []string{"a", "b"},
		"required":    true,
		"required-if": "a is b",
	})

	require.NoError(t, err)
	require.Equal(t, "x", element.Default)
	require.Equal(t, 1, element.MinLength)
	require.Equal(t, 2, element.MaxLength)
	require.Equal(t, "email", element.Format)
	require.Equal(t, []string{"a", "b"}, element.Enum)
	require.True(t, element.Required)
	require.Equal(t, "a is b", element.RequiredIf)
}

func TestStringElement_UnmarshalMap_WrongType(t *testing.T) {
	element := String{}
	require.Error(t, element.UnmarshalMap(map[string]any{"type": "number"}))
}

func TestStringElement_FormatFunctions_DefaultsToNoHTML(t *testing.T) {
	// With no format defined, the no-html format is applied (and accepts the value)
	require.NoError(t, String{}.Validate("<b>hello</b>"))
}
