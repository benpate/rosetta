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
func (unknownElement) validateRequiredIf(Schema, list.List, any) error { return nil }
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

func TestValidate_Boolean_UsesExistingInputFormat(t *testing.T) {
	// A coercible value is accepted and converted to a boolean
	newValue, changed, err := validate(Boolean{}, "true")
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, "true", newValue)
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
	// With no BitSize, the validated result is a plain int.
	value, changed, err := validate(Integer{Minimum: null.NewInt64(100)}, 5)
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, 100, value)
}

func TestValidate_Integer_BitSizeTypes(t *testing.T) {
	// validate types the result to match the element's BitSize.
	cases := []struct {
		bitSize  int
		expected any
	}{
		{0, int(1)},
		{8, int8(1)},
		{16, int16(1)},
		{32, int32(1)},
		{64, int64(1)},
	}

	for _, test := range cases {
		value, _, err := validate(Integer{BitSize: test.bitSize}, 1)
		require.NoError(t, err, test.bitSize)
		require.Equal(t, test.expected, value, test.bitSize)
	}
}

func TestValidate_Integer_NoTruncation(t *testing.T) {
	// A value larger than MaxInt32 must pass through validation without being
	// truncated. BitSize:64 pins the result to int64 on every platform
	// (regression for coercing through a 32-bit `int`).
	const big = int64(5_000_000_000) // > math.MaxInt32

	value, changed, err := validate(Integer{BitSize: 64, Maximum: null.NewInt64(10_000_000_000)}, big)
	require.NoError(t, err)
	require.False(t, changed)
	require.Equal(t, big, value)

	// Clamping to a 64-bit maximum must also preserve full width.
	value, changed, err = validate(Integer{BitSize: 64, Maximum: null.NewInt64(big)}, int64(9_000_000_000))
	require.NoError(t, err)
	require.True(t, changed)
	require.Equal(t, big, value)
}

func TestValidateInteger_Direct_AllWidths(t *testing.T) {
	// validate_Integer accepts every signed integer width
	cases := []any{
		int(1), int8(1), int16(1), int32(1), int64(1),
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

func TestValidate_Integer_BitSizeOverflow(t *testing.T) {

	// A value outside the declared bit size must be rejected, not silently wrapped.
	// Before the fix, int8(300) wrapped to 44 and slipped past validation.
	cases := []struct {
		bitSize int
		value   int
	}{
		{8, 300},            // > MaxInt8 (127)
		{8, -200},           // < MinInt8 (-128)
		{16, 70000},         // > MaxInt16 (32767)
		{32, 5_000_000_000}, // > MaxInt32
	}

	for _, c := range cases {
		_, _, err := validate(Integer{BitSize: c.bitSize}, c.value)
		require.Error(t, err, "bitSize=%d value=%d should be rejected", c.bitSize, c.value)
	}
}

func TestValidate_Integer_BitSizeInRange(t *testing.T) {

	// Values that fit the declared width pass through with the correct typed result.
	value, _, err := validate(Integer{BitSize: 8}, 100)
	require.NoError(t, err)
	require.Equal(t, int8(100), value)

	value, _, err = validate(Integer{BitSize: 16}, 30000)
	require.NoError(t, err)
	require.Equal(t, int16(30000), value)

	value, _, err = validate(Integer{BitSize: 32}, 2_000_000_000)
	require.NoError(t, err)
	require.Equal(t, int32(2_000_000_000), value)
}
