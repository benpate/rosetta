package format

import (
	"strings"
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

// TestMarkdown_AcceptsPureMarkdown confirms that ordinary Markdown passes through
// untouched.  These constructs render identically whether or not raw HTML is honored,
// so none of them may be mistaken for injected markup.
func TestMarkdown_AcceptsPureMarkdown(t *testing.T) {

	markdown := Markdown("")

	values := []string{
		"",
		"# Title",
		"Hello **world** and *emphasis*",
		"[link](https://example.com)",
		"[relative](/about)",
		"[mail](mailto:someone@example.com)",
		"![img](https://example.com/a.png)",
		"a < b && c > d",
		"- one\n  - nested",
		"    indented code block",
		"```go\nif a < b {\n\tfmt.Println(\"<hi>\")\n}\n```",
		"Inline `<script>` in a code span",
		"line one  \nline two",
		"> quoted text",
		"---",
		"| A | B |\n|---|---|\n| 1 | 2 |",
		"Text with  multiple   spaces",
		"Unicode: café 🎉 日本語",

		// A slash directly after the tag name is not valid CommonMark tag syntax, so this
		// is escaped to inert text by every renderer rather than parsed as an element.
		// It is accepted because honoring raw HTML does not change the output.
		"<svg/onload=alert(1)>",
	}

	for _, value := range values {
		result, err := markdown(value)
		require.NoError(t, err, "value: %q", value)
		require.Equal(t, value, result, "value: %q", value)
	}
}

// TestMarkdown_RejectsActiveContent confirms that source whose rendering depends on raw
// HTML or a dangerous URL being honored is refused.
func TestMarkdown_RejectsActiveContent(t *testing.T) {

	markdown := Markdown("")

	values := []string{
		"<script>alert(1)</script>",
		"Hello <script src='https://evil.example/x.js'></script>",
		"<img src=x onerror=alert(1)>",
		"<div onclick=\"alert(1)\">click</div>",
		"<iframe src=\"https://evil.example\"></iframe>",
		"<svg onload=alert(1)>",
		"<svg onload='alert(1)'></svg>",
		"[x](javascript:alert(1))",
		"[x](vbscript:alert(1))",
		"# Title\n\nprose, then <script>alert(1)</script>",
	}

	for _, value := range values {
		result, err := markdown(value)
		require.Error(t, err, "value: %q", value)
		require.Equal(t, "", result, "value: %q", value)
	}
}

// TestMarkdown_RejectsBenignHTML confirms that ALL raw HTML is rejected, not only
// dangerous markup.  Allowing "safe" tags would require a policy, and the point of this
// format is that it needs none -- callers who want raw HTML opt out with "unsafe-any"
// and sanitize at the render boundary instead.
func TestMarkdown_RejectsBenignHTML(t *testing.T) {

	markdown := Markdown("")

	for _, value := range []string{"<b>bold</b>", "Text with <em>inline</em> html", "<br>"} {
		_, err := markdown(value)
		require.Error(t, err, "value: %q", value)
	}
}

// TestMarkdown_ReportsValidationError confirms that a rejection is reported as a
// validation error, so callers surface it as bad input rather than a server fault.
func TestMarkdown_ReportsValidationError(t *testing.T) {

	_, err := Markdown("")("<script>alert(1)</script>")

	require.Error(t, err)
	require.True(t, derp.IsValidationError(err))
}

// TestMarkdown_IgnoresArgument confirms that the format takes no configuration, so an
// argument neither errors nor changes the result.
func TestMarkdown_IgnoresArgument(t *testing.T) {

	for _, arg := range []string{"", "anything", "strict"} {
		result, err := Markdown(arg)("# Title")
		require.NoError(t, err, "arg: %q", arg)
		require.Equal(t, "# Title", result, "arg: %q", arg)
	}
}

// TestMarkdown_LargeInput confirms that a large document is handled without error.
func TestMarkdown_LargeInput(t *testing.T) {

	value := strings.Repeat("# Heading\n\nSome **text** here.\n\n", 1000)

	result, err := Markdown("")(value)

	require.NoError(t, err)
	require.Equal(t, value, result)
}

// FuzzMarkdown confirms that validation never panics, and that whatever it accepts
// renders identically with and without raw HTML honored.
func FuzzMarkdown(f *testing.F) {

	f.Add("")
	f.Add("# Title")
	f.Add("a < b")
	f.Add("<script>alert(1)</script>")
	f.Add("[x](javascript:alert(1))")
	f.Add("```\n<b>x</b>\n```")
	f.Add("\xff\xfe")

	markdown := Markdown("")

	f.Fuzz(func(t *testing.T, value string) {

		result, err := markdown(value)

		if err != nil {
			require.Equal(t, "", result)
			return
		}

		// An accepted value is returned unchanged, and is safe under any renderer
		require.Equal(t, value, result)

		safe, safeErr := renderMarkdown(safeMarkdown(), value)
		unsafe, unsafeErr := renderMarkdown(unsafeMarkdown(), value)

		require.NoError(t, safeErr)
		require.NoError(t, unsafeErr)
		require.Equal(t, safe, unsafe)
	})
}
