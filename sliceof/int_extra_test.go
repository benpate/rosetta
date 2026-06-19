package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt_Manipulations(t *testing.T) {

	x := NewInt(1, 2, 3, 4, 5)

	require.Equal(t, 5, x.Length())
	require.True(t, x.IsLength(5))
	require.False(t, x.IsLength(2))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, 1, x.First())
	require.Equal(t, 5, x.Last())
	require.Equal(t, Int{1, 2, 3}, x.FirstN(3))
	require.Equal(t, x, x.FirstN(99)) // n greater than length returns all
	require.Equal(t, 3, x.At(2))
	require.Equal(t, 0, x.At(99)) // out of bounds returns zero value
}

func TestInt_EmptyAccessors(t *testing.T) {
	x := NewInt()
	require.True(t, x.IsZero())
	require.True(t, x.IsEmpty())
	require.False(t, x.NotEmpty())
	require.Equal(t, 0, x.First())
	require.Equal(t, 0, x.Last())
}

func TestInt_FindFilter(t *testing.T) {

	x := NewInt(1, 2, 3, 4)

	found, ok := x.Find(func(v int) bool { return v > 2 })
	require.True(t, ok)
	require.Equal(t, 3, found)

	_, ok = x.Find(func(v int) bool { return v > 99 })
	require.False(t, ok)

	evens := x.Filter(func(v int) bool { return v%2 == 0 })
	require.Equal(t, Int{2, 4}, evens)
}

func TestInt_Reverse(t *testing.T) {
	x := NewInt(1, 2, 3)
	require.Equal(t, Int{3, 2, 1}, x.Reverse())
}

func TestInt_Range(t *testing.T) {

	x := NewInt(10, 20, 30)

	collected := make([]int, 0)
	for _, value := range x.Range() {
		collected = append(collected, value)
	}

	require.Equal(t, []int{10, 20, 30}, collected)
}

func TestInt_Contains(t *testing.T) {

	x := NewInt(1, 2, 3)

	require.True(t, x.Contains(2))
	require.False(t, x.Contains(9))
	require.True(t, x.NotContains(9))
	require.False(t, x.NotContains(2))

	require.True(t, x.ContainsInterface(2))
	require.True(t, x.ContainsInterface("2")) // coerced
	require.False(t, x.ContainsInterface(9))

	require.True(t, x.ContainsAny(9, 2))
	require.False(t, x.ContainsAny(8, 9))
	require.True(t, x.ContainsAll(1, 2, 3))
	require.False(t, x.ContainsAll(1, 2, 9))
}

func TestInt_Equal(t *testing.T) {
	x := NewInt(1, 2, 3)
	require.True(t, x.Equal([]int{1, 2, 3}))
	require.False(t, x.NotEqual([]int{1, 2, 3}))
	require.False(t, x.Equal([]int{1, 2}))
	require.True(t, x.NotEqual([]int{9}))
}

func TestInt_AppendShuffleKeys(t *testing.T) {

	x := NewInt(1)
	x.Append(2, 3)
	require.Equal(t, Int{1, 2, 3}, x)

	// Shuffle preserves length and elements
	shuffled := NewInt(1, 2, 3, 4, 5).Shuffle()
	require.Equal(t, 5, shuffled.Length())
	require.True(t, shuffled.ContainsAll(1, 2, 3, 4, 5))

	require.Equal(t, []string{"0", "1", "2"}, x.Keys())
}

func TestInt_Getters(t *testing.T) {

	x := NewInt(10, 20, 30)

	require.Equal(t, 20, x.GetInt("1"))
	require.Equal(t, 30, x.GetInt("last"))
	require.Equal(t, any(10), x.GetAny("0"))

	value, ok := x.GetIntOK("2")
	require.Equal(t, 30, value)
	require.True(t, ok)

	_, ok = x.GetIntOK("bogus")
	require.False(t, ok)

	any1, ok := x.GetAnyOK("0")
	require.True(t, ok)
	require.Equal(t, 10, any1)

	indexValue, ok := x.GetIndex(1)
	require.True(t, ok)
	require.Equal(t, 20, indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)
}

func TestInt_Setters(t *testing.T) {

	x := NewInt()

	require.True(t, x.SetInt("0", 100))
	require.True(t, x.SetInt("next", 200))
	require.Equal(t, Int{100, 200}, x)

	require.True(t, x.SetInt("last", 250))
	require.Equal(t, 250, x.Last())

	// Non-numeric, non-keyword keys are rejected
	require.False(t, x.SetInt("bogus", 1))

	require.True(t, x.SetIndex(5, 999))
	require.Equal(t, 999, x.At(5))
	require.Equal(t, 6, x.Length()) // slice grew to fit
}

func TestInt_SetValue(t *testing.T) {
	x := NewInt()
	require.NoError(t, x.SetValue([]int{4, 5, 6}))
	require.Equal(t, Int{4, 5, 6}, x)
}

func TestInt_Remove(t *testing.T) {

	x := NewInt(1, 2, 3)

	require.True(t, x.Remove("1"))
	require.Equal(t, Int{1, 3}, x)

	require.False(t, x.Remove("bogus"))

	require.True(t, x.RemoveAt(0))
	require.Equal(t, Int{3}, x)

	require.False(t, x.RemoveAt(99))
}
