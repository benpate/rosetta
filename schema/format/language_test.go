package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLanguage confirms that the language format rewrites language tags to their canonical form
func TestLanguage(t *testing.T) {

	validate := Language("")

	cases := map[string]string{
		"":      "", // an empty value means "not set"
		"en":    "en",
		"EN-us": "en-US",
		"iw":    "he", // deprecated codes are stored in their current form
		"pt-BR": "pt-BR",
	}

	for input, expected := range cases {
		result, err := validate(input)
		require.NoError(t, err, input)
		require.Equal(t, expected, result, input)
	}
}

// TestLanguage_Rejects confirms that the language format rejects anything that is not a language
func TestLanguage_Rejects(t *testing.T) {

	validate := Language("")

	for _, input := range []string{"not a tag", "und", "mul", "*", "en_US!"} {
		result, err := validate(input)
		require.Error(t, err, input)
		require.Empty(t, result, input)
	}
}
