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

// REGRESSION: Matchable[T].MapOfAny used to call convert.MapOfAny(m). Matchable[T] is a
// map[string]T, so it matches none of that function's concrete map cases and falls through
// to its MapOfAnyGetter case -- which calls straight back into this method. The result was
// unbounded mutual recursion and a fatal stack overflow, for every T. If this regresses,
// the test binary dies rather than reporting a failure.
func TestMatchable_MapOfAny(t *testing.T) {

	m := newTestMatchable()
	result := m.MapOfAny()

	require.Len(t, result, 3)
	require.Equal(t, testMatcher{Name: "apple", Color: "red"}, result["a"])
	require.Equal(t, testMatcher{Name: "banana", Color: "yellow"}, result["b"])
	require.Equal(t, testMatcher{Name: "cherry", Color: "red"}, result["c"])
}

// The returned map must be a COPY -- writing to it must not reach back into the Matchable.
func TestMatchable_MapOfAny_IsACopy(t *testing.T) {

	m := newTestMatchable()
	result := m.MapOfAny()

	result["a"] = "clobbered"
	delete(result, "b")

	require.Equal(t, testMatcher{Name: "apple", Color: "red"}, m["a"])
	require.Len(t, m, 3)
}

func TestMatchable_MapOfAny_EmptyAndNil(t *testing.T) {

	require.Empty(t, NewMatchable[testMatcher]().MapOfAny())

	var nilMap Matchable[testMatcher]
	require.Empty(t, nilMap.MapOfAny(), "a nil map must not panic or recurse")
	require.NotNil(t, nilMap.MapOfAny(), "an empty map is returned, never nil")
}

// MapOfString goes through MapOfAny, so it shared the same stack overflow.
//
// The values come back EMPTY here, which is a separate limitation of convert.String: it
// type-switches on concrete types and has no reflect.String / reflect.Int fallback, so any
// named type (a struct, or even a `type Foo string`) stringifies to "". That is pinned
// below rather than asserted-against, so a future improvement to convert.String shows up
// here as a failure to be looked at.
func TestMatchable_MapOfString(t *testing.T) {

	m := NewMatchable[testMatcher]()
	m["a"] = testMatcher{Name: "apple", Color: "red"}

	result := m.MapOfString()

	require.Len(t, result, 1)
	require.Contains(t, result, "a")
	require.Equal(t, "", result["a"], "KNOWN LIMITATION: convert.String cannot render a named struct type")
}

func TestMatchable_MapOfString_EmptyAndNil(t *testing.T) {

	require.Empty(t, NewMatchable[testMatcher]().MapOfString())

	var nilMap Matchable[testMatcher]
	require.Empty(t, nilMap.MapOfString(), "a nil map must not panic or recurse")
}

// A string-backed Matchable round-trips through MapOfAny with its values intact.
// MapOfString still loses them, for the convert.String reason described above -- a
// `type Foo string` matches no case in its type switch either.
func TestMatchable_Conversions_StringValues(t *testing.T) {

	m := NewMatchable[stringMatcher]()
	m["greeting"] = "hello"

	require.Equal(t, map[string]any{"greeting": stringMatcher("hello")}, m.MapOfAny())
	require.Equal(t, map[string]string{"greeting": ""}, m.MapOfString(),
		"KNOWN LIMITATION: convert.String has no reflect.String fallback for named string types")
}

// stringMatcher is a string-backed Matcher, used to confirm the conversions work for a
// value type that convert.MapOfString can render losslessly.
type stringMatcher string

// Match reports whether this value satisfies the given predicate.
func (s stringMatcher) Match(predicate exp.Predicate) bool {
	return string(s) == predicate.Value
}
