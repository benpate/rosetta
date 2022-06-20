package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoHTML(t *testing.T) {
	value, err := NoHTML("")("<i>Is this HTML?</i>")

	require.Nil(t, err)
	require.Equal(t, "Is this HTML?", value)
}

func TestNoHTML_Text(t *testing.T) {
	value, err := NoHTML("")("This is plain text")

	require.Nil(t, err)
	require.Equal(t, "This is plain text", value)
}

func TestMalformedHTML1(t *testing.T) {
	value, err := NoHTML("")("Do something <strong>bold</strong>.")
	require.Nil(t, err)
	require.Equal(t, "Do something bold.", value)
}

func TestMalformedHTML2(t *testing.T) {
	value, err := NoHTML("")("h1>I broke this</h1>")
	require.Nil(t, err)
	require.Equal(t, "I broke this", value)
}

func TestMalformedHTML3(t *testing.T) {
	value, err := NoHTML("")("This is <a href='#'>>broken link</a>.")
	require.Nil(t, err)
	require.Equal(t, "This is broken link.", value)
}

func TestMalformedHTML4(t *testing.T) {
	value, err := NoHTML("")("I don't know ><where to <<em>start</em> this tag<.")
	require.Nil(t, err)
	require.Equal(t, "start this tag", value)
}

func TestMalformedHTML5(t *testing.T) {
	value, err := NoHTML("")(`<A HREF="http://example.com/comment.cgi?mycomment=<SCRIPT
	SRC='http://bad-site/badfile'></SCRIPT>"> Click here</A>`)
	require.Nil(t, err)
	require.Equal(t, "Click here", value)
}
