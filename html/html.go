package html

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var findHTML *regexp.Regexp
var spaces *regexp.Regexp
var breaks *regexp.Regexp
var paragraphs *regexp.Regexp
var divs *regexp.Regexp
var headings *regexp.Regexp
var styles *regexp.Regexp
var anchors *regexp.Regexp
var tags *regexp.Regexp

func init() {
	findHTML = regexp.MustCompile(`(?i)<[A-Z]+.*?>`)
	spaces = regexp.MustCompile(`[[:space:]]+`)
	breaks = regexp.MustCompile(`(?i)<br[^>]*>`)
	paragraphs = regexp.MustCompile(`(?i)<\/p>`)
	headings = regexp.MustCompile(`(?i)<\/?h[0-9][^>]*>`)
	divs = regexp.MustCompile(`(?i)<\/div>`)
	styles = regexp.MustCompile(`(?i)<style>(.*?)</style>`)
	anchors = regexp.MustCompile(`(?i)<a[^>]*>(.*?)</a>`)
	tags = regexp.MustCompile(`<[^>]+>`)
}

// IsHTML returns TRUE if the string provided "looks like" HTML, in that, it has
// one or more substrings that appear to be an HTML tag
func IsHTML(html string) bool {
	return findHTML.Match([]byte(html))
}

// ToText returns a string that has been converted from HTML into plain text.
// Mostly, this means replacing block level tags (BR, P, DIV) with carriage returns.
func ToText(html string) string {

	result := html

	// Replace HTML tags
	result = spaces.ReplaceAllLiteralString(result, " ")
	result = breaks.ReplaceAllLiteralString(result, "\n")
	result = paragraphs.ReplaceAllLiteralString(result, "\n\n")
	result = headings.ReplaceAllLiteralString(result, "\n\n")
	result = divs.ReplaceAllLiteralString(result, "\n")
	result = styles.ReplaceAllLiteralString(result, "")
	result = tags.ReplaceAllLiteralString(result, "")

	// Replace HTML entities
	result = RemoveSpecialCharacters(result)
	result = strings.Trim(result, " ")

	return result
}

// ToSearchText removes tags in a way that is suitable to text searches.
// This means that it will remove all tags, but adds regular whitespace in between them.
func ToSearchText(html string) string {
	result := html
	result = strings.Replace(result, "<", " <", -1) // Add space before every HTML tag
	result = RemoveTags(result)
	result = RemoveSpecialCharacters(result)
	return result
}

// RemoveSpecialCharacters removes special Unicode characters from a string
func RemoveSpecialCharacters(html string) string {

	result := html

	result = strings.Replace(result, "\u00a0", " ", -1)

	result = strings.Replace(result, "&#60;", "<", -1)
	result = strings.Replace(result, "&lt;", "<", -1)

	result = strings.Replace(result, "&#62;", ">", -1)
	result = strings.Replace(result, "&gt;", ">", -1)

	result = strings.Replace(result, "&#34;", `"`, -1)
	result = strings.Replace(result, "&quot;", `"`, -1)

	result = strings.Replace(result, "&#38;", "&", -1)
	result = strings.Replace(result, "&amp;", "&", -1)

	result = strings.Replace(result, "&#39;", "'", -1)
	result = strings.Replace(result, "&apos;", "'", -1)
	result = strings.Replace(result, "&apos;", "'", -1)
	result = strings.Replace(result, "&lsquo;", "'", -1)
	result = strings.Replace(result, "&rsquo;", "'", -1)

	result = strings.Replace(result, "&#124;", "|", -1)
	result = strings.Replace(result, "&#145;", "'", -1)
	result = strings.Replace(result, "&#146;", "'", -1)
	result = strings.Replace(result, "&#147;", `"`, -1)
	result = strings.Replace(result, "&#148;", `"`, -1)
	result = strings.Replace(result, "&ldquo;", `"`, -1)
	result = strings.Replace(result, "&rdquo;", `"`, -1)

	result = strings.Replace(result, "&ndash;", `-`, -1)
	result = strings.Replace(result, "&mdash;", `-`, -1)
	result = strings.Replace(result, "&#150;", `-`, -1)
	result = strings.Replace(result, "&#151;", `-`, -1)

	result = strings.Replace(result, "&#160;", " ", -1)
	result = strings.Replace(result, "&nbsp;", " ", -1)

	result = strings.Replace(result, "&#169;", "(C)", -1)
	result = strings.Replace(result, "&copy;", "(C)", -1)

	result = strings.Replace(result, "&#171;", "<<", -1)
	result = strings.Replace(result, "&laquo;", "<<", -1)

	result = strings.Replace(result, "&#187;", ">>", -1)
	result = strings.Replace(result, "&raquo;", ">>", -1)

	result = strings.Replace(result, "&#174;", "(R)", -1)
	result = strings.Replace(result, "&reg;", "(R)", -1)

	result = strings.Replace(result, "&#8230;", "...", -1)
	result = strings.Replace(result, "&hellip;", "...", -1)

	result = strings.Replace(result, "&#8249;", "<", -1)
	result = strings.Replace(result, "&lsaquo;", "<", -1)

	result = strings.Replace(result, "&#8250;", ">", -1)
	result = strings.Replace(result, "&rsaquo;", "<", -1)

	return result
}

// RemoveAnchors strips all HTML anchor tags from a string.
func RemoveAnchors(html string) string {
	return string(anchors.ReplaceAll([]byte(html), []byte("$1")))
}

// RemoveTags aggressively strips HTML tags from a string.  It will only keep anything between `>` and `<`.
// From: https://stackoverflow.com/questions/55036156/how-to-replace-all-html-tag-with-empty-string-in-golang
// Original code by: Daniel Morell <https://stackoverflow.com/users/10463261/daniel-morell>
func RemoveTags(html string) string {

	const (
		htmlTagStart = 60 // Unicode `<`
		htmlTagEnd   = 62 // Unicode `>`
	)

	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(html) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range html {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(html) && end >= start {
			builder.WriteString(html[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(html[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}

	return builder.String()
}
