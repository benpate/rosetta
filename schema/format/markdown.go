package format

// Markdown accepts Markdown source, passing it through unchanged.
// It does not sanitize, so rendered output is safe only if the renderer sanitizes it.
func Markdown(arg string) StringFormat {

	// Without this format, an unregistered name falls back to NoHTML, which strips tags and
	// collapses the whitespace that Markdown needs.  Sanitizing here would corrupt the source
	// just as badly, because "a < b" in prose and fenced code samples that quote markup are
	// both legitimate Markdown.
	return func(value string) (string, error) {

		// Untouched, unbothered, unbowed
		return value, nil
	}
}
