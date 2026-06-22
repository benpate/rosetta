package mapof

import (
	"sort"
	"testing"

	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

// testMatcher is a minimal Matcher used to exercise the Matchable map.
type testMatcher struct {
	Name  string
	Color string
}

// Match reports whether this object satisfies the given predicate. It only
// understands "name" and "color" equality predicates.
func (t testMatcher) Match(predicate exp.Predicate) bool {
	switch predicate.Field {
	case "name":
		return t.Name == predicate.Value
	case "color":
		return t.Color == predicate.Value
	}
	return false
}

func newTestMatchable() Matchable[testMatcher] {
	m := NewMatchable[testMatcher]()
	m["a"] = testMatcher{Name: "apple", Color: "red"}
	m["b"] = testMatcher{Name: "banana", Color: "yellow"}
	m["c"] = testMatcher{Name: "cherry", Color: "red"}
	return m
}

func TestMatchable_Match(t *testing.T) {

	m := newTestMatchable()

	reds := m.Match(exp.Equal("color", "red"))
	require.Equal(t, 2, reds.Length())
	require.Contains(t, reds, "a")
	require.Contains(t, reds, "c")

	none := m.Match(exp.Equal("color", "purple"))
	require.Equal(t, 0, none.Length())
}

func TestMatchable_MatchOne(t *testing.T) {

	m := newTestMatchable()

	value, ok := m.MatchOne(exp.Equal("name", "banana"))
	require.True(t, ok)
	require.Equal(t, "banana", value.Name)

	_, ok = m.MatchOne(exp.Equal("name", "missing"))
	require.False(t, ok)
}

func TestMatchable_Manipulations(t *testing.T) {

	m := newTestMatchable()

	require.Equal(t, 3, m.Length())
	require.False(t, m.IsEmpty())
	require.True(t, m.NotEmpty())
	require.Equal(t, []string{"a", "b", "c"}, m.Keys())

	values := m.Values()
	require.Equal(t, 3, len(values))
	names := []string{values[0].Name, values[1].Name, values[2].Name}
	sort.Strings(names)
	require.Equal(t, []string{"apple", "banana", "cherry"}, names)

	require.True(t, m.Remove("a"))
	require.Equal(t, 2, m.Length())
}

func TestMatchable_GetPointer(t *testing.T) {

	m := newTestMatchable()

	value, ok := m.GetPointer("a")
	require.True(t, ok)
	require.Equal(t, "apple", value.(testMatcher).Name)

	_, ok = m.GetPointer("missing")
	require.False(t, ok)
}

func TestMatchable_Empty(t *testing.T) {

	m := NewMatchable[testMatcher]()
	require.True(t, m.IsEmpty())
	require.Equal(t, 0, len(m.Values()))
	require.Equal(t, 0, len(m.Keys()))
}

func TestMatchable_NilRemove(t *testing.T) {
	var m Matchable[testMatcher]
	require.True(t, m.Remove("key"))
	require.NotNil(t, m)
}

func TestMatchable_EqualNotEqual(t *testing.T) {
	m := newTestMatchable()
	// Equal compares against an Any-typed map; a Matchable of structs is never
	// DeepEqual to a map[string]any, so these exercise both branches.
	require.False(t, m.Equal(map[string]any{"a": "apple"}))
	require.True(t, m.NotEqual(map[string]any{"a": "apple"}))
}

// NOTE: Matchable.MapOfAny() / MapOfString() are intentionally NOT tested here.
// They call convert.MapOfAny(m), and because Matchable defines a MapOfAny()
// method, convert re-dispatches back into Matchable.MapOfAny(), causing
// unbounded recursion / stack overflow. This is a bug in the library, not the
// test; testing it would crash the whole package.

func TestMatchable_IsZeroValue(t *testing.T) {
	m := newTestMatchable()
	// compare.IsZero does not recognize arbitrary struct types, so both a
	// populated entry and a missing (zero-struct) entry report not-zero.
	require.False(t, m.IsZeroValue("a"))
	require.False(t, m.IsZeroValue("missing"))
}
