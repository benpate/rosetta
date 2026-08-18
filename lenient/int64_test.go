package lenient

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt64_UnmarshalJSON(t *testing.T) {

	test := func(name string, input string, expected int64, expectError bool) {
		t.Run(name, func(t *testing.T) {

			var value Int64
			err := json.Unmarshal([]byte(input), &value)

			if expectError {
				require.Error(t, err, "input %q should not parse", input)
				return
			}

			require.NoError(t, err, "input %q should parse", input)
			assert.Equal(t, expected, int64(value), "input %q", input)
		})
	}

	// The ugly real-world inputs named in the project plan
	test("plain integer", `480`, 480, false)
	test("quoted integer", `"480"`, 480, false)
	test("float", `480.0`, 480, false)
	test("empty string", `""`, 0, false)
	test("null", `null`, 0, false)

	// More corner cases
	test("negative integer", `-3`, -3, false)
	test("quoted negative", `"-3"`, -3, false)
	test("float truncates toward zero", `480.9`, 480, false)
	test("negative float truncates toward zero", `-480.9`, -480, false)
	test("quoted whitespace-padded", `"  480  "`, 480, false)
	test("whitespace-only string", `"   "`, 0, false)
	test("exponent", `1e3`, 1000, false)
	test("zero", `0`, 0, false)
	test("large number preserved", `2147483648`, 2147483648, false)
	test("large negative preserved", `-2147483648`, -2147483648, false)

	// Integers above 2^53 must be read from their source text. Decoding them
	// through a float64 rounds — 10000000000000001 came back as
	// 10000000000000000, which a fuzz round trip caught.
	test("exact above 2^53", `10000000000000001`, 10000000000000001, false)
	test("exact at 2^53 plus one", `9007199254740993`, 9007199254740993, false)
	test("exact negative above 2^53", `-10000000000000001`, -10000000000000001, false)
	test("max int64", `9223372036854775807`, math.MaxInt64, false)
	test("min int64", `-9223372036854775808`, math.MinInt64, false)
	test("quoted exact above 2^53", `"10000000000000001"`, 10000000000000001, false)

	// Beyond int64 there is no exact answer, so the clamping path takes over
	test("past max int64 clamps", `9223372036854775808`, math.MaxInt64, false)
	test("past min int64 clamps", `-9223372036854775809`, math.MinInt64, false)

	// Postel's law (via rosetta convert.Int64): garbage quietly becomes zero,
	// quoted floats are beyond ParseInt and zero out, huge values clamp
	test("quoted float zeroes", `"480.5"`, 0, false)
	test("quoted exponent zeroes", `"1e3"`, 0, false)
	test("non-numeric string zeroes", `"abc"`, 0, false)
	test("pixel suffix zeroes", `"480px"`, 0, false)
	test("boolean coerces", `true`, 1, false)
	test("single-element array unwraps", `[480]`, 480, false)
	test("object zeroes", `{}`, 0, false)
	test("huge exponent clamps", `1e300`, math.MaxInt64, false)

	// Only structurally invalid JSON is an error
	test("bare garbage", `!!!`, 0, true)
}

func TestInt64_UnmarshalXML(t *testing.T) {

	// wrapper receives the element under test
	type wrapper struct {
		Value Int64 `xml:"value"`
	}

	test := func(name string, input string, expected int64, expectError bool) {
		t.Run(name, func(t *testing.T) {

			var wrapped wrapper
			err := xml.Unmarshal([]byte("<wrapper><value>"+input+"</value></wrapper>"), &wrapped)

			if expectError {
				require.Error(t, err, "input %q should not parse", input)
				return
			}

			require.NoError(t, err, "input %q should parse", input)
			assert.Equal(t, expected, int64(wrapped.Value), "input %q", input)
		})
	}

	test("plain integer", `480`, 480, false)
	test("padded integer", `  480  `, 480, false)
	test("negative", `-3`, -3, false)
	test("empty element", ``, 0, false)
	test("whitespace only", `   `, 0, false)
	test("large number preserved", `2147483648`, 2147483648, false)

	// XML chardata goes through convert.Int64's string path (ParseInt), so
	// non-integer text quietly zeroes rather than erroring
	test("float zeroes", `480.9`, 0, false)
	test("exponent zeroes", `1e3`, 0, false)
	test("non-numeric zeroes", `abc`, 0, false)
}

func TestInt64_MarshalJSON(t *testing.T) {

	test := func(name string, input Int64, expected string) {
		t.Run(name, func(t *testing.T) {
			result, err := json.Marshal(input)
			require.NoError(t, err)
			assert.Equal(t, expected, string(result))
		})
	}

	// Output is always a plain JSON number
	test("positive", Int64(480), `480`)
	test("zero", Int64(0), `0`)
	test("negative", Int64(-3), `-3`)
}

func TestInt64_MarshalXML(t *testing.T) {

	type wrapper struct {
		Value Int64 `xml:"value"`
	}

	result, err := xml.Marshal(wrapper{Value: 480})
	require.NoError(t, err)
	assert.Equal(t, `<wrapper><value>480</value></wrapper>`, string(result))
}

func TestInt64_JSONRoundTrip(t *testing.T) {

	// A marshaled Int64 always re-parses to the same value
	for _, value := range []Int64{
		-100, -1, 0, 1, 480, 2147483647,
		math.MaxInt32 + 1, math.MaxInt64, math.MinInt64, 10000000000000001,
	} {

		data, err := json.Marshal(value)
		require.NoError(t, err)

		var parsed Int64
		require.NoError(t, json.Unmarshal(data, &parsed))
		assert.Equal(t, value, parsed)
	}
}

// TestInt64_RangeIsPlatformIndependent guards the reason this type is backed by
// int64 rather than int. A sender's JSON has no idea what GOARCH you built for,
// so the accepted range must not shrink on a 32-bit target (GOARCH=wasm is one).
// The static half of this guarantee is a `GOARCH=386 go vet ./...` pass, which
// this file would fail to even compile if Int64 were ever narrowed back to int.
func TestInt64_RangeIsPlatformIndependent(t *testing.T) {

	// The type is wide enough to hold the int64 extremes
	assert.Equal(t, int64(math.MaxInt64), int64(Int64(math.MaxInt64)))
	assert.Equal(t, int64(math.MinInt64), int64(Int64(math.MinInt64)))

	// ...and values beyond a 32-bit int survive a full decode/encode trip
	beyond32Bit := []string{`2147483648`, `-2147483649`, `4294967296`, `9223372036854775807`}

	for _, input := range beyond32Bit {

		var value Int64
		require.NoError(t, json.Unmarshal([]byte(input), &value), "input %q", input)

		data, err := json.Marshal(value)
		require.NoError(t, err)
		assert.Equal(t, input, string(data), "input %q must survive verbatim", input)
	}
}
