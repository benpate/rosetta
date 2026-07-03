package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// isNil ----------------------------------------------------------------------

func TestIsNil(t *testing.T) {
	var nilPtr *String
	var nilMap map[string]int
	var nilSlice []int

	// Nillable kinds report their nil-ness correctly.
	require.True(t, isNil(nil))
	require.True(t, isNil(nilPtr))
	require.True(t, isNil(nilMap))
	require.True(t, isNil(nilSlice))

	// Non-nil values of every relevant kind report false...
	require.False(t, isNil(&String{}))
	require.False(t, isNil(map[string]int{}))
	require.False(t, isNil([]int{}))
	require.False(t, isNil("a string"))
	require.False(t, isNil(42))

	// ...including non-nillable kinds (Array, Struct), which previously
	// panicked when passed to reflect.Value.IsNil().
	require.NotPanics(t, func() { require.False(t, isNil([3]int{})) })
	require.NotPanics(t, func() { require.False(t, isNil(String{})) })
}

// validate_String branches via the dispatcher --------------------------------

func TestValidateString_Required(t *testing.T) {
	_, _, err := validate(String{Required: true}, "")
	require.Error(t, err)
}

func TestValidateString_MinValue(t *testing.T) {
	_, _, err := validate(String{MinValue: "m"}, "a")
	require.Error(t, err)
}

func TestValidateString_MaxValue(t *testing.T) {
	_, _, err := validate(String{MaxValue: "m"}, "z")
	require.Error(t, err)
}

func TestValidateString_MinLength(t *testing.T) {
	_, _, err := validate(String{MinLength: 10}, "short")
	require.Error(t, err)
}

func TestValidateString_Enum(t *testing.T) {
	_, _, err := validate(String{Enum: []string{"a", "b"}}, "z")
	require.Error(t, err)
}

func TestValidateString_FormatTransforms(t *testing.T) {
	// The default no-html format strips HTML tags and flags the value as changed
	value, changed, err := validate(String{}, "<b>hi</b>")
	require.NoError(t, err)
	require.NotEmpty(t, changed)
	require.Equal(t, "hi", value)
}

// Typed Get via the PointerGetter interface ----------------------------------

func TestGet_Integer32_ViaPointer(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{}}})

	result, err := schema.Get(&testIntPointer{value: 7}, "value")
	require.NoError(t, err)
	require.Equal(t, 7, result)
}

func TestGet_Integer64_ViaPointer(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{BitSize: 64}}})

	result, err := schema.Get(&testInt64Pointer{value: 7}, "value")
	require.NoError(t, err)
	require.Equal(t, int64(7), result)
}

func TestGet_Number_ViaPointer(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Number{}}})

	result, err := schema.Get(&testFloatPointer{value: 1.5}, "value")
	require.NoError(t, err)
	require.Equal(t, 1.5, result)
}

func TestGet_String_ViaPointer(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": String{}}})

	result, err := schema.Get(&testStringPointer{value: "hello"}, "value")
	require.NoError(t, err)
	require.Equal(t, "hello", result)
}

// ValueGetter / ValueSetter (whole-value access) -----------------------------

// wholeValueObject reads and writes its entire value at the empty path.
type wholeValueObject struct {
	value string
}

func (object wholeValueObject) GetValue() any {
	return object.value
}

func (object *wholeValueObject) SetValue(value any) error {
	object.value = value.(string)
	return nil
}

func TestGet_WholeValue(t *testing.T) {
	// An empty path with a ValueGetter returns the entire value
	result, err := New(String{}).Get(wholeValueObject{value: "hello"}, "")
	require.NoError(t, err)
	require.Equal(t, "hello", result)
}

func TestGet_WholeValue_NoValueGetter(t *testing.T) {
	// An empty path without a ValueGetter returns the object itself
	result, err := New(String{}).Get("plain", "")
	require.NoError(t, err)
	require.Equal(t, "plain", result)
}

func TestSet_WholeValue(t *testing.T) {
	// An empty path with a ValueSetter writes the entire value
	object := &wholeValueObject{}
	err := New(String{}).Set(object, "", "hello")

	require.NoError(t, err)
	require.Equal(t, "hello", object.value)
}

func TestSet_WholeValue_NoValueSetter(t *testing.T) {
	// An empty path without a ValueSetter cannot be set
	require.Error(t, New(String{}).Set("plain", "", "hello"))
}
