package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTML(t *testing.T) {
	value, err := HTML("")("<i>is this HTML?</i>")

	require.Nil(t, err)
	require.Equal(t, "<i>is this HTML?</i>", value)
}
