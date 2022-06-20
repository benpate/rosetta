package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromText(t *testing.T) {

	original1 := `This is a less than "<".
		This is a greater than ">"
		This is a tag that will get displayed <evil tag="true"/><evil>`

	expected1 := `This is a less than &quot;&lt;&quot;.<br> This is a greater than &quot;&gt;&quot;<br> This is a tag that will get displayed &lt;evil tag=&quot;true&quot;/&gt;&lt;evil&gt;`

	require.Equal(t, expected1, FromText(original1))
}

func TestFromTextBR(t *testing.T) {
	original := `This is a test
	of the break tag.`

	expected := "This is a test<br> of the break tag."

	require.Equal(t, expected, FromText(original))
}
