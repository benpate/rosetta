package html

import (
	"strings"
)

// FromText converts plain text into (lightly) formatted HTML
func FromText(text string) string {

	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, `"`, "&quot;")
	text = strings.ReplaceAll(text, "\n", "<br>")
	text = CollapseWhitespace(text)

	return text
}
