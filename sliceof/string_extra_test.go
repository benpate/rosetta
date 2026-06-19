package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString_Manipulations(t *testing.T) {

	x := NewString("a", "b", "c")

	require.Equal(t, 3, x.Length())
	require.True(t, x.IsLength(3))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, "a", x.First())
	require.Equal(t, "c", x.Last())
	require.Equal(t, "b", x.At(1))
	require.Equal(t, "", x.At(99))
}

func TestString_EmptyAccessors(t *testing.T) {
	x := NewString()
	require.True(t, x.IsZero())
	require.Equal(t, "", x.First())
	require.Equal(t, "", x.Last())
}

func TestString_FindFilter(t *testing.T) {

	x := NewString("apple", "banana", "cherry")

	found, ok := x.Find(func(v string) bool { return v == "banana" })
	require.True(t, ok)
	require.Equal(t, "banana", found)

	_, ok = x.Find(func(v string) bool { return v == "missing" })
	require.False(t, ok)
}

func TestString_ReverseRangeJoin(t *testing.T) {

	x := NewString("a", "b", "c")
	require.Equal(t, String{"c", "b", "a"}, x.Reverse())

	require.Equal(t, "c-b-a", x.Join("-"))

	collected := make([]string, 0)
	for _, value := range x.Range() {
		collected = append(collected, value)
	}
	require.Equal(t, []string{"c", "b", "a"}, collected)
}

func TestString_Contains(t *testing.T) {

	x := NewString("a", "b", "c")

	require.True(t, x.Contains("b"))
	require.False(t, x.Contains("z"))
	require.True(t, x.NotContains("z"))
	require.True(t, x.ContainsInterface("b"))
	require.False(t, x.ContainsInterface("z"))

	require.True(t, x.ContainsAny("z", "b"))
	require.False(t, x.ContainsAny("y", "z"))
	require.True(t, x.ContainsAll("a", "b"))
	require.False(t, x.ContainsAll("a", "z"))
}

func TestString_Equal(t *testing.T) {
	x := NewString("a", "b")
	require.True(t, x.Equal([]string{"a", "b"}))
	require.False(t, x.NotEqual([]string{"a", "b"}))
	require.True(t, x.NotEqual([]string{"a"}))
}

func TestString_AppendShuffleKeys(t *testing.T) {

	x := NewString("a")
	x.Append("b", "c")
	require.Equal(t, String{"a", "b", "c"}, x)

	shuffled := NewString("a", "b", "c", "d").Shuffle()
	require.Equal(t, 4, shuffled.Length())
	require.True(t, shuffled.ContainsAll("a", "b", "c", "d"))

	require.Equal(t, []string{"0", "1", "2"}, x.Keys())
}

func TestString_Getters(t *testing.T) {

	x := NewString("one", "two", "three")

	require.Equal(t, "two", x.GetString("1"))
	require.Equal(t, "three", x.GetString("last"))
	require.Equal(t, any("one"), x.GetAny("0"))

	value, ok := x.GetStringOK("2")
	require.Equal(t, "three", value)
	require.True(t, ok)

	_, ok = x.GetStringOK("bogus")
	require.False(t, ok)

	anyValue, ok := x.GetAnyOK("0")
	require.True(t, ok)
	require.Equal(t, "one", anyValue)

	indexValue, ok := x.GetIndex(1)
	require.True(t, ok)
	require.Equal(t, "two", indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)
}

func TestString_SettersAndSetIndex(t *testing.T) {

	x := NewString()

	require.True(t, x.SetString("0", "zero"))
	require.True(t, x.SetString("next", "one"))
	require.Equal(t, String{"zero", "one"}, x)
	require.True(t, x.SetString("last", "ONE"))
	require.Equal(t, "ONE", x.Last())
	require.False(t, x.SetString("bogus", "x"))

	require.True(t, x.SetIndex(4, 42))
	require.Equal(t, "42", x.At(4))
	require.Equal(t, 5, x.Length())
}

func TestString_SetValue(t *testing.T) {
	x := NewString()
	require.NoError(t, x.SetValue([]string{"x", "y"}))
	require.Equal(t, String{"x", "y"}, x)
}

func TestString_RemoveAt(t *testing.T) {

	x := NewString("a", "b", "c")
	require.True(t, x.RemoveAt(1))
	require.Equal(t, String{"a", "c"}, x)
	require.False(t, x.RemoveAt(99))
}
