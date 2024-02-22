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

func TestRemoveAnchors(t *testing.T) {
	require.Equal(t, "Regular string", RemoveAnchors("Regular string"))
	require.Equal(t, "Regular string <b>with tags</b>", RemoveAnchors("Regular string <b>with tags</b>"))
	require.Equal(t, `Regular string <span class="some-class">with tags</span> and <i>attributes</i>.`, RemoveAnchors(`Regular string <span class="some-class">with tags</span> and <i>attributes</i>.`))

	require.Equal(t, "A string with anchors and stuff", RemoveAnchors("A string with <a href='#'>anchors</a> and stuff"))
	require.Equal(t, `A string with anchors and stuff`, RemoveAnchors(`A string with <a href="#">anchors</a> and stuff`))
	require.Equal(t, `Begins with an anchor`, RemoveAnchors(`<a href="#">Begins with</a> an anchor`))
	require.Equal(t, `Begins with an anchor`, RemoveAnchors(`<a href="#">Begins with </a>an anchor`))
	require.Equal(t, `Ends with an anchor`, RemoveAnchors(`Ends with <a href="#">an anchor</a>`))
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
