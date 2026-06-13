package schema

import (
	"net/url"
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestSchema_SetAll(t *testing.T) {
	schema := New(testStructA_Schema())
	value := newTestStructA()

	err := schema.SetAll(&value, map[string]any{
		"name":     "Sarah Connor",
		"latitude": 12.5,
	})

	require.NoError(t, err)
	require.Equal(t, "Sarah Connor", value.Name)
	require.Equal(t, 12.5, value.Latitude)
}

func TestSchema_SetAll_Error(t *testing.T) {
	// Setting an unknown path returns an error
	schema := New(testStructA_Schema())
	value := newTestStructA()

	require.Error(t, schema.SetAll(&value, map[string]any{"missing": "x"}))
}

func TestSchema_SetURLValues(t *testing.T) {
	schema := New(testStructA_Schema())
	value := newTestStructA()

	err := schema.SetURLValues(&value, url.Values{
		"name": []string{"Kyle Reese"},
	})

	require.NoError(t, err)
	require.Equal(t, "Kyle Reese", value.Name)
}

// TestSchema_SetURLValues_CoercesAndWritesBack confirms that values which the
// schema coerces during validation (truncation, clamping) are written back to
// the object, not just the raw input.
func TestSchema_SetURLValues_CoercesAndWritesBack(t *testing.T) {
	schema := New(Object{
		Properties: map[string]Element{
			"name":     String{MaxLength: 5},
			"latitude": Number{BitSize: 64, Maximum: null.NewFloat(10)},
		},
	})
	value := newTestStructA()

	err := schema.SetURLValues(&value, url.Values{
		"name":     []string{"abcdefghij"},
		"latitude": []string{"99.5"},
	})

	require.NoError(t, err)
	require.Equal(t, "abcde", value.Name)  // truncated to MaxLength
	require.Equal(t, 10.0, value.Latitude) // clamped to Maximum
}

// boolPointerObject exposes a single boolean property via the PointerGetter
// interface only (it is not a BoolGetter/BoolSetter).
type boolPointerObject struct {
	value bool
}

func (object *boolPointerObject) GetPointer(_ string) (any, bool) {
	return &object.value, true
}

func TestGetSet_Boolean_ViaPointer(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Boolean{}}})
	object := &boolPointerObject{}

	require.NoError(t, schema.Set(object, "value", true))
	require.True(t, object.value)

	result, err := schema.Get(object, "value")
	require.NoError(t, err)
	require.Equal(t, true, result)
}

func TestGet_Boolean_Unsupported(t *testing.T) {
	// An object that supports neither BoolGetter nor PointerGetter fails
	schema := New(Object{Properties: ElementMap{"value": Boolean{}}})

	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)
}

func TestGet_String_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": String{}}})

	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)
}

func TestGet_Number_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Number{}}})

	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)
}

func TestGet_Integer64_Unsupported(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": Integer{BitSize: 64}}})

	_, err := schema.Get(unsupportedObject{}, "value")
	require.Error(t, err)
}

func TestGet_InvalidProperty(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": String{}}})

	_, err := schema.Get(unsupportedObject{}, "missing")
	require.Error(t, err)
}

func TestSet_InvalidPath(t *testing.T) {
	schema := New(Object{Properties: ElementMap{"value": String{}}})

	require.Error(t, schema.Set(&boolPointerObject{}, "missing", "x"))
}

func TestSet_Unsupported(t *testing.T) {
	// The target supports no setter interface for the property
	schema := New(Object{Properties: ElementMap{"value": String{}}})

	require.Error(t, schema.Set(unsupportedObject{}, "value", "x"))
}

// unsupportedObject implements no getter or setter interfaces.
type unsupportedObject struct{}
