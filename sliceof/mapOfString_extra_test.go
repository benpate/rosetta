package sliceof

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func testMapOfString() MapOfString {
	return MapOfString{
		mapof.String{"name": "alpha"},
		mapof.String{"name": "bravo"},
		mapof.String{"name": "charlie"},
	}
}

func TestMapOfString_Manipulations(t *testing.T) {

	x := testMapOfString()

	require.Equal(t, 3, x.Length())
	require.True(t, x.IsLength(3))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, "alpha", x.First().GetString("name"))
	require.Equal(t, "charlie", x.Last().GetString("name"))
	require.Equal(t, 2, x.FirstN(2).Length())
	require.Equal(t, "bravo", x.At(1).GetString("name"))

	value, ok := x.AtOK(2)
	require.True(t, ok)
	require.Equal(t, "charlie", value.GetString("name"))

	_, ok = x.AtOK(99)
	require.False(t, ok)
}

func TestMapOfString_EmptyAccessors(t *testing.T) {
	x := NewMapOfString()
	require.True(t, x.IsZero())
	require.Equal(t, 0, x.First().Length())
	require.Equal(t, 0, x.Last().Length())
}

func TestMapOfString_FindFilterContains(t *testing.T) {

	x := testMapOfString()

	found, ok := x.Find(func(v mapof.String) bool { return v.GetString("name") == "bravo" })
	require.True(t, ok)
	require.Equal(t, "bravo", found.GetString("name"))

	_, ok = x.Find(func(v mapof.String) bool { return v.GetString("name") == "missing" })
	require.False(t, ok)

	filtered := x.Filter(func(v mapof.String) bool { return v.GetString("name") != "bravo" })
	require.Equal(t, 2, filtered.Length())

	require.True(t, x.Contains(func(v mapof.String) bool { return v.GetString("name") == "alpha" }))
	require.False(t, x.Contains(func(v mapof.String) bool { return v.GetString("name") == "missing" }))
}

func TestMapOfString_ReverseRange(t *testing.T) {

	x := testMapOfString()
	require.Equal(t, "charlie", x.Reverse().First().GetString("name"))

	collected := make([]string, 0)
	for _, value := range testMapOfString().Range() {
		collected = append(collected, value.GetString("name"))
	}
	require.Equal(t, []string{"alpha", "bravo", "charlie"}, collected)
}

func TestMapOfString_AppendShuffleKeys(t *testing.T) {

	x := NewMapOfString()
	x.Append(mapof.String{"name": "x"}, mapof.String{"name": "y"})
	require.Equal(t, 2, x.Length())

	shuffled := testMapOfString().Shuffle()
	require.Equal(t, 3, shuffled.Length())

	require.Equal(t, []string{"0", "1"}, x.Keys())
}

func TestMapOfString_GettersSetters(t *testing.T) {

	x := testMapOfString()

	value, ok := x.GetAnyOK("1")
	require.True(t, ok)
	require.Equal(t, "bravo", value.(mapof.String).GetString("name"))

	require.NotNil(t, x.GetAny("0"))

	indexValue, ok := x.GetIndex(0)
	require.True(t, ok)
	require.NotNil(t, indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)

	pointer, ok := x.GetPointer("0")
	require.True(t, ok)
	require.NotNil(t, pointer)

	require.True(t, x.SetIndex(5, map[string]string{"name": "grown"}))
	require.Equal(t, "grown", x.At(5).GetString("name"))
}

func TestMapOfString_SetValueRemove(t *testing.T) {

	x := NewMapOfString()
	require.NoError(t, x.SetValue(testMapOfString()))
	require.Equal(t, 3, x.Length())

	// A non-MapOfString value is rejected
	require.Error(t, x.SetValue("nope"))

	require.True(t, x.Remove("1"))
	require.Equal(t, 2, x.Length())
	require.False(t, x.Remove("bogus"))

	require.True(t, x.RemoveAt(0))
	require.Equal(t, 1, x.Length())
	require.False(t, x.RemoveAt(99))
}
