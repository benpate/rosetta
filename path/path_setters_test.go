package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {

	d := getTestData()

	// Not implemented yet, so this should just error
	require.NotNil(t, Set(d, "anywhere.doesnt.matter", "1"))
}
