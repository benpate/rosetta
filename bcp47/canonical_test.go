package bcp47

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCanonical confirms the one stored spelling for each usable language tag
func TestCanonical(t *testing.T) {

	cases := map[string]string{
		"en":                 "en",
		"EN-us":              "en-US",
		"en_US":              "en-US",
		"zh-hant-tw":         "zh-Hant-TW",
		"pt-BR":              "pt-BR",
		"iw":                 "he",          // deprecated code for Hebrew
		"i-klingon":          "tlh",         // grandfathered tag
		"ca-valencia":        "ca-valencia", // variants are a different language form, so they stay
		"en-US-u-ca-gregory": "en-US",       // extensions are formatting preferences
		"en-US-x-twain":      "en-US",       // private use is dropped
	}

	for input, expected := range cases {
		actual, ok := Canonical(input)
		require.True(t, ok, input)
		require.Equal(t, expected, actual, input)
	}
}

// TestCanonical_Rejects confirms that values naming no language are refused
func TestCanonical_Rejects(t *testing.T) {

	cases := []string{
		"",
		"*",
		"und",
		"mul",
		"x-foo", // private use only
		"xx",    // well-formed, but not a language
		"not a tag",
		"<script>",
		"en-" + strings.Repeat("a", maxTagLength),
	}

	for _, input := range cases {
		actual, ok := Canonical(input)
		require.False(t, ok, input)
		require.Empty(t, actual, input)
	}
}

// FuzzCanonical confirms that Canonical never panics, and that its output is its own canonical form
func FuzzCanonical(f *testing.F) {

	seeds := []string{"en", "EN-us", "zh-Hant-TW", "en-US-u-ca-gregory", "x-foo", "*", ""}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {

		result, ok := Canonical(input)

		if !ok {
			return
		}

		again, ok := Canonical(result)
		require.True(t, ok, input)
		require.Equal(t, result, again, input)
	})
}
