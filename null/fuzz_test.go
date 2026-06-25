package null

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzBool_UnmarshalJSON feeds arbitrary bytes to Bool.UnmarshalJSON and asserts that it never
// panics, and that any value it accepts marshals back to bytes that unmarshal to an equal value.
func FuzzBool_UnmarshalJSON(f *testing.F) {

	f.Add([]byte(`true`))
	f.Add([]byte(`false`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`"true"`))
	f.Add([]byte(`1`))

	f.Fuzz(func(t *testing.T, data []byte) {

		var value Bool
		if err := value.UnmarshalJSON(data); err != nil {
			return
		}

		assertRoundTrip(t, &value, &Bool{})
	})
}

// FuzzInt_UnmarshalJSON feeds arbitrary bytes to Int.UnmarshalJSON and asserts that it never
// panics, and that any value it accepts marshals back to bytes that unmarshal to an equal value.
func FuzzInt_UnmarshalJSON(f *testing.F) {

	f.Add([]byte(`0`))
	f.Add([]byte(`-1`))
	f.Add([]byte(`2147483647`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`1.5`))
	f.Add([]byte(`"123"`))

	f.Fuzz(func(t *testing.T, data []byte) {

		var value Int
		if err := value.UnmarshalJSON(data); err != nil {
			return
		}

		assertRoundTrip(t, &value, &Int{})
	})
}

// FuzzInt64_UnmarshalJSON feeds arbitrary bytes to Int64.UnmarshalJSON and asserts that it never
// panics, and that any value it accepts marshals back to bytes that unmarshal to an equal value.
func FuzzInt64_UnmarshalJSON(f *testing.F) {

	f.Add([]byte(`0`))
	f.Add([]byte(`-1`))
	f.Add([]byte(`9223372036854775807`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`1.5`))
	f.Add([]byte(`"123"`))

	f.Fuzz(func(t *testing.T, data []byte) {

		var value Int64
		if err := value.UnmarshalJSON(data); err != nil {
			return
		}

		assertRoundTrip(t, &value, &Int64{})
	})
}

// FuzzFloat_UnmarshalJSON feeds arbitrary bytes to Float.UnmarshalJSON and asserts that it never
// panics, and that any value it accepts marshals back to bytes that unmarshal to an equal value.
func FuzzFloat_UnmarshalJSON(f *testing.F) {

	f.Add([]byte(`0`))
	f.Add([]byte(`3.14`))
	f.Add([]byte(`-1e10`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`"123"`))

	f.Fuzz(func(t *testing.T, data []byte) {

		var value Float
		if err := value.UnmarshalJSON(data); err != nil {
			return
		}

		assertRoundTrip(t, &value, &Float{})
	})
}

// assertRoundTrip marshals an accepted value, unmarshals the result into a fresh
// zero value of the same type, and requires that the round-trip succeeds and is
// stable. NaN floats are skipped because they never compare equal to themselves.
func assertRoundTrip(t *testing.T, value, fresh json.Unmarshaler) {
	t.Helper()

	if f, ok := value.(*Float); ok && f.IsPresent() && f.Float() != f.Float() {
		return
	}

	marshaled, err := json.Marshal(value)
	require.NoError(t, err)

	require.NoError(t, fresh.UnmarshalJSON(marshaled))
	require.Equal(t, value, fresh)
}
