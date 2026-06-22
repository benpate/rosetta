package schema

import (
	"sync"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema/format"
	"github.com/stretchr/testify/require"
)

func TestGetLength(t *testing.T) {
	length, ok := getLength(testArrayA{"one", "two"})
	require.True(t, ok)
	require.Equal(t, 2, length)
}

func TestGetLength_NotLengthGetter(t *testing.T) {
	_, ok := getLength("not-an-array")
	require.False(t, ok)
}

func TestGetIndex(t *testing.T) {
	item, ok := getIndex(testArrayA{"one", "two"}, 0)
	require.True(t, ok)
	require.Equal(t, "one", item)
}

func TestGetIndex_NotArrayGetter(t *testing.T) {
	_, ok := getIndex("not-an-array", 0)
	require.False(t, ok)
}

func TestIsMultipleOfInteger(t *testing.T) {
	require.True(t, isMultipleOfInteger(10, 5))
	require.False(t, isMultipleOfInteger(10, 3))

	// A zero divisor is treated as "no constraint" and must not panic.
	require.True(t, isMultipleOfInteger(10, 0))

	// Large int64 values must not be corrupted by a detour through float64.
	// 2^53 + 1 cannot be represented exactly as a float64.
	const big = int64(9_007_199_254_740_993) // 2^53 + 1
	require.True(t, isMultipleOfInteger(big, 1))
	require.False(t, isMultipleOfInteger(big, 2))

	// Unsigned integers are supported too.
	require.True(t, isMultipleOfInteger(uint64(12), uint64(4)))
	require.False(t, isMultipleOfInteger(uint64(12), uint64(5)))
}

func TestIsMultipleOfFloat(t *testing.T) {
	require.True(t, isMultipleOfFloat(10.0, 2.5))
	require.False(t, isMultipleOfFloat(10.0, 3.0))

	// A zero divisor is treated as "no constraint".
	require.True(t, isMultipleOfFloat(10.0, 0.0))
}

func TestType_String(t *testing.T) {
	require.Equal(t, "string", TypeString.String())
}

func TestUseFormat_NilIsIgnored(t *testing.T) {
	// Registering a nil generator is a no-op and must not panic
	require.NotPanics(t, func() { UseFormat("coverage-nil-format", nil) })
	require.NotContains(t, formats, "coverage-nil-format")
}

func TestUseFormat_RegistersGenerator(t *testing.T) {
	// Reset the freeze latch to simulate startup-time registration, since
	// other tests may have already frozen the registry by reading from it.
	formatsFrozen.Store(false)

	UseFormat("coverage-test-format", format.NoHTML)
	require.Contains(t, formats, "coverage-test-format")
	delete(formats, "coverage-test-format")
}

func TestUseFormat_IgnoredAfterFreeze(t *testing.T) {
	// Once the registry has been read (frozen), a late UseFormat is a no-op
	// and must not mutate the map (it is reported, not applied).
	formatsFrozen.Store(true)

	UseFormat("coverage-late-format", format.NoHTML)
	require.NotContains(t, formats, "coverage-late-format")
}

func TestUseFormat_ConcurrentReadsAndLateWrites(t *testing.T) {
	// After freezing, concurrent validations (readers) and late UseFormat
	// calls (no-op writers) must not race. Run under `go test -race`.
	formatsFrozen.Store(true)

	// Silence derp reporting so the many expected "ignored" reports don't
	// flood the test output.
	savedPlugins := derp.Plugins
	derp.Plugins = nil
	defer func() { derp.Plugins = savedPlugins }()

	element := String{Format: "no-html"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); element.formatFunctions() }()
		go func() { defer wg.Done(); UseFormat("coverage-race-format", format.NoHTML) }()
	}
	wg.Wait()

	require.NotContains(t, formats, "coverage-race-format")
}
