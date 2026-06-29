package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// pointerGetterStruct models the real-world case (e.g. Emissary's model.PersonLink) where
// the scalar accessors use POINTER receivers, so a by-value copy of the struct satisfies
// neither StringGetter nor PointerGetter. An array of these is the case that v0.28's
// validate_Array originally could not handle.
type pointerGetterStruct struct {
	Name string
	URL  string
}

func pointerGetterStruct_Schema() Element {
	return Object{
		Properties: map[string]Element{
			"name": String{MaxLength: 8},
			"url":  String{Format: "url"},
		},
	}
}

// RULE: pointer receivers only — a VALUE of this type implements none of the getter/setter interfaces.
func (s *pointerGetterStruct) GetStringOK(name string) (string, bool) {
	switch name {
	case "name":
		return s.Name, true
	case "url":
		return s.URL, true
	}
	return "", false
}

func (s *pointerGetterStruct) GetPointer(name string) (any, bool) {
	switch name {
	case "name":
		return &s.Name, true
	case "url":
		return &s.URL, true
	}
	return nil, false
}

func (s *pointerGetterStruct) SetString(name string, value string) bool {
	switch name {
	case "name":
		s.Name = value
		return true
	case "url":
		s.URL = value
		return true
	}
	return false
}

// pointerGetterSlice is the array wrapper; like sliceof.Object[T], GetIndex returns T BY VALUE.
type pointerGetterSlice []pointerGetterStruct

func (slice pointerGetterSlice) Length() int { return len(slice) }

func (slice pointerGetterSlice) GetIndex(index int) (any, bool) {
	if (index < 0) || (index >= len(slice)) {
		return nil, false
	}
	return slice[index], true // by value, on purpose
}

func (slice *pointerGetterSlice) SetIndex(index int, value any) bool {
	if typed, ok := value.(pointerGetterStruct); ok {
		for index >= len(*slice) {
			*slice = append(*slice, pointerGetterStruct{})
		}
		(*slice)[index] = typed
		return true
	}
	return false
}

// TestValidateArray_ObjectItems_PointerReceivers is the regression test for the v0.28 array-validation
// bug: validating an array whose items are structs with pointer-receiver accessors must succeed by
// validating against an addressable copy. Before the fix this returned "Object must be a StringGetter
// or PointerGetter".
func TestValidateArray_ObjectItems_PointerReceivers(t *testing.T) {

	s := New(Array{Items: pointerGetterStruct_Schema()})

	// A clean value should validate without error and without being rewritten.
	value := pointerGetterSlice{
		{Name: "Alice", URL: "https://example.com/alice"},
		{Name: "Bob", URL: "https://example.com/bob"},
	}

	changed, err := s.Validate(&value)
	require.Nil(t, err)
	require.False(t, changed)
}

// TestValidateArray_ObjectItems_WritesBackRewrite confirms that when nested validation rewrites a
// field inside an array element, the rewrite is persisted back into the slice. Per the top-level
// Validate contract, a value that had to be rewritten is reported as invalid (a non-nil error), but
// the write-back must still reach the original slice element through SetIndex.
func TestValidateArray_ObjectItems_WritesBackRewrite(t *testing.T) {

	s := New(Array{Items: pointerGetterStruct_Schema()})

	// "TooLongName" exceeds MaxLength:8 and must be truncated by validation.
	value := pointerGetterSlice{
		{Name: "TooLongName", URL: "https://example.com/x"},
	}

	// A rewrite is required, so Validate reports the value as not-valid-as-given.
	_, err := s.Validate(&value)
	require.NotNil(t, err)

	// ...but the truncation must have been written back into the slice element.
	require.Equal(t, "TooLongN", value[0].Name) // truncated to 8 runes and written back
}
