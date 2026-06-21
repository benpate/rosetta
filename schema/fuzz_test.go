package schema

import (
	"testing"
)

// FuzzUnmarshalJSON feeds arbitrary bytes to Schema.UnmarshalJSON to confirm that parsing untrusted
// JSON never panics, and that any schema that DOES parse can be re-marshaled without panicking.
func FuzzUnmarshalJSON(f *testing.F) {

	f.Add([]byte(`{"type":"string"}`))
	f.Add([]byte(`{"type":"integer","bitSize":8,"minimum":0,"maximum":100}`))
	f.Add([]byte(`{"type":"object","properties":{"name":{"type":"string"}}}`))
	f.Add([]byte(`{"type":"array","items":{"type":"integer"}}`))
	f.Add([]byte(`{"type":"object","properties":{"a":{"type":"array","items":{"type":"object","properties":{"b":{"type":"number"}}}}}}`))
	f.Add([]byte(`{"type":"unknown"}`))
	f.Add([]byte(`not json`))
	f.Add([]byte(`{`))
	f.Add([]byte(``))
	f.Add([]byte(`{"type":"integer","bitSize":999}`))

	f.Fuzz(func(t *testing.T, data []byte) {

		var s Schema

		// A parse error is an acceptable outcome; a panic is not.
		if err := s.UnmarshalJSON(data); err != nil {
			return
		}

		// A schema that parsed successfully must be safe to re-marshal.
		if _, err := s.MarshalJSON(); err != nil {
			t.Fatalf("UnmarshalJSON accepted %q but MarshalJSON failed: %v", data, err)
		}
	})
}

// FuzzSchemaPaths feeds arbitrary path strings to the path-walking entry points against a fixed,
// realistic schema and object. The danger zone is the string-cutting and array-index parsing in
// GetElement/Get/Set; none of it may panic regardless of how malformed the path is. Set's own
// recover() boundary is part of what this exercises.
func FuzzSchemaPaths(f *testing.F) {

	f.Add("name")
	f.Add("array.0")
	f.Add("array.999999999999999999999") // index overflow
	f.Add("array.-1")                    // negative index
	f.Add("array.notanumber")
	f.Add("")
	f.Add(".")
	f.Add("...")
	f.Add("name.extra.deep.path")
	f.Add("array.0.array.0") // nested descent

	schema := New(testStructA_Schema())

	f.Fuzz(func(t *testing.T, path string) {

		// GetElement walks the schema definition by path; must never panic.
		_, _ = schema.GetElement(path)

		// Get walks an actual object by path; must never panic.
		value := newTestStructA()
		_, _ = schema.Get(&value, path)

		// Set walks and mutates an object by path; must never panic (Set has a recover()
		// boundary, so a panic here would mean that boundary failed).
		_ = schema.Set(&value, path, "fuzz-value")
	})
}
