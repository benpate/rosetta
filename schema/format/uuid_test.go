package format

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

// TestUUID covers empty values, canonical and uppercase UUIDs, pasted URLs, and rejections
func TestUUID(t *testing.T) {

	const canonical = "f4a31f0a-51dd-4fa7-986d-3095c40c5ed9"

	valid := func(input string, expected string) {
		t.Helper()
		result, err := UUID("")(input)
		require.NoError(t, err, input)
		require.Equal(t, expected, result, input)
	}

	invalid := func(input string) {
		t.Helper()
		result, err := UUID("")(input)
		require.Error(t, err, input)
		require.True(t, derp.IsValidationError(err), input)
		require.Equal(t, "", result, input)
	}

	// Empty values are allowed
	valid("", "")
	valid("   ", "")

	// Canonical UUIDs are kept, and others are made canonical
	valid(canonical, canonical)
	valid("F4A31F0A-51DD-4FA7-986D-3095C40C5ED9", canonical)
	valid("  "+canonical+"\n", canonical)

	// A pasted URL keeps only its UUID
	valid("https://musicbrainz.org/release/"+canonical, canonical)
	valid("https://musicbrainz.org/recording/"+canonical+"/", canonical)
	valid("https://musicbrainz.org/artist/"+canonical+"?tab=releases#top", canonical)
	valid("musicbrainz.org/release-group/"+canonical, canonical)

	// Anything else is rejected
	invalid("not-a-uuid")
	invalid("f4a31f0a51dd4fa7986d3095c40c5ed9")
	invalid("{" + canonical + "}")
	invalid("urn:uuid:" + canonical)
	invalid(canonical + "0")
	invalid("g4a31f0a-51dd-4fa7-986d-3095c40c5ed9")
	invalid("https://musicbrainz.org/release/" + canonical + "/discids")
	invalid("https://musicbrainz.org/")
	invalid("/")
}

// TestLastPathSegment covers URLs, bare values, and the edges of the path
func TestLastPathSegment(t *testing.T) {
	require.Equal(t, "abc", lastPathSegment("abc"))
	require.Equal(t, "abc", lastPathSegment("https://example.com/x/abc"))
	require.Equal(t, "abc", lastPathSegment("https://example.com/x/abc//"))
	require.Equal(t, "abc", lastPathSegment("https://example.com/abc?q=1#f"))
	require.Equal(t, "", lastPathSegment("/"))
	require.Equal(t, "", lastPathSegment("?abc"))
}
