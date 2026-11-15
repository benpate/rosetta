package deep

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNamed proves that Named matches whole, case-sensitive names
func TestNamed(t *testing.T) {

	match := Named("bto", "bcc")

	require.True(t, match("bto"))
	require.True(t, match("bcc"))
	require.False(t, match("BCC"))
	require.False(t, match("bccs"))
	require.False(t, match("bc"))
	require.False(t, match(""))

	t.Run("no names match nothing", func(t *testing.T) {
		require.False(t, Named()("bcc"))
		require.False(t, Named()(""))
	})

	t.Run("an empty name matches an empty key", func(t *testing.T) {
		require.True(t, Named("")(""))
		require.False(t, Named("")("a"))
	})

	t.Run("the caller's slice is copied", func(t *testing.T) {

		names := []string{"bcc"}
		match := Named(names...)
		names[0] = "to"

		require.True(t, match("bcc"))
		require.False(t, match("to"))
	})
}

// TestHasPrefix proves that HasPrefix matches names that begin with any of its prefixes
func TestHasPrefix(t *testing.T) {

	match := HasPrefix("x-", "emissary:")

	require.True(t, match("x-custom"))
	require.True(t, match("x-"))
	require.True(t, match("emissary:token"))
	require.False(t, match("X-custom"))
	require.False(t, match("x"))
	require.False(t, match("ax-"))
	require.False(t, match(""))

	t.Run("no prefixes match nothing", func(t *testing.T) {
		require.False(t, HasPrefix()("x-custom"))
	})

	t.Run("an empty prefix is ignored", func(t *testing.T) {

		// Kept, it would begin every name and so delete every property
		require.False(t, HasPrefix("")("bcc"))
		require.False(t, HasPrefix("", "")(""))
		require.True(t, HasPrefix("", "x-")("x-custom"))
		require.False(t, HasPrefix("", "x-")("bcc"))
	})

	t.Run("the caller's slice is copied", func(t *testing.T) {

		prefixes := []string{"x-", ""}
		match := HasPrefix(prefixes...)
		prefixes[0] = "y-"

		require.True(t, match("x-a"))
		require.False(t, match("y-a"))
		require.Equal(t, []string{"y-", ""}, prefixes)
	})
}

// TestAnyOf proves that AnyOf matches what at least one of its Matchers selects
func TestAnyOf(t *testing.T) {

	match := AnyOf(Named("bcc"), HasPrefix("x-"))

	require.True(t, match("bcc"))
	require.True(t, match("x-custom"))
	require.False(t, match("to"))

	t.Run("no matchers match nothing", func(t *testing.T) {
		require.False(t, AnyOf()("bcc"))
	})

	t.Run("a nil matcher is ignored", func(t *testing.T) {
		require.True(t, AnyOf(nil, Named("bcc"))("bcc"))
		require.False(t, AnyOf(nil, nil)("bcc"))
	})
}

// FuzzMatchers proves that the pre-built Matchers agree with a plain reading of their rules
func FuzzMatchers(f *testing.F) {

	f.Add("bcc", "bcc")
	f.Add("x-", "x-custom")
	f.Add("", "anything")
	f.Add("anything", "")
	f.Add("\xff", "\xff\xfe")

	f.Fuzz(func(t *testing.T, rule string, name string) {

		require.Equal(t, rule == name, Named(rule)(name))

		expected := rule != "" && len(name) >= len(rule) && name[:len(rule)] == rule
		require.Equal(t, expected, HasPrefix(rule)(name))

		require.Equal(t, rule == name || expected, AnyOf(Named(rule), HasPrefix(rule))(name))
	})
}
