package format

import (
	"fmt"
	"strings"
	"testing"
)

// cssSeeds returns the stylesheets shared by the targets in this file: ordinary
// rules in both the single-line and the indented spelling, every construct the
// sanitizer is meant to remove, and the malformed shapes that exercise its
// structural walk.
func cssSeeds() []string {

	return []string{
		"",
		" ",
		"\n\t\r\f",
		"body { color: red; }",
		"body {\n\tcolor: red;\n}",
		"body {\r\n\tcolor: red;\r\n}\r\n",
		".a,\n.b {\n\tmargin: 0 auto;\n\tcolor: #336699;\n}\n\nh2 {\n\tfont-size: 2em;\n}",
		"h1 {\n\tbox-shadow: 0 0 4px\n\t\t#000;\n}",
		"@media screen { body { color: red; } }",
		"@media (max-width: 40em) {\n\th1 {\n\t\tfont-size: 1.2em;\n\t}\n}",
		"@supports (display: grid) { .row { display: grid; } }",
		`@import url("//evil.example.com/x.css");`,
		"@font-face { font-family: x; src: url(//evil.example.com/f.woff); }",
		"body { color: expression(alert(1)); }",
		"body { position: fixed; }",
		"body:before { content: 'x'; }",
		"</style><script>alert(1)</script>",
		"body { color: red",
		"{{{{{{",
		"}}}}}}",
		"/*",
		"/* unterminated { color: red; }",
		"body { color: re/**/d; }",
		"body { background-color: url(x) }",
		"body { ;;; color: red ;;; }",
		"@media screen { @media screen { @media screen { @media screen { @media screen { body { color: red; } } } } } }",
		"body { color: \x00red; }",
		"a\\3c /style\\3e { color: red; }",
		"\xff\xfe invalid utf8 { color: red; }",
		strings.Repeat("a", 2048) + " { color: red; }",
	}
}

// FuzzCSS feeds arbitrary strings to the stylesheet sanitizer to confirm that it
// never panics and that its output can never carry the constructs that would let
// a stylesheet escape its <style> element, load a remote resource, or execute
// script -- regardless of how the input was spelled.
func FuzzCSS(f *testing.F) {

	for _, seed := range cssSeeds() {
		f.Add(seed)
	}

	validate := CSS("")

	f.Fuzz(func(t *testing.T, value string) {

		result, err := validate(value)

		// This format sanitizes rather than rejects, so it never returns an error.
		if err != nil {
			t.Fatalf("CSS returned an unexpected error for %q: %v", value, err)
		}

		requireSafeCSSOutput(t, value, result)

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

// FuzzCSS_WhitespaceIsInsignificant asserts that re-indenting a stylesheet never
// changes how much of it survives.  CSS collapses every run of whitespace
// between tokens, so a rule that is kept when written on one line and dropped
// when written across four is being destroyed by formatting alone -- which is
// the shape of the bluemonday trailing-separator bug that normalizeCSSDeclarations
// exists to absorb, and the shape any future one would take.
func FuzzCSS_WhitespaceIsInsignificant(f *testing.F) {

	for _, seed := range cssSeeds() {
		f.Add(seed)
	}

	validate := CSS("")

	f.Fuzz(func(t *testing.T, value string) {

		expected := strings.Count(mustSanitizeCSS(t, validate, value), "{")

		for _, whitespace := range []string{" ", "\t", "\n", "\r\n", "\f", " \n\t\t"} {

			variant := respaceCSS(value, whitespace)

			if actual := strings.Count(mustSanitizeCSS(t, validate, variant), "{"); actual != expected {
				t.Fatalf("re-indenting %q with %q changed the rule count from %d to %d", value, whitespace, expected, actual)
			}
		}
	})
}

// FuzzCSS_SafeRuleAlwaysSurvives builds a stylesheet from a selector, an
// allowlisted property, and a value, and asserts that a rule every check in this
// package accepts is actually emitted.
//
// The security targets all ask what the sanitizer REMOVES.  This one asks what it
// must KEEP, which is the half that had no coverage -- and the half a silent
// sanitizer fails without any test going red.
func FuzzCSS_SafeRuleAlwaysSurvives(f *testing.F) {

	f.Add("h1", uint8(0), "red", uint8(0))
	f.Add(".article > p:first-child", uint8(3), "1.6", uint8(1))
	f.Add("a[target]", uint8(7), "0 auto", uint8(2))
	f.Add("h1, h2", uint8(11), "#336699", uint8(3))
	f.Add("*", uint8(19), "12px/1.5 sans-serif", uint8(1))

	// The indentation styles a person might write the same rule in
	layouts := []string{
		"%s { %s: %s; }",
		"%s {\n\t%s: %s;\n}",
		"%s\n{\n\t%s : %s ;\n}\n",
		"%s {\r\n\t%s: %s\r\n}",
	}

	validate := CSS("")

	f.Fuzz(func(t *testing.T, selector string, propertyIndex uint8, value string, layoutIndex uint8) {

		property := cssAllowedProperties[int(propertyIndex)%len(cssAllowedProperties)]
		layout := layouts[int(layoutIndex)%len(layouts)]

		// RULE: the property is allowlisted by construction, so the rule is a
		// candidate only when the selector and the value pass their own checks too.
		// Bluemonday lowercases a value before the handler sees it.
		if !isSafeCSSPrelude(selector) || !isSafeCSSValue(strings.ToLower(value)) {
			return
		}

		// A prelude or value of nothing but whitespace carries no rule to keep
		if (strings.Trim(selector, cssWhitespace) == "") || (strings.Trim(value, cssWhitespace) == "") {
			return
		}

		// KNOWN GAP: the declaration round trip runs through gorilla/css, which
		// rejects a value whose parentheses do not balance and drops the whole
		// block with it.  That is upstream of this package and cannot be fixed by
		// normalizing, so those values are excluded rather than quietly tolerated.
		if strings.Count(value, "(") != strings.Count(value, ")") {
			return
		}

		stylesheet := fmt.Sprintf(layout, selector, property, value)
		result := mustSanitizeCSS(t, validate, stylesheet)

		if result == "" {
			t.Fatalf("selector %q, property %q, and value %q all passed their own checks, but %q sanitized away to nothing", selector, property, value, stylesheet)
		}

		requireSafeCSSOutput(t, stylesheet, result)
	})
}

// FuzzNormalizeCSSDeclarations asserts that normalizing a declaration list only
// ever rewrites whitespace and drops empty declarations.  A normalizer that
// introduces a character is a hole in every check that runs after it.
func FuzzNormalizeCSSDeclarations(f *testing.F) {

	for _, seed := range cssDeclarationSeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, declarations string) {

		result := normalizeCSSDeclarations(declarations)

		// RULE: normalizing rewrites whitespace and drops empty declarations, and
		// touches nothing else.  Reading both sides with the whitespace and the
		// separators taken out is what states that exactly: the content a browser
		// would apply has to come through unchanged.
		if withoutCSSNoise(result) != withoutCSSNoise(declarations) {
			t.Fatalf("normalizeCSSDeclarations(%q) = %q, which changed the declaration content", declarations, result)
		}

		// A separator is only ever dropped, never invented
		if strings.Count(result, ";") > strings.Count(declarations, ";") {
			t.Fatalf("normalizeCSSDeclarations(%q) = %q, which added a declaration separator", declarations, result)
		}

		// The only whitespace that may survive is a plain space
		if strings.ContainsAny(result, "\t\r\n\f") {
			t.Fatalf("normalizeCSSDeclarations(%q) = %q, which still carries non-space whitespace", declarations, result)
		}

		// RULE: the only character normalizing may introduce is the space it puts
		// in place of a whitespace run.  This is what makes it safe to run the
		// quote and bracket checks AFTER normalizing rather than before.
		for _, r := range result {
			if (r != ' ') && !strings.ContainsRune(declarations, r) {
				t.Fatalf("normalizeCSSDeclarations(%q) = %q, which introduced the character %q", declarations, result, r)
			}
		}

		// No declaration may be empty, because an empty one fails the CSS parse
		// downstream and takes the whole block with it
		for _, declaration := range strings.Split(result, ";") {
			if strings.Trim(declaration, cssWhitespace) == "" && result != "" {
				t.Fatalf("normalizeCSSDeclarations(%q) = %q, which contains an empty declaration", declarations, result)
			}
		}

		// Normalizing is a fixed point
		if again := normalizeCSSDeclarations(result); again != result {
			t.Fatalf("normalizeCSSDeclarations is not idempotent for %q: %q then %q", declarations, result, again)
		}
	})
}

// FuzzSanitizeDeclarations asserts that the declaration-block sanitizer emits
// only text that is safe to write between a selector's braces.
func FuzzSanitizeDeclarations(f *testing.F) {

	for _, seed := range cssDeclarationSeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, declarations string) {

		result := sanitizeDeclarations(declarations)

		if result == "" {
			return
		}

		// Nothing that survives may close the style element, open a block, or
		// carry a CSS escape
		if strings.ContainsAny(result, "\"'<>{}\\") {
			t.Fatalf("sanitizeDeclarations(%q) = %q, which carries a structural character", declarations, result)
		}

		// Every declaration block is written onto one line of the stylesheet
		if strings.ContainsAny(result, "\r\n\f") {
			t.Fatalf("sanitizeDeclarations(%q) = %q, which spans more than one line", declarations, result)
		}

		// A non-empty block is always terminated, so the caller can append to it
		if !strings.HasSuffix(result, ";") {
			t.Fatalf("sanitizeDeclarations(%q) = %q, which is not terminated", declarations, result)
		}

		// Sanitizing is a fixed point
		if again := sanitizeDeclarations(result); again != result {
			t.Fatalf("sanitizeDeclarations is not idempotent for %q: %q then %q", declarations, result, again)
		}
	})
}

// FuzzIsSafeCSSValue asserts that an accepted property value cannot carry
// structure of its own, in any spelling.
func FuzzIsSafeCSSValue(f *testing.F) {

	for _, seed := range cssValueSeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {

		if !isSafeCSSValue(value) {
			return
		}

		// A value is written verbatim after the colon, so it may not end its own
		// declaration, open a block, escape the style element, or open a comment
		if strings.ContainsAny(value, "\"'<>{};:@\\*&") {
			t.Fatalf("isSafeCSSValue accepted %q, which carries a structural character", value)
		}

		// Nor may it name a function that fetches a resource or runs script
		lower := strings.ToLower(value)

		for _, forbidden := range []string{"url", "expression", "javascript", "image", "cross-fade", "element", "attr"} {
			if strings.Contains(lower, forbidden) {
				t.Fatalf("isSafeCSSValue accepted %q, which names %q", value, forbidden)
			}
		}

		// Control characters have no meaning in a value and hide intent
		for _, r := range value {
			if (r < ' ') && (r != '\t') {
				t.Fatalf("isSafeCSSValue accepted %q, which carries the control character %q", value, r)
			}
		}
	})
}

// FuzzIsSafeCSSPrelude asserts that an accepted selector or at-rule condition
// cannot carry structure of its own.  Preludes are emitted verbatim, so this
// check is the whole of their sanitizing.
func FuzzIsSafeCSSPrelude(f *testing.F) {

	f.Add("")
	f.Add("h1")
	f.Add("#main .article > p:first-child")
	f.Add("h1,\nh2")
	f.Add("a[target]")
	f.Add("(max-width: 40em)")
	f.Add("screen and (max-width: 600px)")
	f.Add("</style><script>alert(1)</script>")
	f.Add("bo\\64 y")
	f.Add("a { color: red; } b")
	f.Add("@media")

	f.Fuzz(func(t *testing.T, prelude string) {

		if !isSafeCSSPrelude(prelude) {
			return
		}

		// `>` stays, because it is the child combinator and a `<` can never
		// accompany it.  Everything else that could open a tag, a block, a rule, a
		// comment, or a CSS escape is forbidden.
		if strings.ContainsAny(prelude, "\"'<{};@\\&/") {
			t.Fatalf("isSafeCSSPrelude accepted %q, which carries a structural character", prelude)
		}

		if strings.Contains(prelude, "}") {
			t.Fatalf("isSafeCSSPrelude accepted %q, which can close a block", prelude)
		}

		// The only whitespace a prelude may carry is whitespace CSS recognizes
		for _, r := range prelude {
			if (r < ' ') && !strings.ContainsRune(cssWhitespace, r) {
				t.Fatalf("isSafeCSSPrelude accepted %q, which carries the control character %q", prelude, r)
			}
		}
	})
}

// FuzzSplitCSSBlock asserts that reading one brace-balanced block off the front
// of a string loses nothing: the block and the remainder reassemble into exactly
// what was read.  A split that drops or duplicates text would let the structural
// walk skip past a rule, or read one twice.
func FuzzSplitCSSBlock(f *testing.F) {

	f.Add("{ color: red; }")
	f.Add("{ color: red; } h1 { color: blue; }")
	f.Add("{ a { b } } tail")
	f.Add("{")
	f.Add("}")
	f.Add("")
	f.Add("no leading brace")
	f.Add("{{{{{{")
	f.Add("{}}}}}")
	f.Add("{ \xff\xfe }")

	f.Fuzz(func(t *testing.T, value string) {

		body, rest, ok := splitCSSBlock(value)

		if !ok {
			// A rejected split reports nothing, so the caller cannot use a partial read
			if (body != "") || (rest != "") {
				t.Fatalf("splitCSSBlock(%q) failed but returned body %q and rest %q", value, body, rest)
			}
			return
		}

		if reassembled := "{" + body + "}" + rest; reassembled != value {
			t.Fatalf("splitCSSBlock(%q) reassembles to %q", value, reassembled)
		}

		// The walk consumes the block, so `rest` must be strictly shorter than the
		// input.  Anything else is an infinite loop in sanitizeStylesheet.
		if len(rest) >= len(value) {
			t.Fatalf("splitCSSBlock(%q) returned a remainder %q that did not advance", value, rest)
		}
	})
}

// FuzzSplitCSSAtRule asserts that an at-rule prelude decomposes into a name and a
// condition that both came from the prelude itself.
func FuzzSplitCSSAtRule(f *testing.F) {

	f.Add("@media screen")
	f.Add("@media (min-width: 40em)")
	f.Add("@media(min-width: 40em)")
	f.Add("@media")
	f.Add("@")
	f.Add("")
	f.Add("@media\n\tscreen")
	f.Add("@MEDIA screen and (max-width: 600px)")

	f.Fuzz(func(t *testing.T, prelude string) {

		name, condition := splitCSSAtRule(prelude)

		if !strings.HasPrefix(prelude, name) {
			t.Fatalf("splitCSSAtRule(%q) returned the name %q, which is not its prefix", prelude, name)
		}

		// A name is one token, so it can never carry whitespace or open the condition
		if strings.ContainsAny(name, cssWhitespace+"(") {
			t.Fatalf("splitCSSAtRule(%q) returned the name %q, which is more than one token", prelude, name)
		}

		if (condition != "") && !strings.Contains(prelude, condition) {
			t.Fatalf("splitCSSAtRule(%q) returned the condition %q, which is not part of it", prelude, condition)
		}
	})
}

// FuzzStripCSSComments asserts that comment removal takes every comment out and
// puts nothing back that a later pass would have to parse again.
func FuzzStripCSSComments(f *testing.F) {

	f.Add("")
	f.Add("body { color: red; }")
	f.Add("body { /* a note */ color: red; }")
	f.Add("body { color: re/**/d; }")
	f.Add("/* unterminated")
	f.Add("/*/")
	f.Add("/**/")
	f.Add("/* /* nested */ */")
	f.Add("*/")
	f.Add("/*\n*/body{color:red}")

	f.Fuzz(func(t *testing.T, stylesheet string) {

		result := stripCSSComments(stylesheet)

		// A comment opener that survived would let a later pass hide a brace
		if strings.Contains(result, "/*") {
			t.Fatalf("stripCSSComments(%q) = %q, which still carries a comment opener", stylesheet, result)
		}

		// Removing a comment can only shorten the stylesheet
		if len(result) > len(stylesheet) {
			t.Fatalf("stripCSSComments(%q) = %q, which is longer than its input", stylesheet, result)
		}

		// Stripping is a fixed point: a `*/` must never be spliced back together
		if again := stripCSSComments(result); again != result {
			t.Fatalf("stripCSSComments is not idempotent for %q: %q then %q", stylesheet, result, again)
		}
	})
}

// cssDeclarationSeeds returns the declaration lists shared by the targets that
// exercise one block at a time.
func cssDeclarationSeeds() []string {

	return []string{
		"",
		" ",
		";",
		";;;",
		"color: red",
		"color: red;",
		"color: red;\n",
		"\n\tcolor: red;\n",
		"\n\tline-height: 1.6;\n\tcolor: #333;\n",
		"box-shadow: 0 0 4px\n\t#000;",
		"color: red; behavior: url(x.htc); margin: 0",
		`font-family: "Helvetica Neue"`,
		"content: 'x'",
		"color: red</style><script>alert(1)</script>",
		"position: fixed",
		"position: relative",
		"color: \\65 xpression(alert(1))",
		"color: re/**/d",
		"color: \x00red",
		"margin: 0 auto; padding: 1em 2em",
		strings.Repeat("color: red;", 256),
	}
}

// cssValueSeeds returns the property values shared by the value-level targets.
func cssValueSeeds() []string {

	return []string{
		"",
		" ",
		"red",
		"RED",
		"#336699",
		"1.6",
		"0 auto",
		"12px/1.5 sans-serif",
		"rgba(0, 0, 0, .5)",
		"1em !important",
		"url(x)",
		"URL(x)",
		"expression(alert(1))",
		"javascript:alert(1)",
		"image-set(x)",
		"attr(href)",
		`"quoted"`,
		"red; behavior: url(x)",
		"red</style>",
		"\\65 xpression",
		"a\x00b",
		"\xff\xfe",
		"red\nblue",
	}
}

// withoutCSSNoise returns a declaration list with everything normalizing is
// allowed to change -- whitespace and declaration separators -- removed, leaving
// only the content that has to survive it untouched.
func withoutCSSNoise(declarations string) string {

	return strings.Map(func(r rune) rune {

		if (r == ';') || isCSSWhitespace(r) {
			return -1
		}

		return r
	}, declarations)
}

// mustSanitizeCSS runs the stylesheet sanitizer and fails the test if it reports
// an error, which it is not supposed to be able to do.
func mustSanitizeCSS(t *testing.T, validate StringFormat, value string) string {

	t.Helper()

	result, err := validate(value)

	if err != nil {
		t.Fatalf("CSS returned an unexpected error for %q: %v", value, err)
	}

	return result
}

// requireSafeCSSOutput fails the test if a sanitized stylesheet carries any
// construct that would let it escape its <style> element, load a remote
// resource, or execute script.
func requireSafeCSSOutput(t *testing.T, value string, result string) {

	t.Helper()

	lower := strings.ToLower(result)

	// No output may contain the characters that would end the <style> element or
	// open a construct the sanitizer did not itself write.  `>` is a legal child
	// combinator, so it is not banned outright -- `<` is what makes a closing tag
	// possible, and it is forbidden.
	for _, forbidden := range []string{"<", `"`, "'", `\`, "/*", "*/", "@import", "@charset", "@namespace", "@font-face"} {
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
}
