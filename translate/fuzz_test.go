package translate

import "testing"

// FuzzNewFromJSON confirms that parsing arbitrary JSON strings into a Pipeline
// never panics. Malformed or unrecognized input should return an error.
func FuzzNewFromJSON(f *testing.F) {

	f.Add(`[{"path":"a","target":"b"}]`)
	f.Add(`[{"value":"v","target":"c"}]`)
	f.Add(`[{"expression":"{{.x}}","target":"y"}]`)
	f.Add(`[]`)
	f.Add(`not json`)
	f.Add(`[{"bogus":"data"}]`)

	f.Fuzz(func(t *testing.T, jsonString string) {
		_, _ = NewFromJSON(jsonString)
	})
}
