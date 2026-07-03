package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

// TestValidate_RewriteDetails_FlatProperty confirms that when Validate rejects a value
// because a property had to be rewritten, the error details name the property and show
// its before/after values.
func TestValidate_RewriteDetails_FlatProperty(t *testing.T) {

	s := New(pointerGetterStruct_Schema())

	// "TooLongName" exceeds MaxLength:8, so validation truncates (rewrites) it.
	value := pointerGetterStruct{Name: "TooLongName", URL: "https://example.com/x"}

	_, _, err := Validate(s, &value)
	require.NotNil(t, err)

	details := errorDetails(t, err)
	require.Contains(t, details, `name: "TooLongName" -> "TooLongN"`)
}

// TestValidate_RewriteDetails_ArrayItem confirms that rewrites inside array items are
// reported with their full dot-path, including the array index.
func TestValidate_RewriteDetails_ArrayItem(t *testing.T) {

	s := New(Array{Items: pointerGetterStruct_Schema()})

	value := pointerGetterSlice{
		{Name: "Alice", URL: "https://example.com/alice"},
		{Name: "TooLongName", URL: "https://example.com/x"},
	}

	_, _, err := Validate(s, &value)
	require.NotNil(t, err)

	details := errorDetails(t, err)
	require.Contains(t, details, `1.name: "TooLongName" -> "TooLongN"`)
	require.NotContains(t, details, `0.name`) // the clean item is not reported
}

// errorDetails flattens a derp.Error's details into one string for Contains assertions.
func errorDetails(t *testing.T, err error) string {

	derpError, ok := err.(derp.Error)
	require.True(t, ok, "error must be a derp.Error")

	return fmt.Sprint(derpError.Details)
}

// TestRewriteList_Prefix covers path accumulation as rewrites bubble up through
// nested objects and arrays.
func TestRewriteList_Prefix(t *testing.T) {

	list := rewriteList{{From: "a", To: "b"}}

	// A leaf record has no path until its parent names it...
	list = list.prefix("profileUrl")
	require.Equal(t, "profileUrl", list[0].Path)

	// ...and each enclosing container prepends its own segment.
	list = list.prefix("0")
	list = list.prefix("links")
	require.Equal(t, "links.0.profileUrl", list[0].Path)
}

// TestRewrite_String_Abbreviates confirms that very long values are truncated in
// the human-readable report.
func TestRewrite_String_Abbreviates(t *testing.T) {

	long := strings.Repeat("x", 100)
	item := rewrite{Path: "statusMessage", From: long, To: "short"}

	report := item.String()
	require.Contains(t, report, "statusMessage: ")
	require.Contains(t, report, "…")
	require.Less(t, len(report), 120)
}
