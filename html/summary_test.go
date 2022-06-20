package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSummary_PlainText(t *testing.T) {
	content := "This is the first paragraph of the content.  There is NO HTML."
	require.Equal(t, "This is the first paragraph of the content. There is NO HTML.", Summary(content))
}

func TestSummary1(t *testing.T) {
	content := "<p>This is the first, but really, really long paragraph of the content.  I'm testing to see that we can still truncate to 200 characters even though the first paragraph says otherwise.  I wonder what else I'll have to say in order to make this truncate correctly?</p><p>Here's the second paragraph</p>"
	require.Equal(t, "This is the first, but really, really long paragraph of the content. I'm testing to see that we can still truncate to 200 characters even though the first paragraph says otherwise. I wonder what else ...", Summary(content))
}

func TestSummary2(t *testing.T) {
	content := "<p>This is the first, but really, really long paragraph of the content.  I'm testing to see that we can still truncate to 200 characters even though the first paragraph says otherwise.  I wonder what else I'll have to say in order to make this truncate correctly?</p><p>Here's the second paragraph</p>"
	require.Equal(t, "This is the first, but really, really long paragraph of the content. I'm testing to see that we can still truncate to 200 characters even though the first paragraph says otherwise. I wonder what else ...", Summary(content))
}
