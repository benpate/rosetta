package schema

import (
	"fmt"

	"github.com/benpate/rosetta/convert"
)

// rewrite records a single value that was modified during validation: the path
// where it lives, and its before/after values.
type rewrite struct {
	Path string // Dot-separated path to the rewritten value ("" means the root value)
	From any    // The value as it was provided
	To   any    // The value after formatting/clamping/truncation
}

// rewriteList collects every value that was modified during a validation pass.
// An empty list means the value already conforms to the schema.
type rewriteList []rewrite

// newRewriteList returns a single-item rewriteList when a leaf value has been
// changed, or nil when it has not.  Parent objects and arrays fill in the Path
// as the change bubbles up through them.
func newRewriteList(changed bool, from any, to any) rewriteList {

	if !changed {
		return nil
	}

	return rewriteList{{From: from, To: to}}
}

// prefix prepends a path segment (a property name or array index) to every
// rewrite in the list, and returns the list for chaining.
func (list rewriteList) prefix(segment string) rewriteList {

	for index := range list {
		if list[index].Path == "" {
			list[index].Path = segment
		} else {
			list[index].Path = segment + "." + list[index].Path
		}
	}

	return list
}

// paths returns the dot-separated path of every rewrite in the list, so that
// callers can report which values were modified.
func (list rewriteList) paths() []string {

	result := make([]string, len(list))

	for index, item := range list {
		result[index] = item.Path
	}

	return result
}

// details formats every rewrite in the list for use as derp error details, so
// that validation errors name each rewritten property with its before/after values.
func (list rewriteList) details() []any {

	result := make([]any, len(list))

	for index, item := range list {
		result[index] = item.String()
	}

	return result
}

// String returns a human-readable description of this rewrite, in the form:
// `path: "before" -> "after"`
func (item rewrite) String() string {

	path := item.Path

	if path == "" {
		path = "(root value)"
	}

	return fmt.Sprintf("%s: %q -> %q", path, abbreviate(item.From), abbreviate(item.To))
}

// abbreviate renders a value as a string, truncating it so that very long
// values (looking at you, 2048-character status messages) don't flood the logs.
func abbreviate(value any) string {

	const maxLength = 64

	result := convert.String(value)
	runes := []rune(result)

	if len(runes) <= maxLength {
		return result
	}

	return string(runes[:maxLength]) + "…"
}
