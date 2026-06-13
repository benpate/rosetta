package schema

import (
	"testing"

	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// unknownElement is an Element implementation that is not one of the built-in
// types, used to reach the "invalid element type" branches of the validators.
type unknownElement struct{}

func (unknownElement) DefaultValue() any                               { return nil }
func (unknownElement) IsRequired() bool                                { return false }
func (unknownElement) ValidateRequiredIf(Schema, list.List, any) error { return nil }
func (unknownElement) MarshalMap() map[string]any                      { return map[string]any{} }
func (unknownElement) GetElement(string) (Element, bool)               { return nil, false }
func (unknownElement) Inherit(Element)                                 {}
func (unknownElement) AllProperties() ElementMap                       { return ElementMap{} }

func TestValidate_NilSchema(t *testing.T) {
	_, _, err := Validate(Schema{}, "value")
	require.Error(t, err)
}

func TestValidate_Any(t *testing.T) {
	value, changed, err := validate(Any{}, "anything")
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, "anything", value)
}

func TestValidate_Boolean(t *testing.T) {
	_, _, err := validate(Boolean{}, true)
	require.NoError(t, err)
}

func TestValidate_Boolean_Required(t *testing.T) {
	_, _, err := validate(Boolean{Required: true}, false)
	require.Error(t, err)
}

func TestValidate_Boolean_CoercesInput(t *testing.T) {
	// A coercible value is accepted and converted to a boolean
	newValue, changed, err := validate(Boolean{}, "true")
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, true, newValue)
}

func TestValidate_Boolean_NotCoercible(t *testing.T) {
	// A value that cannot be coerced into a boolean is rejected by the dispatcher
	_, _, err := validate(Boolean{}, "not-a-bool")
	require.Error(t, err)
}

func TestValidateBoolean_Direct_NotBoolean(t *testing.T) {
	// Called directly with an un-coercible value, validate_Boolean reports an error
	_, _, err := validate_Boolean(Boolean{}, "xyz")
	require.Error(t, err)
}

func TestValidate_Integer_Clamps(t *testing.T) {
	value, changed, err := validate(Integer{Minimum: null.NewInt64(100)}, 5)
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, 100, value)
}

func TestValidateInteger_Direct_AllWidths(t *testing.T) {
	// validate_Integer accepts every signed and unsigned integer width
	cases := []any{
		int(1), int8(1), int16(1), int32(1), int64(1),
		uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
	}

	for _, value := range cases {
		_, _, err := validate_Integer(Integer{}, value)
		require.NoError(t, err)
	}
}

func TestValidateInteger_Direct_NotInteger(t *testing.T) {
	// A non-integer value reaches the fall-through error branch
	_, _, err := validate_Integer(Integer{}, "not-an-int")
	require.Error(t, err)
}

func TestValidate_Integer_Required(t *testing.T) {
	_, _, err := validate(Integer{Required: true}, 0)
	require.Error(t, err)
}

func TestValidate_Integer_MultipleOf(t *testing.T) {
	_, _, err := validate(Integer{MultipleOf: null.NewInt64(5)}, 23)
	require.Error(t, err)
}

func TestValidate_Integer_Enum(t *testing.T) {
	_, _, err := validate(Integer{Enum: []int{1, 2, 3}}, 9)
	require.Error(t, err)
}

func TestValidate_Number(t *testing.T) {
	value, changed, err := validate(Number{Maximum: null.NewFloat(10)}, 99.0)
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, 10.0, value)
}

func TestValidate_Number_CoercesInput(t *testing.T) {
	_, _, err := validate(Number{}, "not-a-float")
	require.Error(t, err)
}

func TestValidate_Number_Required(t *testing.T) {
	_, _, err := validate(Number{Required: true}, 0.0)
	require.Error(t, err)
}

func TestValidate_Number_MultipleOf(t *testing.T) {
	_, _, err := validate(Number{MultipleOf: null.NewFloat(2)}, 7.0)
	require.Error(t, err)
}

func TestValidate_Number_Enum(t *testing.T) {
	_, _, err := validate(Number{Enum: []float64{1, 2}}, 9.0)
	require.Error(t, err)
}

func TestValidate_String_CoercesInput(t *testing.T) {
	value, changed, err := validate(String{}, 42)
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, "42", value)
}

func TestValidate_String_Truncates(t *testing.T) {
	value, changed, err := validate(String{MaxLength: 5}, "this-is-too-long")
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, "this-", value)
}

func TestValidate_String_FormatError(t *testing.T) {
	_, _, err := validate(String{Format: "email"}, "not-an-email")
	require.Error(t, err)
}

func TestValidate_InvalidElementType(t *testing.T) {
	// An unknown element type produces an internal error
	_, _, err := validate(unknownElement{}, "value")
	require.Error(t, err)
}

func TestValidate_Array_NotArrayGetterSetter(t *testing.T) {
	_, _, err := validate(Array{Items: String{}}, "not-an-array")
	require.Error(t, err)
}

func TestValidate_Array_LengthRules(t *testing.T) {
	schema := New(Array{Items: String{}, MinLength: 5})

	_, _, err := Validate(schema, &testArrayA{"one"})
	require.Error(t, err)
}
