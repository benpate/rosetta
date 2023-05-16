package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// testBoolGetter tests getting/setting bool values via GetBoolOK/SetBool
type testBoolGetter struct {
	value bool
}

func (t testBoolGetter) GetBoolOK(name string) (bool, bool) {
	return t.value, true
}

func (t *testBoolGetter) SetBool(name string, value bool) bool {
	t.value = value
	return true
}

func TestBoolGetter(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Boolean{},
		},
	})

	var getter testBoolGetter

	value, err := schema.Get(&getter, "value")

	require.Nil(t, err)
	require.Equal(t, false, value)

	require.Nil(t, schema.Set(&getter, "value", true))
	require.Equal(t, true, getter.value)

	require.Nil(t, schema.Validate(&getter))
}

// testBoolPointer tests getting/setting bool values via a pointer
type testBoolPointer struct {
	value bool
}

func (test *testBoolPointer) GetPointer(name string) (any, bool) {
	return &test.value, true
}

func TestBoolPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": Boolean{},
		},
	})

	var getter testBoolPointer

	value, err := schema.Get(&getter, "value")

	require.Nil(t, err)
	require.Equal(t, false, value)

	require.Nil(t, schema.Set(&getter, "value", true))
	require.Equal(t, true, getter.value)

	require.Nil(t, schema.Validate(&getter))
}
