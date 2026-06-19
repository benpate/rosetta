package compare

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsZero_Numerics(t *testing.T) {
	require.True(t, IsZero(int8(0)))
	require.True(t, IsZero(int16(0)))
	require.True(t, IsZero(int32(0)))
	require.True(t, IsZero(int64(0)))
	require.True(t, IsZero(uint8(0)))
	require.True(t, IsZero(uint16(0)))
	require.True(t, IsZero(uint32(0)))
	require.True(t, IsZero(uint64(0)))
	require.True(t, IsZero(float32(0)))
	require.True(t, IsZero(float64(0)))
	require.True(t, IsZero(false))

	require.False(t, IsZero(int64(1)))
	require.False(t, IsZero(float64(0.1)))
	require.False(t, IsZero(true))
}

func TestIsZero_Slices(t *testing.T) {
	require.True(t, IsZero([]string{}))
	require.True(t, IsZero([]int{}))
	require.True(t, IsZero([]float64{}))
	require.True(t, IsZero([]any{}))

	require.False(t, IsZero([]string{"a"}))
	require.False(t, IsZero([]int{1}))
	require.False(t, IsZero([]float64{1}))
	require.False(t, IsZero([]any{1}))
}

// --- Custom types implementing the comparison interfaces ---

type zeroStringer string

func (z zeroStringer) String() string { return string(z) }

type zeroInter int

func (z zeroInter) Int() int { return int(z) }

type zeroFloater float64

func (z zeroFloater) Float() float64 { return float64(z) }

type zeroNuller bool

func (z zeroNuller) IsNull() bool { return bool(z) }

type zeroLength int

func (z zeroLength) Length() int { return int(z) }

func TestIsZero_Interfaces(t *testing.T) {

	require.True(t, IsZero(zeroStringer("")))
	require.False(t, IsZero(zeroStringer("hello")))

	require.True(t, IsZero(zeroInter(0)))
	require.False(t, IsZero(zeroInter(5)))

	require.True(t, IsZero(zeroFloater(0)))
	require.False(t, IsZero(zeroFloater(1.5)))

	require.True(t, IsZero(zeroNuller(true)))
	require.False(t, IsZero(zeroNuller(false)))

	require.True(t, IsZero(zeroLength(0)))
	require.False(t, IsZero(zeroLength(3)))
}

type zeroHexer string

func (z zeroHexer) Hex() string { return string(z) }

func TestIsZero_Hexer(t *testing.T) {
	// "0" parses to int64(0), which is zero
	require.True(t, IsZero(zeroHexer("0")))
	// "ff" parses to int64(255), which is not zero
	require.False(t, IsZero(zeroHexer("ff")))
}

func TestIsZero_Reader(t *testing.T) {
	require.True(t, IsZero(strings.NewReader("")))
	require.False(t, IsZero(strings.NewReader("content")))
}

func TestNotZero(t *testing.T) {
	require.True(t, NotZero(1))
	require.False(t, NotZero(0))
}
