package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMarkdown confirms that Markdown source survives formatting byte-for-byte.
func TestMarkdown(t *testing.T) {

	markdown := Markdown("")

	// Each of these is destroyed by NoHTML, Text, or HTML, which is why this format exists
	values := []string{
		"",
		"# Title",
		"a < b && c > d",
		"Line one\nLine two",
		"Para one\n\nPara two",
		"- one\n  - nested\n    - deeper",
		"    indented code block",
		"```go\nif a < b {\n\tfmt.Println(\"<hi>\")\n}\n```",
		"Inline `<script>` in a code span",
		"| A | B |\n|---|---|\n| 1 | 2 |",
		"Text with  multiple   spaces",
		"Trailing hard break  \nnext line",
		"Unicode: café 🎉 日本語",
	}

	for _, value := range values {
		result, err := markdown(value)
		require.NoError(t, err, "value: %q", value)
		require.Equal(t, value, result, "value: %q", value)
	}
}

// TestMarkdown_PassesThroughRawHTML confirms that raw HTML in the source is neither
// stripped nor escaped.
func TestMarkdown_PassesThroughRawHTML(t *testing.T) {

	markdown := Markdown("")

	// Stripping here would corrupt code samples that quote markup, so the render
	// boundary -- not this format -- is what makes the output safe.
	value := "<script>alert(1)</script>"
	result, err := markdown(value)

	require.NoError(t, err)
	require.Equal(t, value, result)
}

// TestMarkdown_IgnoresArgument confirms that the format takes no configuration today,
// so an argument neither errors nor changes the result.
func TestMarkdown_IgnoresArgument(t *testing.T) {

	value := "# Title"

	for _, arg := range []string{"", "anything", "strict"} {
		result, err := Markdown(arg)(value)
		require.NoError(t, err, "arg: %q", arg)
		require.Equal(t, value, result, "arg: %q", arg)
	}
}

// TestMarkdown_LargeInput confirms that a large document is returned intact.
func TestMarkdown_LargeInput(t *testing.T) {

	value := strings.Repeat("# Heading\n\nSome **text** here.\n\n", 1000)

	result, err := Markdown("")(value)

	require.NoError(t, err)
	require.Equal(t, value, result)
}

// FuzzMarkdown confirms that the format never errors and never alters its input.
func FuzzMarkdown(f *testing.F) {

	f.Add("")
	f.Add("# Title")
	f.Add("a < b")
	f.Add("<script>alert(1)</script>")
	f.Add("\xff\xfe")
	f.Add("    indented\n\n\ttabbed")

	markdown := Markdown("")

	f.Fuzz(func(t *testing.T, value string) {
		result, err := markdown(value)
		require.NoError(t, err)
		require.Equal(t, value, result)
	})
}
