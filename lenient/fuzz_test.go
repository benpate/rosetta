package lenient

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

/******************************************
 * Shared Corpus
 ******************************************/

// jsonCorpus is every JSON shape worth throwing at either type: the
// documented behaviors, the numeric extremes, the encodings that have broken
// real parsers, and plain garbage.
var jsonCorpus = []string{

	// Well-formed scalars
	`480`, `"480"`, `480.0`, `-3`, `"-3"`, `0`, `-0`, `""`, `null`,
	`true`, `false`, `"true"`,

	// Numeric spellings that must survive as text
	`1.0`, `1.500`, `1e3`, `1E3`, `1e+3`, `1e-3`, `-1.0e-10`, `0.0000001`,

	// Precision and range boundaries
	`9223372036854775807`, `9223372036854775808`, `-9223372036854775808`,
	`-9223372036854775809`, `18446744073709551616`, `1e300`, `1e308`, `1e309`,
	`1e-300`, `1e400`, `-1e400`, `2147483647`, `2147483648`, `-2147483648`,
	`0.1`, `0.30000000000000004`, `123456789012345678901234567890`,

	// Whitespace and padding
	`"  480  "`, `"   "`, `" "`, "\t480\n", `  480  `,

	// Strings that look like other things
	`"480px"`, `"abc"`, `"1e3"`, `"480.5"`, `"NaN"`, `"Infinity"`, `"null"`,
	`"0x1F"`, `"+480"`, `"--3"`, `"1_000"`,

	// Escapes and Unicode
	`"café"`, `"🎉"`, `"\ud800"`, `"\udfff"`, `"\ud800\ud800"`,
	"\"a\\u0000b\"", `"\\"`, `"\""`, `"\n\r\t"`, "\"\\u202e\"",

	// Structures — rejected by String, coerced by Int64
	`{}`, `[]`, `[480]`, `[[480]]`, `{"a":1}`, `[null]`, `["480"]`,

	// Deep nesting (stack pressure)
	`[[[[[[[[[[1]]]]]]]]]]`,

	// Structurally invalid
	`!!!`, `NaN`, `Infinity`, `-Infinity`, `+1`, `01`, `.5`, `5.`, `1e`,
	`"abc`, `{`, `[`, ``, ` `, `nul`, `truex`, `"a" "b"`, `00`,

	// Bytes that are not text at all
	"\x00", "\xff\xfe", "\x00480",
}

// xmlCorpus is the chardata shapes worth throwing at the XML decoders,
// including fragments that try to escape their own element.
var xmlCorpus = []string{
	`480`, `  480  `, `-3`, ``, `   `, `abc`, `480.9`, `1e3`,
	`2147483648`, `9223372036854775808`, `1e400`,
	`a &amp; b`, `&lt;script&gt;`, `&#65;`, `&#x41;`, `&nosuchentity;`,
	`</value><value>1`, `<nested>1</nested>`, `<!--comment-->`,
	`<![CDATA[480]]>`, `]]>`, `&`, `<`,
	"\x00\xff", "\ufeff480", "480 ",
}

// seedJSON adds the shared JSON corpus to a fuzz target, as strings.
func seedJSON(f *testing.F) {
	f.Helper()
	for _, seed := range jsonCorpus {
		f.Add(seed)
	}
}

/******************************************
 * Universal Properties
 ******************************************/

// The properties below hold for BOTH types, so a new lenient type gets its
// whole safety net by being added to this one list.

// unmarshalJSON decodes into a fresh value of T and reports the result.
type jsonTarget struct {
	name string

	// decode parses data into a fresh value, returning its JSON re-encoding
	// and its Go-comparable form. A nil error means the decode succeeded.
	decode func(data []byte) (remarshaled []byte, value any, err error)
}

var jsonTargets = []jsonTarget{
	{
		name: "Int64",
		decode: func(data []byte) ([]byte, any, error) {
			var value Int64
			if err := json.Unmarshal(data, &value); err != nil {
				return nil, nil, err
			}
			remarshaled, err := json.Marshal(value)
			return remarshaled, value, err
		},
	},
	{
		name: "String",
		decode: func(data []byte) ([]byte, any, error) {
			var value String
			if err := json.Unmarshal(data, &value); err != nil {
				return nil, nil, err
			}
			remarshaled, err := json.Marshal(value)
			return remarshaled, value, err
		},
	},
}

// FuzzJSONProperties asserts the invariants every lenient type must hold:
// it never panics, a successful decode always re-encodes to valid JSON, and
// that re-encoding is a fixed point — decoding it again yields the same value.
// The fixed-point property is the one that matters: it means a value can make
// a round trip through storage or a proxy without drifting.
func FuzzJSONProperties(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, input string) {

		for _, target := range jsonTargets {

			// Property 1: decoding never panics (implicit — a panic fails here).
			remarshaled, value, err := target.decode([]byte(input))

			if err != nil {
				continue
			}

			// Property 2: a successful decode re-encodes to valid JSON.
			if !json.Valid(remarshaled) {
				t.Fatalf("%s: re-marshaled %q is not valid JSON (input %q)", target.name, remarshaled, input)
			}

			// Property 3: the re-encoding is a fixed point.
			secondPass, secondValue, err := target.decode(remarshaled)

			if err != nil {
				t.Fatalf("%s: re-marshaled %q failed to decode: %v (input %q)", target.name, remarshaled, err, input)
			}

			if secondValue != value {
				t.Fatalf("%s: value drifted on round trip: %#v -> %#v (input %q)", target.name, value, secondValue, input)
			}

			if string(secondPass) != string(remarshaled) {
				t.Fatalf("%s: encoding drifted on round trip: %q -> %q (input %q)", target.name, remarshaled, secondPass, input)
			}
		}
	})
}

// FuzzStructFields runs both types as fields of a real struct, which is how
// callers actually use them. A type can decode correctly in isolation and
// still break a document — a custom unmarshaler that consumes the wrong
// number of tokens desynchronizes the decoder for every field after it.
func FuzzStructFields(f *testing.F) {

	type document struct {
		Before  string `json:"before"`
		Number  Int64  `json:"number"`
		Middle  string `json:"middle"`
		Text    String `json:"text"`
		After   string `json:"after"`
		Missing Int64  `json:"missing"`
	}

	// Seed with documents that wrap the corpus in both tolerant positions
	for _, seed := range jsonCorpus {
		f.Add(`{"before":"a","number":` + seed + `,"middle":"b","after":"c"}`)
		f.Add(`{"before":"a","text":` + seed + `,"middle":"b","after":"c"}`)
	}

	f.Add(`{"number":1,"text":"a","number":2}`)
	f.Add(`{"NUMBER":1,"Text":"a"}`)
	f.Add(`{"before":"a"}`)

	f.Fuzz(func(t *testing.T, input string) {

		var parsed document

		if err := json.Unmarshal([]byte(input), &parsed); err != nil {
			return
		}

		// Property: an absent field is left at its zero value, never garbage
		// picked up from a neighbor.
		if !strings.Contains(input, "missing") && parsed.Missing != 0 {
			t.Fatalf("absent field was populated with %d (input %q)", parsed.Missing, input)
		}

		// Property: the surrounding fields still decode, so the tolerant
		// unmarshalers consumed exactly their own tokens.
		remarshaled, err := json.Marshal(parsed)

		if err != nil {
			t.Fatalf("marshal failed after successful unmarshal of %q: %v", input, err)
		}

		var second document

		if err := json.Unmarshal(remarshaled, &second); err != nil {
			t.Fatalf("re-marshaled %q failed to decode: %v (input %q)", remarshaled, err, input)
		}

		if second != parsed {
			t.Fatalf("document drifted on round trip: %#v -> %#v (input %q)", parsed, second, input)
		}
	})
}

/******************************************
 * Int64
 ******************************************/

// FuzzInt64_UnmarshalJSON pins the numeric properties that make Int64 safe to put
// in a size or offset: it never wraps, and it never flips sign.
func FuzzInt64_UnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, input string) {

		var value Int64

		if err := json.Unmarshal([]byte(input), &value); err != nil {
			return
		}

		// Property: a plain, in-range JSON integer parses to exactly itself.
		// This is the case callers actually depend on, so tolerance elsewhere
		// must never disturb it.
		if exact, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64); err == nil {
			if int64(value) != exact {
				t.Fatalf("exact integer %q parsed as %d", input, value)
			}
		}

		// Property: an unambiguously positive number never becomes negative.
		// Out-of-range values clamp, by policy — clamping must not wrap.
		if trimmed := strings.TrimSpace(input); strings.HasPrefix(trimmed, "1") || strings.HasPrefix(trimmed, "9") {
			if strings.ContainsAny(trimmed, "eE.") {
				return
			}
			if value < 0 {
				t.Fatalf("positive input %q became negative: %d", input, value)
			}
		}
	})
}

// FuzzInt64_UnmarshalXML fuzzes the XML path, which reaches convert.Int64 through
// chardata rather than a decoded JSON value and so has its own parse rules.
func FuzzInt64_UnmarshalXML(f *testing.F) {

	for _, seed := range xmlCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, chardata string) {

		type wrapper struct {
			Value Int64 `xml:"value"`
		}

		var wrapped wrapper

		if err := xml.Unmarshal([]byte("<wrapper><value>"+chardata+"</value></wrapper>"), &wrapped); err != nil {
			return
		}

		// Property: a successful decode re-marshals, and the value survives.
		remarshaled, err := xml.Marshal(wrapped)

		if err != nil {
			t.Fatalf("marshal failed after successful unmarshal of %q: %v", chardata, err)
		}

		var second wrapper

		if err := xml.Unmarshal(remarshaled, &second); err != nil {
			t.Fatalf("re-marshaled %q failed to decode: %v (input %q)", remarshaled, err, chardata)
		}

		if second.Value != wrapped.Value {
			t.Fatalf("value drifted on round trip: %d -> %d (input %q)", wrapped.Value, second.Value, chardata)
		}
	})
}

// FuzzInt64_RoundTrip drives the encode side from arbitrary integers, which the
// decode-first targets cannot reach: json.Marshal must emit something Int64 can
// read back exactly, for every value in the type's range.
func FuzzInt64_RoundTrip(f *testing.F) {

	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(480))
	f.Add(int64(math.MaxInt64))
	f.Add(int64(math.MinInt64))
	f.Add(int64(math.MaxInt32))

	f.Fuzz(func(t *testing.T, seed int64) {

		value := Int64(seed)

		data, err := json.Marshal(value)

		if err != nil {
			t.Fatalf("marshal of %d failed: %v", value, err)
		}

		var parsed Int64

		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("re-parsing own output %q failed: %v", data, err)
		}

		if parsed != value {
			t.Fatalf("round trip changed %d to %d (encoded as %q)", value, parsed, data)
		}
	})
}

/******************************************
 * String
 ******************************************/

// FuzzString_NumberFidelity is the reason String decodes with UseNumber. A
// number arriving in a string field must keep its SOURCE TEXT — the oEmbed
// version "1.0" may not become "1", and no long decimal may lose digits to a
// float64 round trip. This property is what a convert.String implementation
// would fail.
func FuzzString_NumberFidelity(f *testing.F) {

	for _, seed := range jsonCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {

		// Only bare JSON numbers are in scope here. Validity is checked on the
		// UNTRIMMED input: strings.TrimSpace strips Unicode whitespace, which
		// includes form feed and vertical tab, but JSON's whitespace is only
		// space, tab, LF, and CR. Trimming first would let "0\f" through as a
		// "valid number" that the decoder then rightly rejects.
		if !json.Valid([]byte(input)) {
			return
		}

		trimmed := strings.TrimSpace(input)

		if _, err := strconv.ParseFloat(trimmed, 64); err != nil {
			return
		}

		var value String

		if err := json.Unmarshal([]byte(input), &value); err != nil {
			t.Fatalf("valid JSON number %q failed to parse: %v", input, err)
		}

		// Property: the source text is preserved exactly, digit for digit.
		if string(value) != trimmed {
			t.Fatalf("number %q was rewritten as %q", trimmed, value)
		}
	})
}

// FuzzString_UnmarshalJSON pins the properties that keep String honest: it
// never invents content, and it never emits an invalid UTF-8 sequence that
// json.Marshal would then silently replace.
func FuzzString_UnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, input string) {

		var value String

		if err := json.Unmarshal([]byte(input), &value); err != nil {

			// Property: a failed decode leaves the value untouched.
			if value != "" {
				t.Fatalf("failed decode of %q still wrote %q", input, value)
			}
			return
		}

		// Property: a decoded value is always valid UTF-8. Go's JSON decoder
		// replaces malformed escapes with U+FFFD, so anything else here would
		// mean we assembled the string ourselves and got it wrong.
		if !utf8.ValidString(string(value)) {
			t.Fatalf("decoded %q produced invalid UTF-8: %q", input, value)
		}

		// Property: every decoded rune is backed by at least one input byte, so
		// the decoder cannot fabricate content. Counted in RUNES, not bytes:
		// encoding/json substitutes a 3-byte U+FFFD for each malformed byte, so
		// `"\xa8\xa8"` legitimately decodes 4 bytes into 6 — a byte-length
		// comparison would flag standard-library behavior as a defect.
		if utf8.RuneCountInString(string(value)) > len(input) {
			t.Fatalf("decoded %q (%d bytes) into %d runes: %q", input, len(input), utf8.RuneCountInString(string(value)), value)
		}
	})
}

// FuzzString_XML covers the XML path, which String leaves to the standard
// library on purpose — the test exists to catch the day someone adds a custom
// UnmarshalXML and changes that.
func FuzzString_XML(f *testing.F) {

	for _, seed := range xmlCorpus {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, chardata string) {

		type wrapper struct {
			Value String `xml:"value"`
		}

		var wrapped wrapper

		if err := xml.Unmarshal([]byte("<wrapper><value>"+chardata+"</value></wrapper>"), &wrapped); err != nil {
			return
		}

		remarshaled, err := xml.Marshal(wrapped)

		if err != nil {
			t.Fatalf("marshal failed after successful unmarshal of %q: %v", chardata, err)
		}

		var second wrapper

		if err := xml.Unmarshal(remarshaled, &second); err != nil {
			t.Fatalf("re-marshaled %q failed to decode: %v (input %q)", remarshaled, err, chardata)
		}

		// Property: XML round trips exactly — escaping in, unescaping out.
		if second.Value != wrapped.Value {
			t.Fatalf("value drifted on round trip: %q -> %q (input %q)", wrapped.Value, second.Value, chardata)
		}
	})
}

// FuzzString_RoundTrip drives the encode side from arbitrary strings: control
// characters, lone surrogates, and invalid UTF-8 all have to survive being
// marshaled and read back.
func FuzzString_RoundTrip(f *testing.F) {

	f.Add("")
	f.Add("1.0")
	f.Add(`say "hi"`)
	f.Add("café — 🎉")
	f.Add("  padded  ")
	f.Add("<script>alert(1)</script>")
	f.Add("\x00\x01\x02")
	f.Add("\xff\xfe")
	f.Add("line\nbreak")

	f.Fuzz(func(t *testing.T, seed string) {

		value := String(seed)

		data, err := json.Marshal(value)

		if err != nil {
			t.Fatalf("marshal of %q failed: %v", seed, err)
		}

		var parsed String

		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("re-parsing own output %q failed: %v", data, err)
		}

		// json.Marshal replaces invalid UTF-8 with U+FFFD, so an exact match
		// is only required of input that was valid UTF-8 to begin with.
		if utf8.ValidString(seed) && parsed != value {
			t.Fatalf("round trip changed %q to %q (encoded as %q)", value, parsed, data)
		}

		// Property: whatever survives the first trip is stable from here on.
		// The comparison is on VALUES, not bytes: json.Marshal escapes the
		// U+FFFD it substitutes for invalid UTF-8 ("�") but emits a
		// genuine U+FFFD literally, so the encoding of a repaired string
		// legitimately differs from the encoding of the broken original.
		second, err := json.Marshal(parsed)

		if err != nil {
			t.Fatalf("second marshal of %q failed: %v", parsed, err)
		}

		var third String

		if err := json.Unmarshal(second, &third); err != nil {
			t.Fatalf("re-parsing own output %q failed: %v", second, err)
		}

		if third != parsed {
			t.Fatalf("value is not a fixed point: %q -> %q", parsed, third)
		}
	})
}
