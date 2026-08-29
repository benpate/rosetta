package format

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// cssMaxNestingDepth bounds how deeply conditional at-rules (@media, @supports)
// may nest before the remainder of the stylesheet is discarded. A stylesheet is
// a recursive format, so it needs an explicit depth limit: without one, deeply
// nested input drives unbounded recursion here.
const cssMaxNestingDepth = 4

// cssWhitespace lists the characters CSS treats as whitespace between tokens.
const cssWhitespace = " \t\r\n\f"

// cssPolicy sanitizes the declarations inside a single CSS block. Bluemonday is
// an HTML sanitizer, so it reaches CSS only through an element's `style`
// attribute; sanitizeDeclarations wraps a declaration list in a placeholder
// element to run it through this policy. Every property that survives must be
// named here, and every value must satisfy isSafeCSSValue.
var cssPolicy *bluemonday.Policy

// cssAllowedProperties is the allowlist of CSS properties that may appear in a
// sanitized stylesheet. It covers typography, color, spacing, borders, and the
// flex/grid layout primitives -- enough to restyle a page, and nothing that
// loads an external resource or escapes its own stacking context.
//
// Deliberately absent: `behavior` and `-moz-binding` (script execution in
// legacy engines), `content` (can inject text, and needs quoted strings that
// isSafeCSSValue rejects), and `filter` (the legacy IE spelling took
// `progid:...` expressions).
var cssAllowedProperties = []string{

	// Text and fonts
	"color", "direction", "font", "font-family", "font-size", "font-stretch",
	"font-style", "font-variant", "font-weight", "letter-spacing", "line-height",
	"list-style", "list-style-position", "list-style-type", "quotes", "tab-size",
	"text-align", "text-decoration", "text-decoration-color", "text-decoration-line",
	"text-decoration-style", "text-indent", "text-overflow", "text-shadow",
	"text-transform", "vertical-align", "white-space", "word-break",
	"word-spacing", "overflow-wrap", "hyphens",

	// Box model
	"box-sizing", "height", "margin", "margin-bottom", "margin-left",
	"margin-right", "margin-top", "max-height", "max-width", "min-height",
	"min-width", "padding", "padding-bottom", "padding-left", "padding-right",
	"padding-top", "width",

	// Borders and surfaces
	"background-color", "border", "border-bottom", "border-bottom-color",
	"border-bottom-left-radius", "border-bottom-right-radius", "border-bottom-style",
	"border-bottom-width", "border-collapse", "border-color", "border-left",
	"border-left-color", "border-left-style", "border-left-width", "border-radius",
	"border-right", "border-right-color", "border-right-style", "border-right-width",
	"border-spacing", "border-style", "border-top", "border-top-color",
	"border-top-left-radius", "border-top-right-radius", "border-top-style",
	"border-top-width", "border-width", "box-shadow", "opacity", "outline",
	"outline-color", "outline-offset", "outline-style", "outline-width",

	// Layout
	"align-content", "align-items", "align-self", "clear", "column-count",
	"column-gap", "columns", "display", "flex", "flex-basis", "flex-direction",
	"flex-flow", "flex-grow", "flex-shrink", "flex-wrap", "float", "gap",
	"grid-auto-columns", "grid-auto-flow", "grid-auto-rows", "grid-column",
	"grid-column-end", "grid-column-start", "grid-row", "grid-row-end",
	"grid-row-start", "grid-template-areas", "grid-template-columns",
	"grid-template-rows", "justify-content", "justify-items", "justify-self",
	"order", "overflow", "overflow-x", "overflow-y", "row-gap", "visibility",

	// Tables
	"caption-side", "empty-cells", "table-layout",
}

// cssAllowedAtRules is the allowlist of block at-rules whose bodies are
// recursively sanitized and kept. Both are purely conditional: they wrap
// ordinary style rules without loading anything.
//
// Every other at-rule is dropped, which is what removes `@import` (fetches a
// remote stylesheet, leaking the viewer's IP and importing unsanitized rules),
// `@font-face` (fetches a remote font via `src`), `@charset`, and `@namespace`.
var cssAllowedAtRules = []string{"@media", "@supports"}

// init builds the declaration-sanitizing policy. Each allowed property is bound
// to the same value handler, so the property allowlist decides *what* may be
// styled and isSafeCSSValue decides *how* -- rejecting any value that could
// fetch a resource, execute script, or break out of its own declaration.
func init() {
	cssPolicy = bluemonday.NewPolicy()
	cssPolicy.AllowElements(cssPlaceholderElement)
	cssPolicy.AllowAttrs("style").OnElements(cssPlaceholderElement)
	cssPolicy.AllowStyles(cssAllowedProperties...).MatchingHandler(isSafeCSSValue).Globally()

	// RULE: `position` is allowlisted by value, not just by name. A stylesheet is
	// injected into a page it does not own, so `fixed`/`absolute`/`sticky` would
	// let it lift an element out of the flow and cover unrelated page chrome --
	// the setup for a clickjacking overlay. In-flow positioning stays available.
	cssPolicy.AllowStyles("position").MatchingEnum("static", "relative").Globally()
}

// CSS sanitizes a CSS stylesheet, keeping only the rules that are safe to embed
// in a page the stylesheet's author does not otherwise control. Selectors are
// checked against a character allowlist, declarations are filtered by
// bluemonday against a property allowlist and a value allowlist, and everything
// unrecognized -- comments, unknown at-rules, unbalanced blocks, trailing
// garbage -- is discarded rather than repaired.
//
// Like HTML, this format sanitizes rather than rejects: it returns the cleaned
// stylesheet with a nil error, so a value that is partly unsafe is stored
// stripped instead of failing the whole write.
//
// The result is a stylesheet, not a fragment, and it is still global in scope --
// nothing here confines its selectors to one subtree. A caller embedding it in a
// shared page is responsible for scoping the rules it emits.
func CSS(arg string) StringFormat {
	return func(value string) (string, error) {
		return sanitizeStylesheet(value, 0), nil
	}
}

// sanitizeStylesheet walks a list of CSS rules, emitting only those that pass
// validation. It is called once for the whole stylesheet and again for the body
// of each conditional at-rule, with depth counting the nesting so recursion
// stays bounded.
func sanitizeStylesheet(stylesheet string, depth int) string {

	var result strings.Builder

	for remaining := stripCSSComments(stylesheet); ; {

		remaining = strings.TrimLeft(remaining, cssWhitespace)

		if remaining == "" {
			break
		}

		blockIndex := strings.IndexByte(remaining, '{')

		// No block remains, so whatever is left is a truncated rule or a
		// bodyless at-rule. Neither is kept.
		if blockIndex < 0 {
			break
		}

		// A statement that terminates before the block opens is a bodyless
		// at-rule (`@import url(...);`). None are allowed, so skip past it.
		if statementIndex := strings.IndexByte(remaining[:blockIndex], ';'); statementIndex >= 0 {
			remaining = remaining[statementIndex+1:]
			continue
		}

		prelude := strings.TrimSpace(remaining[:blockIndex])
		body, rest, ok := splitCSSBlock(remaining[blockIndex:])

		// An unbalanced block means the rest of the input cannot be parsed
		// structurally. Discard it rather than guessing where it ended.
		if !ok {
			break
		}

		remaining = rest

		if strings.HasPrefix(prelude, "@") {
			writeCSSAtRule(&result, prelude, body, depth)
			continue
		}

		writeCSSStyleRule(&result, prelude, body)
	}

	return result.String()
}

// writeCSSAtRule validates a conditional at-rule and, if it is allowed,
// recursively sanitizes its body and writes the whole rule to result.
func writeCSSAtRule(result *strings.Builder, prelude string, body string, depth int) {

	// RULE: Nesting is capped so a deeply nested stylesheet cannot drive
	// unbounded recursion. At the cap the rule is dropped whole.
	if depth >= cssMaxNestingDepth {
		return
	}

	name, condition := splitCSSAtRule(prelude)

	if !containsFold(cssAllowedAtRules, name) {
		return
	}

	// Both allowed at-rules are conditional, so one carrying no condition is
	// malformed no matter how it is spelled.
	if condition == "" {
		return
	}

	// The condition is read by the CSS parser rather than a selector engine, but
	// it is emitted into the same text context and gets the same check. The name
	// itself needs no check: it matched the allowlist above.
	if !isSafeCSSPrelude(condition) {
		return
	}

	inner := sanitizeStylesheet(body, depth+1)

	// An at-rule whose body sanitized away entirely has nothing left to apply.
	if inner == "" {
		return
	}

	// The rule is rebuilt from its validated parts rather than echoed, so the
	// output spelling cannot differ from what was actually checked.
	result.WriteString(strings.ToLower(name))
	result.WriteString(" ")
	result.WriteString(condition)
	result.WriteString(" {\n")
	result.WriteString(inner)
	result.WriteString("}\n")
}

// writeCSSStyleRule validates a selector and its declarations, writing the rule
// to result only when both survive sanitizing.
func writeCSSStyleRule(result *strings.Builder, selector string, body string) {

	if !isSafeCSSPrelude(selector) {
		return
	}

	declarations := sanitizeDeclarations(body)

	if declarations == "" {
		return
	}

	result.WriteString(selector)
	result.WriteString(" { ")
	result.WriteString(declarations)
	result.WriteString(" }\n")
}

// cssPlaceholderElement is the element name sanitizeDeclarations wraps a
// declaration list in. It exists only to give bluemonday a `style` attribute to
// sanitize and never appears in any output.
const cssPlaceholderElement = "x"

// sanitizeDeclarations filters one block of CSS declarations through the
// bluemonday policy, returning only the declarations whose property is
// allowlisted and whose value passes isSafeCSSValue.
//
// Bluemonday cannot sanitize bare CSS, so the declarations make the round trip
// as an HTML `style` attribute: they are wrapped in a placeholder element,
// sanitized, and read back out. Bluemonday HTML-escapes what it emits, so the
// recovered value is unescaped before being returned to CSS context.
func sanitizeDeclarations(declarations string) string {

	// A quote would let a declaration close the attribute early and change what
	// bluemonday parses, so no value containing one may enter the round trip.
	// isSafeCSSValue rejects quotes too; this is the check that protects the
	// wrapping itself, before the policy ever runs.
	if strings.ContainsAny(declarations, `"'<>`) {
		return ""
	}

	sanitized := cssPolicy.Sanitize(`<` + cssPlaceholderElement + ` style="` + declarations + `"></` + cssPlaceholderElement + `>`)

	// Recover the attribute value from the sanitized element. A policy that
	// stripped every declaration emits the element with no attribute at all.
	_, after, found := strings.Cut(sanitized, `style="`)

	if !found {
		return ""
	}

	value, _, found := strings.Cut(after, `"`)

	if !found {
		return ""
	}

	value = strings.TrimSpace(html.UnescapeString(value))

	if value == "" {
		return ""
	}

	// RULE: Re-check the recovered text. Unescaping can only shorten the value or
	// reintroduce characters that were escaped on the way out, so the result is
	// verified in the form it will actually be written to the stylesheet.
	if strings.ContainsAny(value, `"'<>{}\`) {
		return ""
	}

	return strings.TrimSuffix(value, ";") + ";"
}

// isSafeCSSValue reports whether a single CSS property value is safe to emit.
// Bluemonday calls it for every declaration whose property is allowlisted, with
// the value already lowercased and stripped of unicode escapes.
//
// Validation is two-stage, matching the approach used for single property
// values elsewhere: a conservative character allowlist, then a check for the
// function names whose spellings fall entirely inside that allowlist.
func isSafeCSSValue(value string) bool {

	if value == "" {
		return false
	}

	if strings.IndexFunc(value, isUnsafeCSSValueRune) >= 0 {
		return false
	}

	// The character allowlist admits `(` and `)`, so it cannot by itself block
	// the CSS functions that fetch a resource or execute script. Reject those by
	// name: `url` loads a remote (or `javascript:`) resource, `expression` runs
	// script in legacy IE, and `image`/`image-set`/`cross-fade` are alternate
	// spellings of a fetch.
	lower := strings.ToLower(value)

	for _, unsafeName := range []string{"url", "expression", "javascript", "image", "cross-fade", "element", "attr"} {
		if strings.Contains(lower, unsafeName) {
			return false
		}
	}

	return true
}

// isUnsafeCSSValueRune reports whether a rune is outside the allowlist of
// characters permitted in a CSS property value. The set admits identifiers,
// numbers, colors, units, function calls, and `!important`.
//
// Notably absent are `"` `'` `;` `{` `}` `<` `>` `@` `\` `*` `&` and control
// characters -- respectively the building blocks of attribute and string
// breakouts, extra declarations, block breakouts, a `</style>` escape, new
// at-rules, CSS escape sequences, and comment openers. `/` is allowed (the
// `font: 12px/1.5` and `grid-area: 1 / 2` separators) because without `*` it
// cannot open a comment.
func isUnsafeCSSValueRune(r rune) bool {

	switch {

	case r >= 'a' && r <= 'z':
		return false

	case r >= 'A' && r <= 'Z':
		return false

	case r >= '0' && r <= '9':
		return false

	case strings.ContainsRune(" \t,.()#%+-/_!", r):
		return false
	}

	return true
}

// isSafeCSSPrelude reports whether a selector or at-rule condition contains only
// allowlisted characters. Preludes are emitted verbatim, so the check is what
// keeps them from carrying structure of their own.
//
// The set covers type, class, id, attribute, pseudo, and grouped selectors along
// with every combinator. Absent are `"` `'` `<` `/` `\` `{` `}` `;` `@` and `&`,
// which between them cover a `</style>` escape, comment openers, CSS escapes,
// and any attempt to open a block or rule the parser did not account for. `>`
// stays (the child combinator) because a `<` can never accompany it.
func isSafeCSSPrelude(prelude string) bool {

	if prelude == "" {
		return false
	}

	for _, r := range prelude {

		switch {

		case r >= 'a' && r <= 'z':
			continue

		case r >= 'A' && r <= 'Z':
			continue

		case r >= '0' && r <= '9':
			continue

		case strings.ContainsRune(" \t\r\n,.#:()[]=+~>*_-", r):
			continue
		}

		return false
	}

	return true
}

// splitCSSAtRule splits an at-rule prelude into its name (including the leading
// `@`) and its condition. The two are separated by whitespace or by the opening
// parenthesis of the condition, so `@media screen`, `@media (min-width: 40em)`,
// and `@media(min-width: 40em)` all split the same way.
func splitCSSAtRule(prelude string) (name string, condition string) {

	index := strings.IndexAny(prelude, cssWhitespace+"(")

	if index < 0 {
		return prelude, ""
	}

	return prelude[:index], strings.TrimSpace(prelude[index:])
}

// splitCSSBlock reads one brace-balanced block from the front of value, which
// must begin with `{`. It returns the block's contents and whatever follows it.
// A block that never closes returns ok false, which tells the caller to discard
// the remainder rather than guess where it ended.
func splitCSSBlock(value string) (body string, rest string, ok bool) {

	if !strings.HasPrefix(value, "{") {
		return "", "", false
	}

	depth := 0

	for index, r := range value {

		switch r {

		case '{':
			depth++

		case '}':
			depth--

			if depth == 0 {
				return value[1:index], value[index+1:], true
			}
		}
	}

	return "", "", false
}

// stripCSSComments removes every `/* ... */` comment from a stylesheet.
//
// Comments are removed before parsing rather than skipped during it, so a
// comment can never hide a brace or semicolon from the structural walk, and a
// `*/` can never be spliced back together by a later step. An unterminated
// comment consumes the rest of the input, which is what a CSS parser does too.
func stripCSSComments(stylesheet string) string {

	if !strings.Contains(stylesheet, "/*") {
		return stylesheet
	}

	var result strings.Builder

	for remaining := stylesheet; ; {

		before, after, found := strings.Cut(remaining, "/*")

		result.WriteString(before)

		if !found {
			break
		}

		_, afterComment, closed := strings.Cut(after, "*/")

		if !closed {
			break
		}

		// A comment separates tokens, so it leaves a space behind rather than
		// joining the text on either side into one identifier.
		result.WriteString(" ")
		remaining = afterComment
	}

	return result.String()
}

// containsFold reports whether values contains target, comparing without regard
// to case. At-rule names are case-insensitive in CSS.
func containsFold(values []string, target string) bool {

	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}

	return false
}
