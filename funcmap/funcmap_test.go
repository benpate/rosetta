package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctions_DollarFormat(t *testing.T) {

	f := All()

	dollarFormat := f["dollarFormat"].(func(any) string)

	require.Equal(t, "$12.34", dollarFormat(1234))
}
