package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// testStringGetter tests getting/setting string values
type testStringGetter struct {
	value string
}

func (t testStringGetter) GetStringOK(name string) (string, bool) {
	return t.value, true
}

func (t *testStringGetter) SetString(name string, value string) bool {
	t.value = value
	return true
}

func TestStringGetter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Required: true},
		},
	})

	var getter testStringGetter

	value, err := schema.Get(&getter, "value")

	require.Nil(t, err)
	require.Equal(t, "", value)
	require.Error(t, schema.Validate(&getter)) // Invalid because field is required

	require.Nil(t, schema.Set(&getter, "value", "hello"))
	require.Equal(t, "hello", getter.value)
	require.Nil(t, schema.Validate(&getter))
}

func TestStringValidator(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{MinLength: 10, MaxLength: 20},
		},
	})

	var getter testStringGetter

	// Test Min-Length
	require.Nil(t, schema.Set(&getter, "value", "too-short"))
	require.Equal(t, "too-short", getter.value)

	require.Error(t, schema.Validate(&getter))

	// Test Max-Length
	require.Nil(t, schema.Set(&getter, "value", "this-one-is-way-too-long-so-it-shouldn't-validate"))
	require.Error(t, schema.Validate(&getter))

	// Test OK
	require.Nil(t, schema.Set(&getter, "value", "this-is-just-right"))
	require.Nil(t, schema.Validate(&getter))

}

// testStringPointer tests getting/setting bool values via a pointer
type testStringPointer struct {
	value string
}

func (test *testStringPointer) GetPointer(name string) (any, bool) {
	return &test.value, true
}

func TestStringPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Enum: []string{"hello", "there", "general", "kenobi"}},
		},
	})

	var getter testStringPointer

	// Validate correct value
	require.Nil(t, schema.Set(&getter, "value", "kenobi"))
	require.Equal(t, "kenobi", getter.value)
	require.Nil(t, schema.Validate(&getter))

	// Validate incorrect value
	require.Nil(t, schema.Set(&getter, "value", "invalid"))
	require.Equal(t, "invalid", getter.value)
	require.Error(t, schema.Validate(&getter))
}
