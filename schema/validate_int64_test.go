package schema

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// testInt64Getter tests getting/setting int64 values
type testInt64Getter struct {
	value int64
}

func (t testInt64Getter) GetInt64OK(_ string) (int64, bool) {
	return t.value, true
}

func (t *testInt64Getter) SetInt64(_ string, value int64) bool {
	t.value = value
	return true
}

func TestInt64Getter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Integer{BitSize: 64},
		},
	})

	var getter testInt64Getter

	value, err := schema.Get(&getter, "value")

	require.NoError(t, err)
	require.Equal(t, int64(0), value)

	require.NoError(t, schema.Set(&getter, "value", 12345678))
	require.Equal(t, int64(12345678), getter.value)
}

// testInt64Pointer tests getting/setting bool values via a pointer
type testInt64Pointer struct {
	value int64
}

func (test *testInt64Pointer) GetPointer(_ string) (any, bool) {
	return &test.value, true
}

func TestInt64Pointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Integer{BitSize: 64, Minimum: null.NewInt64(100), Maximum: null.NewInt64(1000000)}},
	},
	)

	var getter testInt64Pointer

	// Validate correct value
	require.NoError(t, schema.Set(&getter, "value", 123456))
	require.Equal(t, int64(123456), getter.value)
	_, _, err := Validate(schema, &getter)
	require.NoError(t, err)

	// Value below the minimum is clamped to the minimum (does not fail)
	require.NoError(t, schema.Set(&getter, "value", int64(1)))
	require.Equal(t, int64(100), getter.value)
	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)

	// Value above the maximum is clamped to the maximum (does not fail)
	require.NoError(t, schema.Set(&getter, "value", int64(1000000000)))
	require.Equal(t, int64(1000000), getter.value)
	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)
}
