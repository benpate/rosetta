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

func TestFromText_Ampersand(t *testing.T) {

	// A literal ampersand must be escaped to &amp;
	require.Equal(t, "Tom &amp; Jerry", FromText("Tom & Jerry"))

	// "&" is escaped FIRST, so an entity in the input is escaped literally rather than double-escaped.
	require.Equal(t, "&amp;lt; is a less-than entity", FromText("&lt; is a less-than entity"))

	// All escapable characters together, in order.
	require.Equal(t, "a &amp; b &lt; c &gt; d &quot; e &#39; f", FromText(`a & b < c > d " e ' f`))
}

func TestFromText_SingleQuote(t *testing.T) {
	require.Equal(t, "it&#39;s mine", FromText("it's mine"))
}

// TestFromText_EdgeCases pins exact output for boundary and tricky inputs.
func TestFromText_EdgeCases(t *testing.T) {

	cases := map[string]struct {
		input  string
		expect string
	}{
		"empty":                  {"", ""},
		"plain":                  {"hello", "hello"},
		"only newline":           {"\n", "<br>"},
		"leading/trailing space": {"  hi  ", "hi"},
		"collapsed whitespace":   {"a \t b", "a b"},
		"newline then text":      {"\nhi", "<br>hi"},
		"ampersand only":         {"&", "&amp;"},
		"entity in input":        {"&amp;", "&amp;amp;"},
		"all specials":           {`& < > " '`, "&amp; &lt; &gt; &quot; &#39;"},
		"script tag":             {"<script>x</script>", "&lt;script&gt;x&lt;/script&gt;"},
		"attribute injection":    {`" onload="evil`, "&quot; onload=&quot;evil"},
	}

	for name, c := range cases {
		require.Equal(t, c.expect, FromText(c.input), name)
	}
}
