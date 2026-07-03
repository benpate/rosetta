package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// nestedContent is a struct property that is reached via PointerGetter and does NOT
// implement ValueSetter -- exactly like Emissary's model.Content.  Its "html" field uses
// a rewriting format, so validating it produces a rewrite that the parent must write back.
type nestedContent struct {
	Format string
	HTML   string
}

func (c *nestedContent) GetPointer(name string) (any, bool) {
	switch name {
	case "format":
		return &c.Format, true
	case "html":
		return &c.HTML, true
	}
	return nil, false
}

func (c *nestedContent) GetStringOK(name string) (string, bool) {
	switch name {
	case "format":
		return c.Format, true
	case "html":
		return c.HTML, true
	}
	return "", false
}

// nestedOuter holds a nestedContent by value and exposes it via a pointer.
type nestedOuter struct {
	Content nestedContent
}

func (o *nestedOuter) GetPointer(name string) (any, bool) {
	if name == "content" {
		return &o.Content, true
	}
	return nil, false
}

func nestedOuterSchema() Schema {
	return New(Object{Properties: ElementMap{
		"content": Object{Properties: ElementMap{
			"format": String{},
			"html":   String{Format: "no-html"}, // rewrites: strips HTML tags
		}},
	}})
}

// TestNormalize_NestedObject_Rewrite_WritesBack reproduces the Emissary bandwagon-news
// "Cannot set values on empty path" crash: a nested struct property (no ValueSetter) whose
// child field is rewritten during Normalize.  The parent must write the rewritten child
// back through its pointer -- without a ValueSetter -- via reflection.
func TestNormalize_NestedObject_Rewrite_WritesBack(t *testing.T) {

	s := nestedOuterSchema()

	outer := &nestedOuter{
		Content: nestedContent{
			Format: "HTML",
			HTML:   "Content <b>content</b> content", // the <b> tags will be stripped
		},
	}

	paths, err := s.Normalize(outer)

	// Before the fix this failed with "Cannot set values on empty path".
	require.NoError(t, err)

	// The html field was rewritten in place (tags stripped)...
	require.Equal(t, "Content content content", outer.Content.HTML)

	// ...and the rewrite was reported against the correct dotted path.
	require.Contains(t, paths, "content.html")

	// The sibling field is untouched.
	require.Equal(t, "HTML", outer.Content.Format)
}

// TestNormalize_NestedObject_NoRewrite_Clean confirms the happy path: when nothing needs
// rewriting, Normalize succeeds and reports no changed paths.
func TestNormalize_NestedObject_NoRewrite_Clean(t *testing.T) {

	s := nestedOuterSchema()

	outer := &nestedOuter{
		Content: nestedContent{Format: "HTML", HTML: "Clean content"},
	}

	paths, err := s.Normalize(outer)

	require.NoError(t, err)
	require.Empty(t, paths)
	require.Equal(t, "Clean content", outer.Content.HTML)
}
