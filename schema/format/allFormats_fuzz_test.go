package format

import (
	"strings"
	"testing"
)

// allGenerators returns every format Generator in this package, paired with a representative
// argument.  The schema package registers these by name; testing them here, generically, means
// a newly-added format is one line away from full property coverage.
func allGenerators() map[string]Generator {

	return map[string]Generator{
		"HTML":            HTML,
		"NoHTML":          NoHTML,
		"Text":            Text,
		"CSS":             CSS,
		"CSSDeclarations": CSSDeclarations,
		"ISO8601":         ISO8601,
		"Date":            Date,
		"DateTime":        DateTime,
		"Time":            Time,
		"Color":           Color,
		"Email":           Email,
		"IPv4":            IPv4,
		"IPv6":            IPv6,
		"Hostname":        Hostname,
		"URI":             URI,
		"URL":             URL,
		"Markdown":        Markdown,
		"HasLowercase":    HasLowercase,
		"HasUppercase":    HasUppercase,
		"HasNumbers":      HasNumbers,
		"ObjectID":        ObjectID,
		"Token":           Token,
		"MatchRegex":      MatchRegex,
		"UnsafeAny":       UnsafeAny,
		"In":              In,
		"NotIn":           NotIn,
		"Username":        Username,
		"WebFinger":       WebFinger,
	}
}

// formatArgs returns a representative configuration argument for each format.  Most ignore it;
// the set- and count-based formats do not.
func formatArgs() map[string]string {

	return map[string]string{
		"HasLowercase": "2",
		"HasUppercase": "2",
		"HasNumbers":   "2",
		"MatchRegex":   "^[a-z]+$",
		"In":           "red,green",
		"NotIn":        "red,green",
	}
}

// seedValues returns the hostile inputs shared by every target in this file.
func seedValues() []string {

	return []string{
		"",
		"a",
		"red",
		"blue",
		"2026-03-04",
		"2026-03-04T13:02:00Z",
		"2026-03-04T13:02:00+01:00",
		"#ff8800",
		"user@example.com",
		"@user@example.com",
		"192.168.1.1",
		"::1",
		"https://example.com/path?q=1",
		"000000000000000000000001",
		"<script>alert(1)</script>",
		"a\x00b",
		"\xff\xfe invalid utf8",
		strings.Repeat("a", 4096),
		"  leading and trailing  ",
		"line\nbreak\ttab",
		"h1 { color: red; }",
		"h1 {\n\tcolor: red;\n}",
	}
}

// FuzzFormats_NeverPanic feeds arbitrary values to EVERY format in this package.  These
// validators exist to be pointed at untrusted input, so a panic in any one of them is a
// denial-of-service in whatever is validating a request body.
func FuzzFormats_NeverPanic(f *testing.F) {

	for _, value := range seedValues() {
		f.Add(value)
	}

	generators := allGenerators()
	args := formatArgs()

	f.Fuzz(func(t *testing.T, value string) {

		for name, generator := range generators {

			// A panic here fails the fuzz target, which is the whole point
			validate := generator(args[name])
			_, _ = validate(value)
		}
	})
}

// FuzzFormats_AcceptedValueIsDerivedFromInput asserts the property that every format shares:
// a value it ACCEPTS must be returned as the value itself, or as some transformation of it --
// never as content invented from the format's own configuration.
//
// This is the property that caught NotIn returning its comma-separated option list in place of
// whatever the caller actually supplied, silently overwriting the field on every save.
func FuzzFormats_AcceptedValueIsDerivedFromInput(f *testing.F) {

	for _, value := range seedValues() {
		f.Add(value)
	}

	generators := allGenerators()
	args := formatArgs()

	// The formats that legitimately REWRITE their input rather than returning it verbatim.
	// Each strips or renders content, so the result is checked only for not being the config.
	transforms := map[string]bool{
		"HTML":            true,
		"NoHTML":          true,
		"Text":            true,
		"CSS":             true,
		"CSSDeclarations": true,
		"Markdown":        true,
		"URL":             true,
		"URI":             true,
		"Email":           true,
	}

	f.Fuzz(func(t *testing.T, value string) {

		for name, generator := range generators {

			arg := args[name]
			result, err := generator(arg)(value)

			// A rejected value tells us nothing about the result
			if err != nil {
				continue
			}

			// RULE: An accepted value must never come back as the format's own configuration.
			// An empty arg is skipped, because "" is a legitimate result for many formats.
			if (arg != "") && (result == arg) && (value != arg) {
				t.Fatalf("%s accepted %q but returned its own config %q", name, value, result)
			}

			if transforms[name] {
				continue
			}

			// Every other format either returns the value untouched, or rejects it
			if (result != value) && (result != "") {
				t.Fatalf("%s accepted %q but returned unrelated value %q", name, value, result)
			}
		}
	})
}

// FuzzFormats_Idempotent asserts that running an accepted value through its format a second
// time does not change it again.  A format that keeps rewriting is one that cannot be applied
// safely on both the read and write paths, which is exactly how stored values drift.
func FuzzFormats_Idempotent(f *testing.F) {

	for _, value := range seedValues() {
		f.Add(value)
	}

	generators := allGenerators()
	args := formatArgs()

	f.Fuzz(func(t *testing.T, value string) {

		for name, generator := range generators {

			validate := generator(args[name])

			once, err := validate(value)

			if err != nil {
				continue
			}

			twice, err := validate(once)

			if err != nil {
				t.Fatalf("%s accepted %q -> %q, then REJECTED its own output", name, value, once)
			}

			if once != twice {
				t.Fatalf("%s is not idempotent for %q: %q then %q", name, value, once, twice)
			}
		}
	})
}
