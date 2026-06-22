package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestInteger_DefaultValue_BitSizes(t *testing.T) {
	def := null.NewInt64(5)

	require.Equal(t, int8(5), Integer{BitSize: 8, Default: def}.DefaultValue())
	require.Equal(t, int16(5), Integer{BitSize: 16, Default: def}.DefaultValue())
	require.Equal(t, int32(5), Integer{BitSize: 32, Default: def}.DefaultValue())
	require.Equal(t, int64(5), Integer{BitSize: 64, Default: def}.DefaultValue())
	require.Equal(t, int(5), Integer{Default: def}.DefaultValue())
}

func TestInteger_IsRequired(t *testing.T) {
	require.False(t, Integer{}.IsRequired())
	require.True(t, Integer{Required: true}.IsRequired())
}

func TestInteger_Validate_OK(t *testing.T) {
	require.NoError(t, Integer{}.Validate(42))
	require.NoError(t, Integer{Minimum: null.NewInt64(10), Maximum: null.NewInt64(100)}.Validate(50))
	require.NoError(t, Integer{MultipleOf: null.NewInt64(5)}.Validate(25))
	require.NoError(t, Integer{Enum: []int{1, 2, 3}}.Validate(2))
}

func TestInteger_Validate_NotInteger(t *testing.T) {
	require.Error(t, Integer{}.Validate("not-an-int"))
}

func TestInteger_Validate_Required(t *testing.T) {
	require.Error(t, Integer{Required: true}.Validate(0))
}

func TestInteger_Validate_Minimum(t *testing.T) {
	require.Error(t, Integer{Minimum: null.NewInt64(10)}.Validate(5))
}

func TestInteger_Validate_Maximum(t *testing.T) {
	require.Error(t, Integer{Maximum: null.NewInt64(10)}.Validate(15))
}

func TestInteger_Validate_MultipleOf(t *testing.T) {
	require.Error(t, Integer{MultipleOf: null.NewInt64(5)}.Validate(23))
}

func TestInteger_Validate_Enum(t *testing.T) {
	require.Error(t, Integer{Enum: []int{1, 2, 3}}.Validate(9))
}

func TestInteger_GetElement(t *testing.T) {
	element := Integer{}

	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)

	_, ok = element.GetElement("nope")
	require.False(t, ok)
}

func TestInteger_Inherit(t *testing.T) {
	require.NotPanics(t, func() { Integer{}.Inherit(String{}) })
}

func TestInteger_AllProperties(t *testing.T) {
	element := Integer{}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestInteger_Enumerate(t *testing.T) {
	require.Equal(t, []string{"1", "2", "3"}, Integer{Enum: []int{1, 2, 3}}.Enumerate())
}

func TestInteger_MarshalMap(t *testing.T) {
	result := Integer{
		Default:    null.NewInt64(1),
		Minimum:    null.NewInt64(2),
		Maximum:    null.NewInt64(3),
		MultipleOf: null.NewInt64(4),
		Enum:       []int{5, 6},
	}.MarshalMap()

	require.Equal(t, TypeInteger, result["type"])
	require.Equal(t, int64(1), result["default"])
	require.Equal(t, int64(2), result["minimum"])
	require.Equal(t, int64(3), result["maximum"])
	require.Equal(t, int64(4), result["multipleOf"])
	require.Equal(t, []int{5, 6}, result["enum"])
}

func TestInteger_MarshalMap_Empty(t *testing.T) {
	result := Integer{}.MarshalMap()

	require.Equal(t, TypeInteger, result["type"])
	require.NotContains(t, result, "default")
	require.NotContains(t, result, "minimum")
}

func TestInteger_UnmarshalMap(t *testing.T) {
	element := Integer{}
	err := element.UnmarshalMap(map[string]any{
		"type":       "integer",
		"default":    1,
		"minimum":    2,
		"maximum":    3,
		"multipleOf": 4,
		"required":   true,
		"enum":       []int{5, 6},
	})

	require.NoError(t, err)
	require.Equal(t, int64(1), element.Default.Int64())
	require.Equal(t, int64(2), element.Minimum.Int64())
	require.Equal(t, int64(3), element.Maximum.Int64())
	require.Equal(t, int64(4), element.MultipleOf.Int64())
	require.True(t, element.Required)
	require.Equal(t, []int{5, 6}, element.Enum)
}

func TestInteger_UnmarshalMap_WrongType(t *testing.T) {
	element := Integer{}
	require.Error(t, element.UnmarshalMap(map[string]any{"type": "string"}))
}

func TestToInt64(t *testing.T) {
	cases := []any{int(1), int8(1), int16(1), int32(1), int64(1)}

	for _, value := range cases {
		result, ok := toInt64(value)
		require.True(t, ok)
		require.Equal(t, int64(1), result)
	}
}

func TestToInt64_NotAnInteger(t *testing.T) {
	result, ok := toInt64("nope")
	require.False(t, ok)
	require.Zero(t, result)
}
