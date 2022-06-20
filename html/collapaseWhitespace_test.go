package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollapseWhitespace(t *testing.T) {

	content1 := `
		This is the first paragraph of the content.  
		This is the second paragraph of the content.
	`

	require.Equal(t, "This is the first paragraph of the content. This is the second paragraph of the content.", CollapseWhitespace(content1))
}
