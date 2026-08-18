package null

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {

	var s String

	require.True(t, s.IsNull())
	require.False(t, s.IsPresent())
	require.Equal(t, "", s.String())
	require.Nil(t, s.Interface())

	s.Set("Baker Street")
	require.False(t, s.IsNull())
	require.True(t, s.IsPresent())
	require.Equal(t, "Baker Street", s.String())
	require.Equal(t, "Baker Street", s.Interface())

	s.Set("221B")
	require.False(t, s.IsNull())
	require.True(t, s.IsPresent())
	require.Equal(t, "221B", s.String())
	require.Equal(t, "221B", s.Interface())

	s.Unset()
	require.True(t, s.IsNull())
	require.False(t, s.IsPresent())
	require.Equal(t, "", s.String())
	require.Nil(t, s.Interface())
}

func TestNewString(t *testing.T) {

	s := NewString("")

	// A present-but-empty string is NOT null
	require.False(t, s.IsNull())
	require.True(t, s.IsPresent())
	require.Equal(t, "", s.String())
	require.Equal(t, "", s.Interface())

	s.Set("Baker Street")
	require.False(t, s.IsNull())
	require.True(t, s.IsPresent())
	require.Equal(t, "Baker Street", s.String())

	s.Unset()
	require.True(t, s.IsNull())
	require.False(t, s.IsPresent())
	require.Equal(t, "", s.String())
}

func TestString_IsNil(t *testing.T) {

	// IsNil must track IsNull exactly
	var value String
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Set("Watson")
	require.False(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Unset()
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())
}

func TestString_IsZero(t *testing.T) {

	// A null value is always zero
	var value String
	require.True(t, value.IsZero())

	// A present-but-empty value is ALSO zero
	value.Set("")
	require.True(t, value.IsZero())
	require.True(t, value.IsPresent())

	// A present, non-empty value is not
	value.Set("Watson")
	require.False(t, value.IsZero())

	value.Unset()
	require.True(t, value.IsZero())

	// The constructors agree
	require.True(t, NewString("").IsZero())
	require.False(t, NewString("Watson").IsZero())
}

func TestString_MarshalJSON(t *testing.T) {

	var value String

	// A null value marshals to the literal null
	result, err := json.Marshal(value)
	require.Nil(t, err)
	require.Equal(t, `null`, string(result))

	// A present-but-empty value marshals to an empty JSON string
	value.Set("")
	result, err = json.Marshal(value)
	require.Nil(t, err)
	require.Equal(t, `""`, string(result))

	// Escapes are encoding/json's, not strconv's — which means HTML-sensitive
	// runes are \u-escaped, exactly as they would be in a plain string field
	value.Set(`He said "no" <script>`)
	result, err = json.Marshal(value)
	require.Nil(t, err)
	require.Equal(t, `"He said \"no\" \u003cscript\u003e"`, string(result))
}

func TestString_UnmarshalJSON(t *testing.T) {

	var value String

	// Both empty bytes and the null literal read as null
	require.Nil(t, value.UnmarshalJSON([]byte(``)))
	require.True(t, value.IsNull())

	value.Set("Watson")
	require.Nil(t, value.UnmarshalJSON([]byte(`null`)))
	require.True(t, value.IsNull())

	// A JSON string is unquoted and unescaped
	require.Nil(t, value.UnmarshalJSON([]byte(`"221B <Baker>"`)))
	require.True(t, value.IsPresent())
	require.Equal(t, "221B <Baker>", value.String())

	// An empty JSON string is present, not null
	require.Nil(t, value.UnmarshalJSON([]byte(`""`)))
	require.True(t, value.IsPresent())
	require.True(t, value.IsZero())

	// RULE: a bare number or boolean is NOT coerced into text
	require.NotNil(t, value.UnmarshalJSON([]byte(`123`)))
	require.NotNil(t, value.UnmarshalJSON([]byte(`true`)))
	require.NotNil(t, value.UnmarshalJSON([]byte(`{"a":1}`)))
}

func TestString_Nullable(t *testing.T) {

	// String must satisfy the Nullable interface
	var value Nullable = String{}
	require.True(t, value.IsNull())

	value = NewString("")
	require.False(t, value.IsNull(), "a present-but-empty string is not null")
}

func TestString_MarshalJSON_InvalidUTF8(t *testing.T) {

	// GOTCHA: encoding/json replaces invalid UTF-8 with U+FFFD rather than failing,
	// so a String holding raw bytes does NOT survive a marshal/unmarshal round trip.
	value := NewString("\xff\xfe")

	result, err := value.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"��\"", string(result))

	var restored String
	require.Nil(t, restored.UnmarshalJSON(result))
	require.True(t, restored.IsPresent())
	require.Equal(t, "��", restored.String())
	require.NotEqual(t, value.String(), restored.String(), "the raw bytes are lost, not preserved")
}

func TestString_MarshalJSON_ControlCharacters(t *testing.T) {

	// Control characters are escaped, so the output is always a legal JSON string
	value := NewString("tab\tnewline\nnul\x00")

	result, err := value.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, `"tab\tnewline\nnul"`, string(result))

	// ...and unlike invalid UTF-8, they survive the round trip intact
	var restored String
	require.Nil(t, restored.UnmarshalJSON(result))
	require.Equal(t, value.String(), restored.String())
}

func TestString_MarshalJSON_Unicode(t *testing.T) {

	// Valid multi-byte runes pass through unescaped and round-trip exactly
	value := NewString("héllo wörld \U0001F389")

	result, err := value.MarshalJSON()
	require.Nil(t, err)
	require.Equal(t, "\"héllo wörld \U0001F389\"", string(result))

	var restored String
	require.Nil(t, restored.UnmarshalJSON(result))
	require.Equal(t, value.String(), restored.String())
}

func TestString_UnmarshalJSON_Escapes(t *testing.T) {

	var value String

	// Escape sequences are decoded, not preserved literally
	require.Nil(t, value.UnmarshalJSON([]byte(`"line\nbreak"`)))
	require.Equal(t, "line\nbreak", value.String())

	require.Nil(t, value.UnmarshalJSON([]byte(`"Aé"`)))
	require.Equal(t, "Aé", value.String())

	// An unpaired surrogate escape decodes to the replacement character
	require.Nil(t, value.UnmarshalJSON([]byte(`"\ud800"`)))
	require.Equal(t, "�", value.String())
}

func TestString_UnmarshalJSON_ErrorKeepsPriorValue(t *testing.T) {

	// A failed parse leaves the previous value untouched (matching the other null types)
	value := NewString("keep me")

	require.NotNil(t, value.UnmarshalJSON([]byte(`123`)))
	require.True(t, value.IsPresent())
	require.Equal(t, "keep me", value.String())
}

func TestString_UnmarshalJSON_PaddedNull(t *testing.T) {

	// GOTCHA: only the exact literal "null" is recognized.  A padded null is handed to
	// encoding/json, which rejects it as a string -- so unlike Object[T], String errors here.
	var value String

	require.NotNil(t, value.UnmarshalJSON([]byte(` null`)))
	require.True(t, value.IsNull())
}

func TestString_UnmarshalJSON_Whitespace(t *testing.T) {

	// Whitespace INSIDE the quotes is content, and is preserved exactly
	var value String

	require.Nil(t, value.UnmarshalJSON([]byte(`"   "`)))
	require.True(t, value.IsPresent())
	require.False(t, value.IsZero())
	require.Equal(t, "   ", value.String())
}
