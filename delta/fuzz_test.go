package delta

import (
	"encoding/json"
	"testing"
)

// FuzzSliceUnmarshalJSON confirms that decoding arbitrary bytes into a Slice
// never panics. Invalid input should return an error, not crash.
func FuzzSliceUnmarshalJSON(f *testing.F) {

	f.Add([]byte("[1,2,3]"))
	f.Add([]byte("[]"))
	f.Add([]byte("null"))
	f.Add([]byte("not json"))
	f.Add([]byte(`["a","b"]`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var s Slice[int]
		_ = json.Unmarshal(data, &s)
	})
}
