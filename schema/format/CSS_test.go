package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCSS verifies that ordinary stylesheet constructs survive sanitizing.
func TestCSS(t *testing.T) {

	validate := CSS("")

	table := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "simple rule",
			value:    "body { color: red; }",
			expected: "body { color: red; }\n",
		},
		{
			name:     "multiple declarations",
			value:    "p { color: #336699; margin-top: 1.5em; }",
			expected: "p { color: #336699; margin-top: 1.5em; }\n",
		},
		{
			name:     "multiple rules",
			value:    "h1 { font-size: 2em; } h2 { font-size: 1.5em; }",
			expected: "h1 { font-size: 2em; }\nh2 { font-size: 1.5em; }\n",
		},
		{
			name:     "class, id, and descendant selectors",
			value:    "#main .article > p:first-child { line-height: 1.6; }",
			expected: "#main .article > p:first-child { line-height: 1.6; }\n",
		},
		{
			name:     "attribute selector",
			value:    "a[target] { text-decoration: underline; }",
			expected: "a[target] { text-decoration: underline; }\n",
		},
		{
			name:     "grouped selectors",
			value:    "h1, h2, h3 { font-weight: 700; }",
			expected: "h1, h2, h3 { font-weight: 700; }\n",
		},
		{
			name:     "comments are removed",
			value:    "body { /* a note */ color: red; }",
			expected: "body { color: red; }\n",
		},
		{
			name:     "flex layout",
			value:    ".row { display: flex; gap: 8px; justify-content: space-between; }",
			expected: ".row { display: flex; gap: 8px; justify-content: space-between; }\n",
		},
		{
			name:     "value with a slash",
			value:    "p { font: 12px/1.5 sans-serif; }",
			expected: "p { font: 12px/1.5 sans-serif; }\n",
		},
		{
			name:     "media query is kept and its body sanitized",
			value:    "@media screen and (max-width: 600px) { body { color: red; } }",
			expected: "@media screen and (max-width: 600px) {\nbody { color: red; }\n}\n",
		},
		{
			name:     "supports at-rule",
			value:    "@supports (display: grid) { .row { display: grid; } }",
			expected: "@supports (display: grid) {\n.row { display: grid; }\n}\n",
		},
		{
			name:     "empty stylesheet",
			value:    "",
			expected: "",
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := validate(testCase.value)
			require.Nil(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}
}

// TestCSS_MultiLine verifies that a stylesheet formatted the way people
// actually write one -- declarations indented on their own lines, closing brace
// on the last -- survives sanitizing.
//
// Every positive case in TestCSS above is written on a single line, and that is
// exactly how a stylesheet ending `red;\n}` came to sanitize away to nothing
// while the suite stayed green.  See normalizeCSSDeclarations.
func TestCSS_MultiLine(t *testing.T) {

	validate := CSS("")

	table := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "one declaration, closing brace on its own line",
			value:    "h1 {\n\tcolor: red;\n}",
			expected: "h1 { color: red; }\n",
		},
		{
			name:     "several declarations",
			value:    ".article p {\n\tline-height: 1.6;\n\tcolor: #333;\n}",
			expected: ".article p { line-height: 1.6; color: #333; }\n",
		},
		{
			name:     "several rules separated by a blank line",
			value:    "h1 {\n\tfont-size: 2em;\n}\n\nh2 {\n\tfont-size: 1.5em;\n}\n",
			expected: "h1 { font-size: 2em; }\nh2 { font-size: 1.5em; }\n",
		},
		{
			name:     "windows line endings",
			value:    "h1 {\r\n\tcolor: red;\r\n}\r\n",
			expected: "h1 { color: red; }\n",
		},
		{
			name:     "value wrapped across two lines",
			value:    "h1 {\n\tbox-shadow: 0 0 4px\n\t\t#000;\n}",
			expected: "h1 { box-shadow: 0 0 4px #000; }\n",
		},
		{
			name:     "indented media query",
			value:    "@media (max-width: 40em) {\n\th1 {\n\t\tfont-size: 1.2em;\n\t}\n}",
			expected: "@media (max-width: 40em) {\nh1 { font-size: 1.2em; }\n}\n",
		},
		{
			name:     "grouped selectors on separate lines",
			value:    "h1,\nh2 {\n\tcolor: red;\n}",
			expected: "h1,\nh2 { color: red; }\n",
		},
		{
			name:     "comment on its own line",
			value:    "/* headings */\nh1 {\n\tcolor: red;\n}",
			expected: "h1 { color: red; }\n",
		},
		{
			name:     "form feed is whitespace, not a reason to drop the rule",
			value:    "h1\f{\f\fcolor: red;\f}",
			expected: "h1 { color: red; }\n",
		},
		{
			name:     "empty and repeated declaration separators",
			value:    "h1 {\n\t;;\n\tcolor: red;;\n}",
			expected: "h1 { color: red; }\n",
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := validate(testCase.value)
			require.Nil(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}
}

// TestCSS_WhitespaceIsInsignificant verifies the general property behind
// TestCSS_MultiLine: re-indenting a stylesheet must not change how many rules
// survive it.  CSS treats every run of whitespace between tokens as one space,
// so a sanitizer that disagrees is destroying valid input.
func TestCSS_WhitespaceIsInsignificant(t *testing.T) {

	validate := CSS("")

	for _, value := range []string{
		"h1 { color: red; }",
		"h1 { color: red; font-size: 2em; }",
		"h1 { color: red; } h2 { color: blue; }",
		"@media screen { h1 { color: red; } }",
		"#main .article > p:first-child { line-height: 1.6; margin: 0 auto; }",
		"h1 { behavior: url(x.htc); }",
		"h1 { color: red; position: fixed; }",
	} {
		expected := countCSSBlocks(t, validate, value)

		for _, whitespace := range []string{" ", "\t", "\n", "\r\n", "\f", " \n\t\t", "\n\n\n"} {

			variant := respaceCSS(value, whitespace)
			actual := countCSSBlocks(t, validate, variant)

			require.Equal(t, expected, actual, "re-indenting %q with %q changed how many rules survived", value, whitespace)
		}
	}
}

// TestCSS_NormalizeDeclarations covers the declaration-list normalizer directly.
func TestCSS_NormalizeDeclarations(t *testing.T) {

	table := []struct {
		name     string
		value    string
		expected string
	}{
		{"empty", "", ""},
		{"only whitespace", " \t\r\n\f", ""},
		{"only separators", ";;;", ""},
		{"already normal", "color: red", "color: red"},
		{"trailing separator", "color: red;", "color: red"},
		{"trailing newline after separator", "color: red;\n", "color: red"},
		{"leading separator", "; color: red", "color: red"},
		{"indented block", "\n\tcolor: red;\n\tmargin: 0;\n", "color: red; margin: 0"},
		{"value wrapped across lines", "box-shadow: 0 0 4px\n\t#000", "box-shadow: 0 0 4px #000"},
		{"repeated separators", "color: red;;;margin: 0", "color: red; margin: 0"},
		{"whitespace around the colon", "color\n:\nred", "color : red"},
		{"unsafe characters are left for the checks that own them", `font-family: "X"`, `font-family: "X"`},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.expected, normalizeCSSDeclarations(testCase.value))
		})
	}
}

// countCSSBlocks returns the number of rule blocks a stylesheet sanitizes down
// to, which is the measure of how much of it survived.
func countCSSBlocks(t *testing.T, validate StringFormat, value string) int {

	t.Helper()

	result, err := validate(value)
	require.Nil(t, err)

	return strings.Count(result, "{")
}

// respaceCSS replaces every run of CSS whitespace in a stylesheet with the given
// run.  Whitespace is only ever exchanged for other whitespace, never removed,
// so the result means exactly what the original meant.
func respaceCSS(value string, whitespace string) string {

	var result strings.Builder

	for index := 0; index < len(value); index++ {

		if !strings.ContainsRune(cssWhitespace, rune(value[index])) {
			result.WriteByte(value[index])
			continue
		}

		// Consume the whole run, so the replacement is one-for-one
		for (index+1 < len(value)) && strings.ContainsRune(cssWhitespace, rune(value[index+1])) {
			index++
		}

		result.WriteString(whitespace)
	}

	return result.String()
}

// TestCSS_Unsafe verifies that every construct capable of loading a resource,
// executing script, or escaping the stylesheet is removed.
func TestCSS_Unsafe(t *testing.T) {

	validate := CSS("")

	table := []struct {
		name  string
		value string
	}{
		{"import at-rule", `@import url("//evil.example.com/x.css");`},
		{"import without url function", `@import "//evil.example.com/x.css";`},
		{"font-face fetches a remote font", `@font-face { font-family: x; src: url(//evil.example.com/f.woff); }`},
		{"charset at-rule", `@charset "utf-8";`},
		{"namespace at-rule", `@namespace svg "http://www.w3.org/2000/svg";`},
		{"unknown at-rule", `@evil { body { color: red; } }`},
		{"quoted at-rule condition", `@media "</style><script>alert(1)</script>" { body { color: red; } }`},
		{"style tag breakout in a selector", `</style><script>alert(1)</script> { color: red; }`},
		{"style tag breakout in a value", `body { color: red</style><script>alert(1)</script>; }`},
		{"expression value", `body { color: expression(alert(1)); }`},
		{"url value", `body { background-color: url(//evil.example.com/x.png); }`},
		{"javascript url value", `body { background-color: url(javascript:alert(1)); }`},
		{"behavior property", `body { behavior: url(//evil.example.com/x.htc); }`},
		{"moz-binding property", `body { -moz-binding: url(//evil.example.com/x.xml); }`},
		{"content property", `body:before { content: "injected"; }`},
		{"fixed position overlay", `body { position: fixed; }`},
		{"absolute position overlay", `body { position: absolute; }`},
		{"css escape sequence", `body { color: \65 xpression(alert(1)); }`},
		{"comment splices a value", `body { color: re/**/d; }`},
		{"unterminated comment", `body { color: red; } /* h1 { color: blue; }`},
		{"unbalanced block", `body { color: red;`},
		{"quote in a value", `body { font-family: "Helvetica Neue"; }`},
		{"backslash in a selector", `bo\64 y { color: red; }`},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {

			result, err := validate(testCase.value)
			require.Nil(t, err)

			// Whatever survives must not carry the construct that made the input
			// dangerous, in any spelling.
			lower := strings.ToLower(result)

			for _, forbidden := range []string{"url", "expression", "javascript", "@import", "@font-face", "behavior", "binding", "content", "position", "<", ">", `"`, "'", `\`} {
				require.NotContains(t, lower, forbidden, "input %q produced %q", testCase.value, result)
			}
		})
	}
}

// TestCSS_PartialSanitizing verifies that an unsafe declaration is dropped
// without taking the safe declarations around it with it.
func TestCSS_PartialSanitizing(t *testing.T) {

	validate := CSS("")

	result, err := validate("body { color: red; behavior: url(x.htc); margin: 0; }")

	require.Nil(t, err)
	require.Contains(t, result, "color: red")
	require.Contains(t, result, "margin: 0")
	require.NotContains(t, result, "behavior")
}

// TestCSS_NestingDepth verifies that at-rule nesting past the cap is discarded
// rather than driving unbounded recursion.
func TestCSS_NestingDepth(t *testing.T) {

	validate := CSS("")

	deep := strings.Repeat("@media screen { ", 50) + "body { color: red; }" + strings.Repeat(" }", 50)

	result, err := validate(deep)

	require.Nil(t, err)
	require.Equal(t, "", result, "nesting past the cap should sanitize away entirely")
}

// TestCSS_EmptyRulesRemoved verifies that a rule whose declarations all
// sanitized away is not emitted as an empty shell.
func TestCSS_EmptyRulesRemoved(t *testing.T) {

	validate := CSS("")

	result, err := validate("body { behavior: url(x.htc); }")

	require.Nil(t, err)
	require.Equal(t, "", result)
}

// TestCSS_Idempotent verifies that sanitizing an already-sanitized stylesheet
// changes nothing. A sanitizer that is not idempotent is either destroying safe
// input or failing to reach a fixed point on unsafe input.
func TestCSS_Idempotent(t *testing.T) {

	validate := CSS("")

	for _, value := range []string{
		"body { color: red; }",
		"@media screen and (max-width: 600px) { body { color: red; } }",
		"#main .article > p { line-height: 1.6; }",
		`@import url("//evil.example.com/x.css"); body { color: red; }`,
	} {
		once, err := validate(value)
		require.Nil(t, err)

		twice, err := validate(once)
		require.Nil(t, err)

		require.Equal(t, once, twice, "sanitizing %q was not idempotent", value)
	}
}

// TestCSS_CaseInsensitive verifies that CSS's case-insensitive spellings are
// recognized, and that an at-rule name is normalized to the form that was
// actually checked against the allowlist.
func TestCSS_CaseInsensitive(t *testing.T) {

	validate := CSS("")

	result, err := validate("BODY { COLOR: RED; }")
	require.Nil(t, err)
	require.Equal(t, "BODY { COLOR: RED; }\n", result)

	result, err = validate("@MEDIA screen { body { color: red; } }")
	require.Nil(t, err)
	require.Equal(t, "@media screen {\nbody { color: red; }\n}\n", result)
}

// TestCSS_MalformedAtRules verifies the at-rule paths that carry no usable
// condition.
func TestCSS_MalformedAtRules(t *testing.T) {

	validate := CSS("")

	for _, value := range []string{
		"@media { body { color: red; } }",
		"@media",
		"@supports { body { color: red; } }",
		"@media screen { }",
		"@media screen { body { behavior: url(x.htc); } }",
	} {
		result, err := validate(value)
		require.Nil(t, err)
		require.Equal(t, "", result, "malformed at-rule %q should sanitize away", value)
	}
}

// TestCSS_Helpers covers the guard clauses in the sanitizer's helpers that the
// stylesheet-level tests cannot reach directly.
func TestCSS_Helpers(t *testing.T) {

	t.Run("splitCSSBlock requires a leading brace", func(t *testing.T) {
		_, _, ok := splitCSSBlock("body { color: red; }")
		require.False(t, ok)
	})

	t.Run("splitCSSBlock reports an unbalanced block", func(t *testing.T) {
		_, _, ok := splitCSSBlock("{ color: red;")
		require.False(t, ok)
	})

	t.Run("splitCSSBlock returns the remainder", func(t *testing.T) {
		body, rest, ok := splitCSSBlock("{ a { b } } tail")
		require.True(t, ok)
		require.Equal(t, " a { b } ", body)
		require.Equal(t, " tail", rest)
	})

	t.Run("empty value is unsafe", func(t *testing.T) {
		require.False(t, isSafeCSSValue(""))
	})

	t.Run("empty prelude is unsafe", func(t *testing.T) {
		require.False(t, isSafeCSSPrelude(""))
	})

	t.Run("empty declarations sanitize to empty", func(t *testing.T) {
		require.Equal(t, "", sanitizeDeclarations(""))
	})

	// Bluemonday lowercases a value before handing it to the policy handler, so
	// an uppercase value only reaches isSafeCSSValue on a direct call.
	t.Run("uppercase values are accepted on their merits", func(t *testing.T) {
		require.True(t, isSafeCSSValue("RED"))
		require.False(t, isSafeCSSValue("URL(x)"))
	})

	t.Run("splitCSSAtRule handles a missing condition", func(t *testing.T) {
		name, condition := splitCSSAtRule("@media")
		require.Equal(t, "@media", name)
		require.Equal(t, "", condition)
	})

	t.Run("splitCSSAtRule splits on a parenthesis", func(t *testing.T) {
		name, condition := splitCSSAtRule("@media(min-width: 40em)")
		require.Equal(t, "@media", name)
		require.Equal(t, "(min-width: 40em)", condition)
	})
}

// TestCSSDeclarations verifies the declaration-list format: the same allowlists as
// CSS, applied to the contents of a `style` attribute instead of a stylesheet.
func TestCSSDeclarations(t *testing.T) {

	validate := CSSDeclarations("")

	table := []struct {
		name     string
		value    string
		expected string
	}{
		// Values come back with a trailing ";" -- bluemonday emits one, and it is
		// valid in a style attribute, so it is left rather than trimmed.
		{"empty", "", ""},
		{"single declaration", "color:red", "color: red;"},
		{"several declarations", "color:red; padding:8px", "color: red; padding: 8px;"},
		{"layout properties survive", "display:flex; gap:8px", "display: flex; gap: 8px;"},
		{"written across lines", "color: red;\n\tmargin: 0;\n", "color: red; margin: 0;"},

		// RULE: the properties CSS refuses in a stylesheet are refused here too. A
		// `style` attribute is injected into a page its author does not own, exactly
		// as a stylesheet is, so the two cannot disagree about what is safe.
		{"fixed positioning is dropped", "position:fixed; color:red", "color: red;"},
		{"in-flow positioning survives", "position:relative; color:red", "position: relative; color: red;"},
		{"resource loading is dropped", "background:url(https://example.com/x); color:red", "color: red;"},
		{"legacy script hooks are dropped", "behavior:url(x.htc); color:red", "color: red;"},
		{"unknown properties are dropped", "made-up:1; color:red", "color: red;"},

		// A quote is what the sanitizer's own round trip is delimited by, so a value
		// carrying one cannot be partly recovered -- the whole list is discarded.
		{"a quote discards the whole list", `color:red" onmouseover="alert(1)`, ""},
		{"a single quote discards it too", "font-family:'X'; color:red", ""},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := validate(testCase.value)
			require.Nil(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}
}

// TestCSSDeclarations_NotInterchangeableWithCSS pins the reason this format exists.
//
// The two formats read the same input differently, and each destroys what the other
// accepts: CSS looks for selectors and blocks, so a declaration list has no rule in
// it and vanishes; CSSDeclarations looks for declarations, so a stylesheet's braces
// and selector are not properties and vanish too. Pointing a schema at the wrong one
// does not error -- it silently stores an empty string.
func TestCSSDeclarations_NotInterchangeableWithCSS(t *testing.T) {

	const declarations = "color:red; padding:8px"
	const stylesheet = ".foo { color:red }"

	declarationsFormat := CSSDeclarations("")
	stylesheetFormat := CSS("")

	// Each format keeps what it is for
	result, err := declarationsFormat(declarations)
	require.Nil(t, err)
	require.Equal(t, "color: red; padding: 8px;", result)

	result, err = stylesheetFormat(stylesheet)
	require.Nil(t, err)
	require.Contains(t, result, "color: red")

	// ...and empties what it is not
	result, err = stylesheetFormat(declarations)
	require.Nil(t, err)
	require.Equal(t, "", result)

	result, err = declarationsFormat(stylesheet)
	require.Nil(t, err)
	require.Equal(t, "", result)
}
