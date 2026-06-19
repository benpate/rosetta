package funcmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyFuncs(t *testing.T) {

	dollarFormat := All()["dollarFormat"].(func(any) string)

	require.Equal(t, "$12.34", dollarFormat(1234))
	require.Equal(t, "$0.05", dollarFormat(5))       // padding with leading zeros
	require.Equal(t, "$1.00", dollarFormat(1.0))      // float64 multiplied by 100
	require.Equal(t, "$2.50", dollarFormat(float32(2.5)))
}
