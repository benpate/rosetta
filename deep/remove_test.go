package deep

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestRemove_MapOfString proves that names are removed from a mapof.String, at any depth
func TestRemove_MapOfString(t *testing.T) {

	document := map[string]any{
		"direct": mapof.String{"bcc": "x", "to": "y"},
		"listed": sliceof.MapOfString{{"bcc": "x", "to": "y"}},
	}

	require.NoError(t, Remove(document, "bcc"))

	require.Equal(t, mapof.String{"to": "y"}, document["direct"])
	require.Equal(t, sliceof.MapOfString{{"to": "y"}}, document["listed"])
}
