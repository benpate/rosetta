package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestNumber_DefaultValue_BitSizes(t *testing.T) {
	def := null.NewFloat(5)

	require.Equal(t, float32(5), Number{BitSize: 32, Default: def}.DefaultValue())
	require.Equal(t, float64(5), Number{BitSize: 64, Default: def}.DefaultValue())
	require.Equal(t, float64(5), Number{Default: def}.DefaultValue())
}

func TestNumber_IsRequired(t *testing.T) {
	require.False(t, Number{}.IsRequired())
	require.True(t, Number{Required: true}.IsRequired())
}

func TestNumber_Validate_OK(t *testing.T) {
	require.NoError(t, Number{}.Validate(3.14))
	require.NoError(t, Number{Minimum: null.NewFloat(1), Maximum: null.NewFloat(10)}.Validate(5.0))
	require.NoError(t, Number{MultipleOf: null.NewFloat(2)}.Validate(8.0))
	require.NoError(t, Number{Enum: []float64{1.5, 2.5}}.Validate(2.5))
}

func TestNumber_Validate_NotFloat(t *testing.T) {
	require.Error(t, Number{}.Validate("not-a-float"))
}

func TestNumber_Validate_Required(t *testing.T) {
	require.Error(t, Number{Required: true}.Validate(0.0))
}

func TestNumber_Validate_Minimum(t *testing.T) {
	require.Error(t, Number{Minimum: null.NewFloat(10)}.Validate(5.0))
}

func TestNumber_Validate_Maximum(t *testing.T) {
	require.Error(t, Number{Maximum: null.NewFloat(10)}.Validate(15.0))
}

func TestNumber_Validate_MultipleOf(t *testing.T) {
	require.Error(t, Number{MultipleOf: null.NewFloat(2)}.Validate(7.0))
}

func TestNumber_Validate_Enum(t *testing.T) {
	require.Error(t, Number{Enum: []float64{1.5, 2.5}}.Validate(9.9))
}

func TestNumber_GetElement(t *testing.T) {
	element := Number{}

	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)

	_, ok = element.GetElement("nope")
	require.False(t, ok)
}

func TestNumber_Inherit(t *testing.T) {
	require.NotPanics(t, func() { Number{}.Inherit(String{}) })
}

func TestNumber_AllProperties(t *testing.T) {
	element := Number{}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestNumber_Enumerate(t *testing.T) {
	require.Equal(t, []string{"1.5", "2.5"}, Number{Enum: []float64{1.5, 2.5}}.Enumerate())
}

func TestNumber_MarshalMap(t *testing.T) {
	result := Number{
		Default:  null.NewFloat(1),
		Minimum:  null.NewFloat(2),
		Maximum:  null.NewFloat(3),
		Enum:     []float64{4, 5},
		Required: true,
	}.MarshalMap()

	require.Equal(t, TypeNumber, result["type"])
	require.Equal(t, float64(1), result["default"])
	require.Equal(t, float64(2), result["minimum"])
	require.Equal(t, float64(3), result["maximum"])
	require.Equal(t, []float64{4, 5}, result["enum"])
	require.Equal(t, true, result["required"])
}

func TestNumber_MarshalMap_RequiredIf(t *testing.T) {
	result := Number{RequiredIf: "a is b"}.MarshalMap()
	require.Equal(t, "a is b", result["required-if"])
}

func TestNumber_UnmarshalMap(t *testing.T) {
	element := Number{}
	err := element.UnmarshalMap(map[string]any{
		"type":        "number",
		"default":     1.0,
		"minimum":     2.0,
		"maximum":     3.0,
		"enum":        []float64{4, 5},
		"required":    true,
		"required-if": "a is b",
	})

	require.NoError(t, err)
	require.Equal(t, float64(1), element.Default.Float())
	require.Equal(t, float64(2), element.Minimum.Float())
	require.Equal(t, float64(3), element.Maximum.Float())
	require.Equal(t, []float64{4, 5}, element.Enum)
	require.True(t, element.Required)
	require.Equal(t, "a is b", element.RequiredIf)
}

func TestNumber_UnmarshalMap_WrongType(t *testing.T) {
	element := Number{}
	require.Error(t, element.UnmarshalMap(map[string]any{"type": "string"}))
}

func TestToFloat(t *testing.T) {
	cases := []any{int(1), int8(1), int16(1), int32(1), int64(1), float32(1), float64(1)}

	for _, value := range cases {
		result, ok := toFloat(value)
		require.True(t, ok)
		require.Equal(t, float64(1), result)
	}
}

func TestToFloat_NotANumber(t *testing.T) {
	result, ok := toFloat("nope")
	require.False(t, ok)
	require.Zero(t, result)
}
