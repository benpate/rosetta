package schema

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

// TestSet_PathDepthIsBounded confirms that Set refuses a path with more than maxPathSegments
// segments as a validation error, and does not repeat the path in that error
func TestSet_PathDepthIsBounded(t *testing.T) {

	s := New(Object{Properties: ElementMap{"data": Object{Wildcard: Any{}}}})

	// "data", the segments in between, then "x"
	deepest := "data." + strings.Repeat("a.", maxPathSegments-2) + "x"
	require.NoError(t, s.Set(&mapof.Any{}, deepest, "v"))

	tooDeep := "data." + strings.Repeat("a.", maxPathSegments-1) + "x"
	require.True(t, derp.IsValidationError(s.Set(&mapof.Any{}, tooDeep, "v")))

	// Before 2026-10-22, about 1 MB of such a path overflowed the stack, which no recover catches
	long := "data." + strings.Repeat("a.", 50_000) + "x"
	err := s.Set(&mapof.Any{}, long, "v")
	require.True(t, derp.IsValidationError(err))

	encoded, _ := json.Marshal(err)
	require.Less(t, len(encoded), 1000, "the error does not repeat the path")
}

// TestSet_ValidationErrorNamesField confirms that a validation failure in Set is
// re-rooted with the failing path prefixed, so handlers that surface the ROOT
// message of a 422 chain can tell the user WHICH field failed.
func TestSet_ValidationErrorNamesField(t *testing.T) {

	schema := New(Object{Properties: ElementMap{
		"masterKey": String{MinLength: 64, MaxLength: 64},
	}})

	object := mapStringObject{}

	err := schema.Set(&object, "masterKey", "")
	require.Error(t, err)
	require.True(t, derp.IsValidationError(err))
	require.Equal(t, "masterKey: Minimum length is 64", derp.RootMessage(err))
}

// TestSet_ValidationErrorSurvivesWrapping confirms that the field-named root message
// still surfaces after the error is wrapped by pipeline layers (the shape produced by
// SetAll and its callers).
func TestSet_ValidationErrorSurvivesWrapping(t *testing.T) {

	schema := New(Object{Properties: ElementMap{
		"masterKey": String{MinLength: 64, MaxLength: 64},
	}})

	object := mapStringObject{}

	err := schema.SetAll(&object, map[string]any{"masterKey": ""})
	require.Error(t, err)

	// Mimic a handler-level wrap on top of SetAll's own wrap
	wrapped := derp.Wrap(err, "handler.Location", "Setting config values")

	require.True(t, derp.IsValidationError(wrapped))
	require.Equal(t, "masterKey: Minimum length is 64", derp.RootMessage(wrapped))
}

// TestSet_UnknownPathIsValidation confirms that setting an unknown path reports a
// validation error that names the field, rather than a generic 400.
func TestSet_UnknownPathIsValidation(t *testing.T) {

	schema := New(Object{Properties: ElementMap{
		"name": String{},
	}})

	object := mapStringObject{}

	err := schema.Set(&object, "missing", "x")
	require.Error(t, err)
	require.True(t, derp.IsValidationError(err))
	require.Equal(t, "Unknown field: missing", derp.RootMessage(err))
}

// TestSet_ValidValueStillSets confirms the happy path is unchanged by the error
// re-rooting: a valid value lands in the object with no error.
func TestSet_ValidValueStillSets(t *testing.T) {

	schema := New(Object{Properties: ElementMap{
		"name": String{MaxLength: 16},
	}})

	object := mapStringObject{}

	require.NoError(t, schema.Set(&object, "name", "Sarah Connor"))
	require.Equal(t, "Sarah Connor", object["name"])
}

// mapStringObject is a minimal StringSetter/PointerGetter target for Set tests.
type mapStringObject map[string]string

func (object *mapStringObject) GetPointer(name string) (any, bool) {
	value := (*object)[name]
	return &value, true
}

func (object *mapStringObject) SetString(name string, value string) bool {
	if *object == nil {
		*object = mapStringObject{}
	}
	(*object)[name] = value
	return true
}
