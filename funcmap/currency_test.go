package funcmap

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyFuncs(t *testing.T) {

	dollarFormat := All()["dollarFormat"].(func(any) string)

	require.Equal(t, "$12.34", dollarFormat(1234))
	require.Equal(t, "$0.05", dollarFormat(5))   // padding with leading zeros
	require.Equal(t, "$1.00", dollarFormat(1.0)) // float64 multiplied by 100
	require.Equal(t, "$2.50", dollarFormat(float32(2.5)))
}

func TestCurrencyFuncs_Negative(t *testing.T) {

	dollarFormat := All()["dollarFormat"].(func(any) string)

	// The sign belongs outside the currency symbol, and must not be counted as a digit
	require.Equal(t, "-$0.05", dollarFormat(-5))
	require.Equal(t, "-$0.50", dollarFormat(-50))
	require.Equal(t, "-$1.50", dollarFormat(-150))
	require.Equal(t, "-$12.34", dollarFormat(-1234))
	require.Equal(t, "-$1.50", dollarFormat(-1.5))
	require.Equal(t, "-$2.50", dollarFormat(float32(-2.5)))

	// The most negative int64 must not overflow while its sign is removed
	require.Equal(t, "-$92233720368547758.08", dollarFormat(int64(math.MinInt64)))
}

func TestCurrencyFuncs_EdgeCases(t *testing.T) {

	dollarFormat := All()["dollarFormat"].(func(any) string)

	require.Equal(t, "$0.00", dollarFormat(0))
	require.Equal(t, "$0.01", dollarFormat(1))
	require.Equal(t, "-$0.01", dollarFormat(-1))
	require.Equal(t, "$92233720368547758.07", dollarFormat(int64(math.MaxInt64)))

	// Unconvertible values fall back to zero rather than panicking
	require.Equal(t, "$0.00", dollarFormat("not-a-number"))
	require.Equal(t, "$0.00", dollarFormat(nil))
}
