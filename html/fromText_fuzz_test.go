package html

import (
	"strings"
	"testing"
)

// FuzzFromText feeds arbitrary text to FromText and asserts the core safety invariant: the output
// is HTML-safe. FromText legitimately introduces <br> tags (from newlines), so we strip those
// first; whatever remains must contain NO raw "<" or ">" and every "&" must begin a known entity.
// A failure here means plain text could break out of HTML text context (an injection bug).
func FuzzFromText(f *testing.F) {

	f.Add("")
	f.Add("hello world")
	f.Add("a & b")
	f.Add("<script>alert(1)</script>")
	f.Add(`quotes " and ' here`)
	f.Add("&lt;already escaped&gt;")
	f.Add("line one\nline two")
	f.Add("&amp;&amp;&amp;")
	f.Add("<<<>>>&&&")
	f.Add("tabs\tand\r\nnewlines")

	f.Fuzz(func(t *testing.T, text string) {

		output := FromText(text)

		// Remove the only tag FromText is allowed to introduce.
		stripped := strings.ReplaceAll(output, "<br>", "")

		// INVARIANT 1: no raw angle brackets survive outside of the <br> tags.
		if strings.ContainsAny(stripped, "<>") {
			t.Fatalf("FromText(%q) = %q contains a raw angle bracket after removing <br>", text, output)
		}

		// INVARIANT 2: every "&" begins one of the entities FromText emits.
		assertWellFormedEntities(t, text, stripped)
	})
}

// FuzzFromTextRoundTrip asserts that FromText is lossless: decoding the entities it produces and
// turning <br> back into a space reproduces the same whitespace-collapsed text as collapsing the
// input directly. This proves FromText only escapes/formats content, never drops or corrupts it.
func FuzzFromTextRoundTrip(f *testing.F) {

	f.Add("hello world")
	f.Add("a & b < c > d")
	f.Add(`he said "hi" & left`)
	f.Add("it's a test")
	f.Add("line\nbreak")
	f.Add("&lt;not a tag&gt;")
	f.Add("\n")              // leading/trailing newline -> <br> at the edge (whitespace-collapse edge case)
	f.Add("&amp;lt;")        // entity-looking input must not decode back to "<"
	f.Add("  spaced  out  ") // runs of whitespace collapse to a single space

	f.Fuzz(func(t *testing.T, text string) {

		output := FromText(text)

		// Reverse FromText's transformations. <br> came from "\n" (whitespace), so decode it back to
		// "\n"; decode entities back to literals. "&amp;" is decoded LAST so it does not interfere
		// with the other entities (e.g. "&amp;lt;" must become "&lt;", not "<").
		decoded := strings.ReplaceAll(output, "<br>", "\n")
		decoded = strings.ReplaceAll(decoded, "&lt;", "<")
		decoded = strings.ReplaceAll(decoded, "&gt;", ">")
		decoded = strings.ReplaceAll(decoded, "&quot;", `"`)
		decoded = strings.ReplaceAll(decoded, "&#39;", "'")
		decoded = strings.ReplaceAll(decoded, "&amp;", "&")

		// FromText collapses whitespace as its final step and "<br>" is opaque to that collapsing,
		// so the only faithful comparison is to collapse whitespace on BOTH the round-tripped output
		// and the original input. Equality then proves no NON-whitespace content was lost or altered.
		if CollapseWhitespace(decoded) != CollapseWhitespace(text) {
			t.Fatalf("FromText round-trip mismatch:\n  input:   %q\n  output:  %q\n  decoded: %q", text, output, decoded)
		}
	})
}

// assertWellFormedEntities fails if any "&" in s does not begin one of the entities that FromText
// produces (&amp; &lt; &gt; &quot; &#39;). This proves no "&" was left unescaped.
func assertWellFormedEntities(t *testing.T, input string, s string) {
	t.Helper()

	known := []string{"&amp;", "&lt;", "&gt;", "&quot;", "&#39;"}

	for i := 0; i < len(s); i++ {
		if s[i] != '&' {
			continue
		}

		matched := false
		for _, entity := range known {
			if strings.HasPrefix(s[i:], entity) {
				matched = true
				break
			}
		}

		if !matched {
			t.Fatalf("FromText(%q) produced an unescaped '&' at offset %d in %q", input, i, s)
		}
	}
}
