package lenient

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString_UnmarshalJSON(t *testing.T) {

	test := func(name string, input string, expected string, expectError bool) {
		t.Run(name, func(t *testing.T) {

			var value String
			err := json.Unmarshal([]byte(input), &value)

			if expectError {
				require.Error(t, err, "input %q should not parse", input)
				return
			}

			require.NoError(t, err, "input %q should parse", input)
			assert.Equal(t, expected, string(value), "input %q", input)
		})
	}

	// The ordinary case
	test("quoted string", `"1.0"`, "1.0", false)
	test("empty string", `""`, "", false)
	test("null", `null`, "", false)
	test("whitespace preserved", `"  padded  "`, "  padded  ", false)
	test("unicode", `"café — 🎉"`, "café — 🎉", false)
	test("escaped quotes", `"say \"hi\""`, `say "hi"`, false)

	// Postel's law: a scalar in the wrong shape is stringified, and numbers
	// keep their exact source text — this is why SoundCloud's 1.0 stays "1.0"
	test("float keeps source text", `1.0`, "1.0", false)
	test("integer", `1`, "1", false)
	test("trailing zeroes preserved", `1.500`, "1.500", false)
	test("negative", `-3`, "-3", false)
	test("exponent keeps source text", `1e3`, "1e3", false)
	test("huge number is not clamped", `1e400`, "1e400", false)
	test("boolean true", `true`, "true", false)
	test("boolean false", `false`, "false", false)

	// RULE: objects and arrays are structurally wrong, not merely mistyped
	test("object", `{"nested": true}`, "", true)
	test("empty object", `{}`, "", true)
	test("array", `["a"]`, "", true)
	test("empty array", `[]`, "", true)

	// Structurally invalid JSON is an error
	test("bare garbage", `!!!`, "", true)
	test("unterminated string", `"abc`, "", true)
}

func TestString_MarshalJSON(t *testing.T) {

	test := func(name string, input String, expected string) {
		t.Run(name, func(t *testing.T) {
			result, err := json.Marshal(input)
			require.NoError(t, err)
			assert.Equal(t, expected, string(result))
		})
	}

	// Output is always a plain, escaped JSON string
	test("plain", String("1.0"), `"1.0"`)
	test("empty", String(""), `""`)
	test("quotes escaped", String(`say "hi"`), `"say \"hi\""`)
	test("unicode", String("café"), `"café"`)
}

func TestString_XML(t *testing.T) {

	// String carries no custom XML methods; the default chardata handling is
	// already correct, including whitespace, which a title may legitimately own
	type wrapper struct {
		Value String `xml:"value"`
	}

	test := func(name string, input string, expected string) {
		t.Run(name, func(t *testing.T) {

			var wrapped wrapper
			require.NoError(t, xml.Unmarshal([]byte("<wrapper><value>"+input+"</value></wrapper>"), &wrapped))
			assert.Equal(t, expected, string(wrapped.Value), "input %q", input)
		})
	}

	test("plain text", `hello`, "hello")
	test("padding preserved", `  padded  `, "  padded  ")
	test("empty element", ``, "")
	test("entities decoded", `a &amp; b`, "a & b")
}

func TestString_JSONRoundTrip(t *testing.T) {

	// A marshaled String always re-parses to the same value
	for _, value := range []String{"", "1.0", `say "hi"`, "café — 🎉", "  padded  ", "<script>"} {

		data, err := json.Marshal(value)
		require.NoError(t, err)

		var parsed String
		require.NoError(t, json.Unmarshal(data, &parsed))
		assert.Equal(t, value, parsed)
	}
}
