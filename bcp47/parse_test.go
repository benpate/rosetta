package bcp47

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// TestParse confirms that Parse returns the same tag that Canonical spells
func TestParse(t *testing.T) {

	tag, ok := Parse("pt-br-u-ca-gregory")
	require.True(t, ok)
	require.Equal(t, language.MustParse("pt-BR"), tag)

	tag, ok = Parse("und")
	require.False(t, ok)
	require.Equal(t, language.Und, tag)
}
