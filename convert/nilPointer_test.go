package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// nilPointerInputs returns every flavor of nil that reaches this package's reflection paths.
// reflect.Value.Elem() hands back the ZERO Value for a nil pointer, and calling Interface() on
// that panics -- so any function that dereferences without checking IsNil first will crash on
// one of these.
func nilPointerInputs() map[string]any {

	var nilSlicePtr *[]string
	var nilMapPtr *map[string]any
	var nilIntPtr *int
	var nilStringPtr *string
	nilPointerToPointer := &nilSlicePtr

	return map[string]any{
		"untyped nil":        nil,
		"nil *[]string":      nilSlicePtr,
		"nil *map":           nilMapPtr,
		"nil *int":           nilIntPtr,
		"nil *string":        nilStringPtr,
		"pointer to nil ptr": nilPointerToPointer,
		"nil map":            map[string]any(nil),
		"nil slice":          []string(nil),
		"nil chan":           (chan int)(nil),
		"nil func":           (func())(nil),
	}
}

// TestNilPointer_NoPanics runs every exported conversion that uses reflection against every
// flavor of nil.  Each one must return a zero value rather than panicking.
func TestNilPointer_NoPanics(t *testing.T) {

	conversions := map[string]func(any){
		"Bool":               func(v any) { Bool(v) },
		"Bytes":              func(v any) { Bytes(v) },
		"Element":            func(v any) { Element(v) },
		"Float":              func(v any) { Float(v) },
		"Int":                func(v any) { Int(v) },
		"Int32":              func(v any) { Int32(v) },
		"Int64":              func(v any) { Int64(v) },
		"IsMap":              func(v any) { IsMap(v) },
		"IsSlice":            func(v any) { IsSlice(v) },
		"JoinString":         func(v any) { JoinString(v, ",") },
		"MapOfAny":           func(v any) { MapOfAny(v) },
		"MapOfInt":           func(v any) { MapOfInt(v) },
		"MapOfInt32":         func(v any) { MapOfInt32(v) },
		"MapOfSliceOfString": func(v any) { MapOfSliceOfString(v) },
		"MapOfString":        func(v any) { MapOfString(v) },
		"SliceLength":        func(v any) { SliceLength(v) },
		"SliceOf[any]":       func(v any) { SliceOf[any](v) },
		"SliceOf[string]":    func(v any) { SliceOf[string](v) },
		"SliceOfAny":         func(v any) { SliceOfAny(v) },
		"SliceOfFloat":       func(v any) { SliceOfFloat(v) },
		"SliceOfInt":         func(v any) { SliceOfInt(v) },
		"SliceOfInt64":       func(v any) { SliceOfInt64(v) },
		"SliceOfMap":         func(v any) { SliceOfMap(v) },
		"SliceOfString":      func(v any) { SliceOfString(v) },
		"String":             func(v any) { String(v) },
		"Time":               func(v any) { Time(v) },
	}

	for conversionName, conversion := range conversions {
		for inputName, input := range nilPointerInputs() {
			require.NotPanics(t, func() { conversion(input) }, "%s(%s)", conversionName, inputName)
		}
	}
}

// TestNilPointer_Results pins what a nil pointer converts INTO, so that "does not panic" cannot
// quietly become "returns something surprising".
func TestNilPointer_Results(t *testing.T) {

	var nilSlicePtr *[]string

	{ // A nil pointer holds no collection, so conversion reports failure
		result, ok := SliceOfStringOk(nilSlicePtr)
		require.False(t, ok)
		require.Empty(t, result)

		_, ok = SliceOfAnyOk(nilSlicePtr)
		require.False(t, ok)

		_, ok = SliceOfIntOk(nilSlicePtr)
		require.False(t, ok)
	}

	{ // ...and it is neither a slice nor a map, and has no length
		require.False(t, IsSlice(nilSlicePtr))
		require.False(t, IsMap(nilSlicePtr))
		require.Zero(t, SliceLength(nilSlicePtr))
	}

	// Element dereferences to nothing rather than to a zero Value
	require.Nil(t, Element(nilSlicePtr))
}

// TestIsMap_DereferencesToIsMap covers a pointer to a map.  The pointer branch of IsMap used to
// hand off to IsSlice, so a *map reported FALSE.
func TestIsMap_DereferencesToIsMap(t *testing.T) {

	value := map[string]any{"a": 1}

	require.True(t, IsMap(value))
	require.True(t, IsMap(&value))

	// A pointer to a slice is still not a map
	slice := []string{"a"}
	require.False(t, IsMap(&slice))
	require.True(t, IsSlice(&slice))
}
