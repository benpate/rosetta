package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveTags(t *testing.T) {
	require.Equal(t, "Regular string", RemoveTags("Regular string"))
	require.Equal(t, "Regular string with tags", RemoveTags("Regular string <b>with tags</b>"))
	require.Equal(t, "Regular string with tags and attributes.", RemoveTags(`Regular string <span class="some-class">with tags</span> and <i>attributes</i>.`))

	require.Equal(t, "Do something bold.", RemoveTags("Do something <strong>bold</strong>."))
	require.Equal(t, "I broke this", RemoveTags("h1>I broke this</h1>"))
	require.Equal(t, "This is broken link.", RemoveTags("This is <a href='#'>>broken link</a>."))
	require.Equal(t, "start this tag", RemoveTags("I don't know ><where to <<em>start</em> this tag<."))
}

func TestIsHTML(t *testing.T) {
	require.Equal(t, true, IsHTML("<b>This is HTML</b>"))
	require.Equal(t, false, IsHTML("This is not HTML"))
}

func TestToText(t *testing.T) {
	require.Equal(t, "Hello Gordon", ToText("<i>Hello</i> <b>Gordon</b>"))
	require.Equal(t, "Without Stylesheets", ToText("<style>* {font-weight:bold;}</style> Without Stylesheets"))
	require.Equal(t, "...", ToText("&hellip;"))
}
