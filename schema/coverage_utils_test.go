package schema

import (
	"sync"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema/format"
	"github.com/stretchr/testify/require"
)

func TestIndirect_Pointer(t *testing.T) {
	value := "hello"
	require.Equal(t, "hello", indirect(&value))
}

func TestIndirect_NotPointer(t *testing.T) {
	require.Equal(t, "hello", indirect("hello"))
}

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

func TestIsMultipleOf(t *testing.T) {
	require.True(t, isMultipleOf(10, 5))
	require.False(t, isMultipleOf(10, 3))
}

func TestNotMultipleOf(t *testing.T) {
	require.False(t, notMultipleOf(10, 5))
	require.True(t, notMultipleOf(10, 3))
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
