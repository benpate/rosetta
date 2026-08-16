package format

import (
	"bytes"
	"sync"

	"github.com/benpate/derp"
	"github.com/yuin/goldmark"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// safeMarkdown renders Markdown the careful way: raw HTML is omitted, and dangerous
// URL schemes are blanked.  Its output is safe without any further sanitizing.
var safeMarkdown = sync.OnceValue(func() goldmark.Markdown {
	return goldmark.New()
})

// unsafeMarkdown renders Markdown permissively, passing raw HTML and URLs through
// verbatim.  Used only to compare against safeMarkdown -- never to produce output.
var unsafeMarkdown = sync.OnceValue(func() goldmark.Markdown {
	return goldmark.New(goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()))
})

// Markdown accepts Markdown source that renders safely under any renderer, and
// rejects source carrying raw HTML or a dangerous URL.
func Markdown(arg string) StringFormat {

	return func(value string) (string, error) {

		const location = "schema.format.Markdown"

		// Render the same source both carefully and permissively.  Cost is two renders,
		// linear in length -- pair this format with a MaxLength, which the schema checks
		// before it reaches any format function.
		safe, err := renderMarkdown(safeMarkdown(), value)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to render Markdown")
		}

		unsafe, err := renderMarkdown(unsafeMarkdown(), value)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to render Markdown")
		}

		// RULE: The two renderings must agree.  Where they differ, the source contains raw
		// HTML or a dangerous URL that only the permissive renderer would emit -- so the value
		// is unsafe in front of any renderer that does not sanitize.  Comparing renderings
		// needs no allow-list of its own, so this can never disagree with the sanitizing
		// policy that an application applies later.
		if safe != unsafe {
			return "", derp.Validation("Markdown must not contain HTML tags or unsafe links")
		}

		return value, nil
	}
}

// renderMarkdown converts Markdown source into HTML using the provided renderer.
func renderMarkdown(renderer goldmark.Markdown, value string) (string, error) {

	const location = "schema.format.renderMarkdown"

	var buffer bytes.Buffer

	if err := renderer.Convert([]byte(value), &buffer); err != nil {
		return "", derp.Wrap(err, location, "Error converting Markdown to HTML")
	}

	return buffer.String(), nil
}
