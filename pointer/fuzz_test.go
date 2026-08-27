package pointer

import (
	"reflect"
	"strings"
	"testing"
)

// carriesOwnReference returns TRUE for the kinds that already hold a reference, and which To()
// therefore hands straight back instead of wrapping in a new pointer.
func carriesOwnReference(value any) bool {

	switch reflect.ValueOf(value).Kind() {

	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return true

	default:
		return false
	}
}

// fuzzValues returns the awkward values worth handing to To(), built from fuzzed text.
func fuzzValues(text string) []any {

	var nilPointer *int
	var nilMap map[string]any
	var nilSlice []string
	var nilInterface any

	return []any{
		nil,
		nilInterface,
		nilPointer,
		nilMap,
		nilSlice,
		text,
		len(text),
		[]string{text},
		map[string]any{text: text},
		struct{ Name string }{Name: text},
		&struct{ Name string }{Name: text},
		make(chan int),
		func() {},
	}
}

// FuzzTo_NeverPanics feeds every awkward value to the reflection-based pointer helper.
// To() builds a pointer with reflect.New, which panics on a value carrying no type, so an
// untyped nil has to be recognized before it ever reaches there.
func FuzzTo_NeverPanics(f *testing.F) {

	f.Add("")
	f.Add("value")
	f.Add("\x00")
	f.Add(strings.Repeat("a", 1024))

	f.Fuzz(func(t *testing.T, text string) {

		for _, value := range fuzzValues(text) {

			result := To(value)

			// RULE: A nil has no type to point at, and comes straight back
			if value == nil {
				if result != nil {
					t.Fatalf("To(nil) returned %#v", result)
				}
				continue
			}

			if carriesOwnReference(value) {
				continue
			}

			// Everything else comes back as a pointer that dereferences to an equal value
			if reflect.ValueOf(result).Kind() != reflect.Pointer {
				t.Fatalf("To(%#v) returned a %s, not a pointer",
					value, reflect.ValueOf(result).Kind())
			}

			if got := reflect.ValueOf(result).Elem().Interface(); !reflect.DeepEqual(got, value) {
				t.Fatalf("To(%#v) dereferenced to %#v", value, got)
			}
		}
	})
}
