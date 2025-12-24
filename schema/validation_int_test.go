package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// testIntGetter tests getting/setting int values
type testIntGetter struct {
	value int
}

// GetIntOK gets an int value
func (t testIntGetter) GetIntOK(_ string) (int, bool) {
	return t.value, true
}

// SetInt sets an int value
func (t *testIntGetter) SetInt(_ string, value int) bool {
	t.value = value
	return true
}

// TestIntGetter tests getting/setting bool values
func TestIntGetter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Integer{},
		},
	})

	var getter testIntGetter

	value, err := schema.Get(&getter, "value")

	require.Nil(t, err)
	require.Equal(t, 0, value)

	require.Nil(t, schema.Set(&getter, "value", 12345678))
	require.Equal(t, 12345678, getter.value)
}

// testIntPointer tests getting/setting bool values via a pointer
type testIntPointer struct {
	value int
}

func (test *testIntPointer) GetPointer(_ string) (any, bool) {
	return &test.value, true
}

func TestIntPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Integer{Minimum: null.NewInt64(100), Maximum: null.NewInt64(1000000)}},
	},
	)

	var getter testIntPointer

	// Validate correct value
	require.Nil(t, schema.Set(&getter, "value", 123456))
	require.Equal(t, int(123456), getter.value)
	require.Nil(t, schema.Validate(&getter))

	// Validate incorrect value (to small)
	require.Nil(t, schema.Set(&getter, "value", int(1)))
	require.Equal(t, int(1), getter.value)
	require.Error(t, schema.Validate(&getter))

	// Validate incorrect value (too large)
	require.Nil(t, schema.Set(&getter, "value", int(1000000000)))
	require.Equal(t, int(1000000000), getter.value)
	require.Error(t, schema.Validate(&getter))
}
