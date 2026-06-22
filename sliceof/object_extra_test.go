package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObject_Manipulations(t *testing.T) {

	x := NewObject[int](1, 2, 3, 4)

	require.Equal(t, 4, x.Length())
	require.True(t, x.IsLength(4))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, 1, x.First())
	require.Equal(t, 4, x.Last())
	require.Equal(t, Object[int]{1, 2}, x.FirstN(2))
	require.Equal(t, 3, x.At(2))
	require.Equal(t, 0, x.At(99))

	value, ok := x.AtOK(1)
	require.True(t, ok)
	require.Equal(t, 2, value)

	_, ok = x.AtOK(99)
	require.False(t, ok)
}

func TestObject_EmptyAccessors(t *testing.T) {
	x := NewObject[int]()
	require.True(t, x.IsZero())
	require.Equal(t, 0, x.First())
	require.Equal(t, 0, x.Last())
}

func TestObject_FindFilterContains(t *testing.T) {

	x := NewObject[int](1, 2, 3, 4)

	found, ok := x.Find(func(v int) bool { return v > 2 })
	require.True(t, ok)
	require.Equal(t, 3, found)

	_, ok = x.Find(func(v int) bool { return v > 99 })
	require.False(t, ok)

	require.Equal(t, Object[int]{2, 4}, x.Filter(func(v int) bool { return v%2 == 0 }))

	require.True(t, x.Contains(func(v int) bool { return v == 3 }))
	require.False(t, x.Contains(func(v int) bool { return v == 99 }))
}

func TestObject_Range(t *testing.T) {

	collected := make([]int, 0)
	for _, value := range NewObject[int](7, 8, 9).Range() {
		collected = append(collected, value)
	}
	require.Equal(t, []int{7, 8, 9}, collected)
}

func TestObject_AppendShuffleKeys(t *testing.T) {

	x := NewObject[int](1)
	x.Append(2, 3)
	require.Equal(t, Object[int]{1, 2, 3}, x)

	shuffled := NewObject[int](1, 2, 3, 4).Shuffle()
	require.Equal(t, 4, shuffled.Length())

	require.Equal(t, []string{"0", "1", "2"}, x.Keys())
}

func TestObject_GettersSetters(t *testing.T) {

	x := NewObject[string]("a", "b", "c")

	require.Equal(t, any("b"), x.GetAny("1"))

	value, ok := x.GetAnyOK("2")
	require.True(t, ok)
	require.Equal(t, "c", value)

	_, ok = x.GetAnyOK("bogus")
	require.False(t, ok)

	indexValue, ok := x.GetIndex(0)
	require.True(t, ok)
	require.Equal(t, "a", indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)

	pointer, ok := x.GetPointer("0")
	require.True(t, ok)
	require.NotNil(t, pointer)

	require.True(t, x.SetIndex(5, "grown"))
	require.Equal(t, "grown", x.At(5))
}

func TestObject_SetValue(t *testing.T) {

	// SetValue only accepts Object[T] (or *Object[T]) values
	x := NewObject[string]()
	require.NoError(t, x.SetValue(Object[string]{"x", "y"}))
	require.Equal(t, 2, x.Length())

	source := Object[string]{"a"}
	require.NoError(t, x.SetValue(&source))
	require.Equal(t, 1, x.Length())

	// A plain slice cannot be converted
	require.Error(t, x.SetValue([]string{"nope"}))
}

func TestObject_Remove(t *testing.T) {

	x := NewObject[int](1, 2, 3)

	require.True(t, x.Remove("1"))
	require.Equal(t, Object[int]{1, 3}, x)
	require.False(t, x.Remove("bogus"))

	require.True(t, x.RemoveAt(0))
	require.Equal(t, Object[int]{3}, x)
	require.False(t, x.RemoveAt(99))
}
