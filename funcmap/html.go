package funcmap

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/url"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func addHTMLFuncs(target map[string]any) {

	target["addQueryParams"] = func(extraParams string, url string) string {
		if strings.Contains(url, "?") {
			return url + "&" + extraParams
		}
		return url + "?" + extraParams
	}

	// attr marks a string as a trusted HTML attribute value, which html/template
	// emits with NO escaping. To keep that safe by construction it first confirms
	// the value cannot break out of the attribute or inject a new one: any value
	// containing a quote, angle bracket, `=`, `&`, backtick, or whitespace is
	// rejected and returns empty. Legitimate attribute values (class/id tokens,
	// data values, ARIA states) pass through unchanged.
	target["attr"] = func(value string) template.HTMLAttr {
		return safeAttr(value)
	}

	// css marks a string as trusted CSS, so html/template emits it with NO
	// escaping. It is UNSAFE by construction: only pass owner- or admin-scoped
	// values (e.g. a profile's custom stylesheet). Never pass remote, federated,
	// or query-string data — use `cssValue` for a single computed property value.
	target["css"] = func(value string) template.CSS {
		return template.CSS(value)
	}

	// cssValue marks a single CSS property value as trusted, but ONLY after
	// confirming it contains no context-breakout characters. It is the safe
	// choice when the value is computed from data (colors, sizes, gradients)
	// rather than authored by the owner. Anything containing a character outside
	// the conservative property-value allowlist is rejected and returns empty,
	// so a crafted value cannot escape the `style` attribute or the property it
	// sits in. The allowlist covers identifiers, numbers, colors (#rrggbb),
	// units (%), functions (rgba(), linear-gradient()), and their separators.
	target["cssValue"] = func(value string) template.CSS {
		return safeCSSValue(value)
	}

	target["highlight"] = func(text string, search string) template.HTML {

		// Both `text` and `search` are untrusted (display names, query params), so they are
		// HTML-escaped before anything is marked safe. Escaping happens FIRST so that the only
		// live markup in the result is the <b> wrapper this function adds. Matching is done on
		// the already-escaped strings so needle and haystack share the same encoding.
		text = template.HTMLEscapeString(text)

		if search == "" {
			return template.HTML(text)
		}

		search = template.HTMLEscapeString(search)
		result := strings.ReplaceAll(text, search, `<b class="highlight">`+search+"</b>")
		return template.HTML(result)
	}

	target["domainOnly"] = func(value string) string {
		result := strings.TrimPrefix(value, "http://")
		result = strings.TrimPrefix(result, "https://")
		result, _, _ = strings.Cut(result, "/")
		return result
	}

	target["hasImage"] = func(value string) bool {
		if strings.Contains(value, "<img") {
			return true
		}

		if strings.Contains(value, "<picture") {
			return true
		}

		return false
	}

	target["html"] = func(value string) template.HTML {
		return template.HTML(value)
	}

	target["htmlMinimal"] = func(value string) template.HTML {
		return template.HTML(html.Minimal(value))
	}

	target["js"] = func(value string) string {
		return template.JSEscapeString(value)
	}

	target["json"] = func(value any) string {
		result, err := json.Marshal(value)

		if err != nil {
			derp.Report(derp.WrapIF(err, "tools.templates.functions.json", "Error marshaling JSON"))
		}

		return string(result)
	}

	target["jsonIndent"] = func(value any) string {
		result, err := json.MarshalIndent(value, "", "    ")
		if err != nil {
			derp.Report(derp.WrapIF(err, "tools.templates.functions.jsonIndent", "Error marshaling JSON with indentation"))
		}
		return string(result)
	}

	target["markdown"] = func(value any) template.HTML {

		valueBytes := convert.Bytes(value)

		// https://github.com/yuin/goldmark#built-in-extensions
		var buffer bytes.Buffer

		md := goldmark.New(
			goldmark.WithExtensions(
				extension.Table,
				extension.Linkify,
				extension.Typographer,
				extension.DefinitionList,
			),
			goldmark.WithRendererOptions(),
		)

		if err := md.Convert(valueBytes, &buffer); err != nil {
			derp.Report(derp.Wrap(err, "tools.templates.functions.markdown", "Error converting Markdown to HTML"))
		}

		return template.HTML(buffer.String())
	}

	target["queryEscape"] = func(value string) string {
		return url.QueryEscape(value)
	}

	target["summary"] = html.Summary

	target["stripProtocol"] = func(value string) string {
		if after, found := strings.CutPrefix(value, "http://"); found {
			return after
		}
		if after, found := strings.CutPrefix(value, "https://"); found {
			return after
		}
		return value
	}

	target["textOnly"] = html.RemoveTags

	target["text"] = func(value string) template.HTML {
		return template.HTML(html.FromText(value))
	}
}

// safeAttr validates an HTML attribute value and returns it as a trusted
// template.HTMLAttr only if it cannot break out of the attribute or inject a new
// one; otherwise it returns empty. Because template.HTMLAttr is emitted without
// escaping, any value carrying a quote, angle bracket, `=`, `&`, backtick, or
// whitespace — the characters needed to close the attribute, open a tag, or start
// an adjacent attribute like `onmouseover=` — is rejected wholesale.
func safeAttr(value string) template.HTMLAttr {

	if strings.IndexFunc(value, isUnsafeAttrRune) >= 0 {
		return template.HTMLAttr("")
	}

	return template.HTMLAttr(value)
}

// isUnsafeAttrRune reports whether a rune must not appear in a raw (unescaped)
// HTML attribute value. It rejects the attribute-breakout characters directly,
// plus all whitespace and control runes (which could split the value into an
// additional attribute). Everything else — the characters that make up class and
// id tokens, data values, and ARIA states — is permitted.
func isUnsafeAttrRune(r rune) bool {

	switch r {

	case '"', '\'', '`', '<', '>', '=', '&':
		return true
	}

	return r <= ' ' || r == 0x7f
}

// safeCSSValue validates a single CSS property value and returns it as trusted
// CSS only if it cannot break out of the value, the property, or the enclosing
// `style` attribute; otherwise it returns empty. Callers use it for values
// computed from data (colors, sizes, gradients) rather than authored by a
// trusted owner. Validation is two-stage: a conservative character allowlist,
// then a name check for the two dangerous functions (`expression`, `url`) whose
// spellings happen to fall entirely within that allowlist.
func safeCSSValue(value string) template.CSS {

	if strings.IndexFunc(value, isUnsafeCSSValueRune) >= 0 {
		return template.CSS("")
	}

	// The charset allowlist admits `()`, so it cannot block the two dangerous
	// CSS functions spelled with allowlisted characters: `expression(...)`
	// (legacy-IE script execution) and `url(...)` (can load remote or
	// `javascript:` resources). Reject them by name.
	lower := strings.ToLower(value)
	if strings.Contains(lower, "expression") || strings.Contains(lower, "url") {
		return template.CSS("")
	}

	return template.CSS(value)
}

// isUnsafeCSSValueRune reports whether a rune is outside the conservative
// allowlist of characters permitted in a `cssValue` property value. The allowed
// set is intentionally narrow: it admits identifiers, numbers, colors, units,
// simple functions, and their separators, and nothing that could terminate the
// value, the property, or the surrounding `style` attribute. Notably absent are
// `"` `'` `;` `<` `>` `{` `}` `@` `\` `/` `*` and control characters — the
// building blocks of quote/attribute breakouts, extra declarations, `@import`,
// comments, and CSS escapes.
func isUnsafeCSSValueRune(r rune) bool {

	switch {

	case r >= 'a' && r <= 'z':
		return false

	case r >= 'A' && r <= 'Z':
		return false

	case r >= '0' && r <= '9':
		return false

	case strings.ContainsRune(" \t,.()#%+-", r):
		return false
	}

	return true
}
