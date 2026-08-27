package mapof

import (
	"strings"
	"testing"

	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/schema"
)

// fuzzSchema returns a small nested schema used to drive SetObject down a real path.
//
// Every leaf is a String on purpose. SetObject is the low-level writer -- it stores whatever
// value it is handed without consulting the element type, because schema.Schema.Set validates
// and converts upstream. A non-string leaf here would only test that mismatch, not the path
// traversal this target exists to exercise.
func fuzzSchema() schema.Element {

	return schema.Object{
		Properties: schema.ElementMap{
			"name": schema.String{},
			"age":  schema.String{},
			"child": schema.Object{
				Properties: schema.ElementMap{
					"name": schema.String{},
					"grandchild": schema.Object{
						Properties: schema.ElementMap{
							"name": schema.String{},
						},
					},
				},
			},
		},
	}
}

// FuzzAny_SetObject drives arbitrary dotted paths through the recursive object setter.
// SetObject walks a caller-supplied path, creating child maps as it descends, so a malformed
// or hostile path is exactly the input it must survive: schema.Set routes form and JSON writes
// through here, and the path segments come from whatever the client posted.
func FuzzAny_SetObject(f *testing.F) {

	f.Add("name")
	f.Add("child.name")
	f.Add("child.grandchild.name")
	f.Add("")
	f.Add(".")
	f.Add("..")
	f.Add("...")
	f.Add(".name")
	f.Add("name.")
	f.Add("child..name")
	f.Add("nonexistent.path")
	f.Add("child.nonexistent")
	f.Add(strings.Repeat("child.", 64) + "name")
	f.Add("\x00")
	f.Add("child.\x00.name")

	element := fuzzSchema()

	f.Fuzz(func(t *testing.T, path string) {

		object := NewAny()

		// A hostile path may legitimately be REFUSED, but must never panic or hang
		if err := object.SetObject(element, list.ByDot(path), "VALUE"); err != nil {
			return
		}

		// RULE: A write to a path the SCHEMA defines must be readable back through that same
		// path. Silently accepting a write that goes nowhere is how a form field vanishes.
		//
		// Paths the schema does NOT define are excluded: SetObject writes a single-segment key
		// straight into the map without consulting the schema, and schema.Schema.Set is what
		// rejects unknown fields before they ever reach here.
		if _, defined := schema.New(element).GetElement(path); !defined {
			return
		}

		result, err := schema.New(element).Get(&object, path)

		if err != nil {
			t.Fatalf("SetObject accepted path %q, but Get could not read it back: %v", path, err)
		}

		if result == nil {
			t.Fatalf("SetObject accepted path %q, but Get read back nil", path)
		}
	})
}

// FuzzAny_GetSetRoundTrip confirms that a value written through the typed setters reads back
// through the matching typed getter, for arbitrary keys.
func FuzzAny_GetSetRoundTrip(f *testing.F) {

	f.Add("name", "value")
	f.Add("", "")
	f.Add(".", "x")
	f.Add("a.b", "x")
	f.Add("\x00", "\x00")
	f.Add(strings.Repeat("k", 4096), strings.Repeat("v", 4096))
	f.Add("\xff\xfe", "\xff\xfe")

	f.Fuzz(func(t *testing.T, key string, value string) {

		object := NewAny()

		if !object.SetString(key, value) {
			return
		}

		// RULE: This map is deliberately SPARSE -- storing a zero value deletes the key
		// instead of writing it, so an empty string round-trips as "absent", not as "".
		if value == "" {
			if object.Length() != 0 {
				t.Fatalf("SetString(%q, \"\") stored a key instead of deleting it", key)
			}
			return
		}

		if got := object.GetString(key); got != value {
			t.Fatalf("SetString(%q, %q) then GetString read back %q", key, value, got)
		}

		// The key must also appear in Keys(), or callers that enumerate will miss it
		found := false

		for _, candidate := range object.Keys() {
			if candidate == key {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("SetString(%q) succeeded but the key is absent from Keys()", key)
		}

		// ...and Remove must actually remove it
		object.Remove(key)

		if object.Length() != 0 {
			t.Fatalf("Remove(%q) left %d keys behind", key, object.Length())
		}
	})
}

// FuzzAny_TypedGettersNeverPanic reads every key through every typed getter.  The getters
// coerce whatever happens to be stored, so a value of an unexpected type must produce a zero
// value rather than a panic.
func FuzzAny_TypedGettersNeverPanic(f *testing.F) {

	f.Add("key", "value")
	f.Add("key", "")
	f.Add("", "")
	f.Add("key", "12345")
	f.Add("key", "3.14159")
	f.Add("key", "true")
	f.Add("key", "2026-03-04T13:02:00Z")

	f.Fuzz(func(t *testing.T, key string, value string) {

		object := Any{key: value}

		_ = object.GetAny(key)
		_ = object.GetBool(key)
		_ = object.GetFloat(key)
		_ = object.GetInt(key)
		_ = object.GetInt64(key)
		_ = object.GetString(key)
		_ = object.GetTime(key)
		_ = object.GetSliceOfAny(key)
		_ = object.GetSliceOfString(key)
		_ = object.GetSliceOfInt(key)
		_ = object.GetSliceOfFloat(key)
		_ = object.GetMap(key)
		_ = object.IsZeroValue(key)

		// The same getters must survive a MISSING key just as well
		_ = object.GetTime("missing")
		_ = object.GetMap("missing")
		_ = object.GetSliceOfString("missing")
	})
}

// FuzzAny_SparseContract pins the sparse-map invariant across every typed setter: storing a
// ZERO value deletes the key rather than writing it, so "absent" and "zero" are the same
// state.  Anything that persists this map (JSON, BSON) depends on that, and a setter that
// quietly started writing zeroes would change every stored document.
func FuzzAny_SparseContract(f *testing.F) {

	f.Add("key")
	f.Add("")
	f.Add(".")
	f.Add("\x00")
	f.Add(strings.Repeat("k", 1024))

	f.Fuzz(func(t *testing.T, key string) {

		{ // Zero values are deleted, never stored
			object := NewAny()
			object.SetString(key, "")
			object.SetInt(key, 0)
			object.SetInt64(key, 0)
			object.SetFloat(key, 0)
			object.SetAny(key, nil)

			if object.Length() != 0 {
				t.Fatalf("a zero value was stored under %q: %v", key, object)
			}
		}

		{ // ...and a non-zero value under the same key IS stored
			object := NewAny()
			object.SetString(key, "value")

			if object.Length() != 1 {
				t.Fatalf("a non-zero value was NOT stored under %q", key)
			}

			// Writing a zero over it removes it again
			object.SetString(key, "")

			if object.Length() != 0 {
				t.Fatalf("writing a zero over %q did not remove it", key)
			}
		}

		// SetBool is the deliberate exception: FALSE is a real value, not an absence
		{
			object := NewAny()
			object.SetBool(key, false)

			if object.Length() != 1 {
				t.Fatalf("SetBool(%q, false) must store, not delete", key)
			}
		}
	})
}
