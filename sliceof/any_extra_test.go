package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAny_Manipulations(t *testing.T) {

	x := NewAny(1, 2, 3, 4)

	require.Equal(t, 4, x.Length())
	require.True(t, x.IsLength(4))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, any(1), x.First())
	require.Equal(t, any(4), x.Last())
	require.Equal(t, Any{1, 2}, x.FirstN(2))
	require.Equal(t, any(3), x.At(2))
	require.Nil(t, x.At(99))
}

func TestAny_EmptyAccessors(t *testing.T) {
	x := NewAny()
	require.True(t, x.IsZero())
	require.Nil(t, x.First())
	require.Nil(t, x.Last())
}

func TestAny_FindFilterReverseRange(t *testing.T) {

	x := NewAny(1, 2, 3, 4)

	found, ok := x.Find(func(v any) bool { return v == 3 })
	require.True(t, ok)
	require.Equal(t, 3, found)

	_, ok = x.Find(func(v any) bool { return v == 99 })
	require.False(t, ok)

	filtered := x.Filter(func(v any) bool { return v.(int) > 2 })
	require.Equal(t, Any{3, 4}, filtered)

	require.Equal(t, Any{4, 3, 2, 1}, x.Reverse())

	collected := make([]any, 0)
	for _, value := range NewAny("a", "b").Range() {
		collected = append(collected, value)
	}
	require.Equal(t, []any{"a", "b"}, collected)
}

func TestAny_Contains(t *testing.T) {

	x := NewAny(1, "two", 3.0)

	require.True(t, x.Contains(1))
	require.False(t, x.Contains("missing"))
	require.True(t, x.ContainsInterface(1))

	require.True(t, x.ContainsAny("missing", 1))
	require.False(t, x.ContainsAny("missing", "nope"))
	require.True(t, x.ContainsAll(1, "two"))
	require.False(t, x.ContainsAll(1, "missing"))
}

func TestAny_EqualAppendShuffleKeys(t *testing.T) {

	x := NewAny(1, 2)
	require.True(t, x.Equal([]any{1, 2}))
	require.False(t, x.Equal([]any{1}))

	x.Append(3)
	require.Equal(t, Any{1, 2, 3}, x)

	shuffled := NewAny(1, 2, 3, 4).Shuffle()
	require.Equal(t, 4, shuffled.Length())

	require.Equal(t, []string{"0", "1", "2"}, x.Keys())
}

func TestAny_TypedGetters(t *testing.T) {

	x := Any{true, 42, int64(99), 3.14, "hello"}

	require.True(t, x.GetBool("0"))
	require.Equal(t, 42, x.GetInt("1"))
	require.Equal(t, int64(99), x.GetInt64("2"))
	require.Equal(t, 3.14, x.GetFloat("3"))
	require.Equal(t, "hello", x.GetString("4"))
	require.Equal(t, any(true), x.GetAny("0"))

	b, ok := x.GetBoolOK("0")
	require.True(t, b)
	require.True(t, ok)

	i, ok := x.GetIntOK("1")
	require.Equal(t, 42, i)
	require.True(t, ok)

	i64, ok := x.GetInt64OK("2")
	require.Equal(t, int64(99), i64)
	require.True(t, ok)

	f, ok := x.GetFloatOK("3")
	require.Equal(t, 3.14, f)
	require.True(t, ok)

	s, ok := x.GetStringOK("4")
	require.Equal(t, "hello", s)
	require.True(t, ok)

	_, ok = x.GetStringOK("bogus")
	require.False(t, ok)

	anyValue, ok := x.GetAnyOK("0")
	require.True(t, ok)
	require.Equal(t, true, anyValue)

	indexValue, ok := x.GetIndex(1)
	require.True(t, ok)
	require.Equal(t, 42, indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)

	// GetPointer returns a pointer to the (any-typed) element
	pointer, ok := x.GetPointer("1")
	require.True(t, ok)
	require.NotNil(t, pointer)

	_, ok = x.GetPointer("bogus")
	require.False(t, ok)
}

func TestAny_Setters(t *testing.T) {

	x := NewAny()

	// Any setters only accept numeric string keys (no "next"/"last")
	require.True(t, x.SetBool("0", true))
	require.True(t, x.SetInt("1", 42))
	require.True(t, x.SetInt64("2", 99))
	require.True(t, x.SetFloat("3", 3.14))
	require.True(t, x.SetString("4", "hello"))
	require.True(t, x.SetAny("5", "world"))

	require.False(t, x.SetString("next", "rejected"))

	require.Equal(t, 6, x.Length())
	require.True(t, x.GetBool("0"))
	require.Equal(t, "world", x.GetString("5"))

	require.True(t, x.SetIndex(10, "grown"))
	require.Equal(t, "grown", x.At(10))
}

func TestAny_SetValue(t *testing.T) {
	x := NewAny()
	require.NoError(t, x.SetValue([]any{"a", "b"}))
	require.Equal(t, 2, x.Length())
}

func TestAny_Remove(t *testing.T) {

	x := NewAny(1, 2, 3)

	require.True(t, x.Remove("1"))
	require.Equal(t, Any{1, 3}, x)
	require.False(t, x.Remove("bogus"))

	require.True(t, x.RemoveAt(0))
	require.Equal(t, Any{3}, x)
	require.False(t, x.RemoveAt(99))
}
