package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// testFloatGetter tests getting/setting float values
type testFloatGetter struct {
	value float64
}

func (t testFloatGetter) GetFloatOK(_ string) (float64, bool) {
	return t.value, true
}

func (t *testFloatGetter) SetFloat(_ string, value float64) bool {
	t.value = value
	return true
}

func TestFloatGetter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Number{},
		},
	})

	var getter testFloatGetter

	value, err := schema.Get(&getter, "value")

	require.Nil(t, err)
	require.Equal(t, float64(0), value)

	require.Nil(t, schema.Set(&getter, "value", 12345678))
	require.Equal(t, float64(12345678), getter.value)
}

// testFloatPointer tests getting/setting bool values via a pointer
type testFloatPointer struct {
	value float64
}

// GetPointer gets a pointer to the float value
func (test *testFloatPointer) GetPointer(_ string) (any, bool) {
	return &test.value, true
}

// TestFloatPointer tests getting/setting float values via a pointer
func TestFloatPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Number{BitSize: 64, Minimum: null.NewFloat(100), Maximum: null.NewFloat(1000000)}},
	},
	)

	var getter testFloatPointer

	// Validate correct value
	require.Nil(t, schema.Set(&getter, "value", 123456))
	require.Equal(t, float64(123456), getter.value)
	require.Nil(t, schema.Validate(&getter))

	// Validate incorrect value (to small)
	require.Nil(t, schema.Set(&getter, "value", float64(1)))
	require.Equal(t, float64(1), getter.value)
	require.Error(t, schema.Validate(&getter))

	// Validate incorrect value (too large)
	require.Nil(t, schema.Set(&getter, "value", float64(1000000000)))
	require.Equal(t, float64(1000000000), getter.value)
	require.Error(t, schema.Validate(&getter))
}
