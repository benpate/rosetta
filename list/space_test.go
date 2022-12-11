package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpace_Empty(t *testing.T) {
	list := Space("nodelimiterhere")
	require.Equal(t, "nodelimiterhere", list.Head())
	require.Equal(t, Space(""), list.Tail())

	require.True(t, Space("").IsEmpty())
	require.False(t, Space("not-empty").IsEmpty())

	require.True(t, Space("").IsEmptyTail())
	require.True(t, Space("empty-tail").IsEmptyTail())
	require.True(t, Space("empty-tail ").IsEmptyTail())
	require.False(t, Space("empty tail").IsEmptyTail())
	require.False(t, Space("empty tail ").IsEmptyTail())
}

func TestSpace_Head(t *testing.T) {

	list := Space("hello there general kenobi")
	require.Equal(t, list.Head(), list.First())
	require.Equal(t, "hello", list.Head())

	empty := Space("")
	require.Equal(t, "", empty.Head())
}

func TestSpace_Tail(t *testing.T) {

	list := Space("hello there general kenobi")
	require.Equal(t, Space("there general kenobi"), list.Tail())

	empty := Space("")
	require.Equal(t, Space(""), empty.Tail())
}

func TestSpace_Last(t *testing.T) {

	list := Space("hello there general kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Space("")
	require.Equal(t, "", empty.Last())
}

func TestSpace_At(t *testing.T) {

	list := Space("hello there general kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestSpace_Push(t *testing.T) {

	list := BySpace("")
	list = list.PushHead("")
	require.Equal(t, Space(""), list)

	list = list.PushTail("")
	require.Equal(t, Space(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Space("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Space("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Space("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Space("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Space("THERE GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Space("HELLO THERE GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Space("HELLO THERE GENERAL KENOBI"), list)
}

func TestSpace_Remove(t *testing.T) {

	list := Space("hello there general kenobi")

	require.Equal(t, Space("there general kenobi"), list.Tail())
	require.Equal(t, Space("hello there general"), list.RemoveLast())

	list2 := Space("")
	require.Equal(t, Space(""), list2.Tail())
	require.Equal(t, Space(""), list2.RemoveLast())
}

func TestSpace_Split(t *testing.T) {

	empty := Space("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Space([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Space([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Space("hello there general kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Space("there general kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Space("hello there general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Space("there general kenobi"), list.Tail())
}

func TestSpace_String(t *testing.T) {
	require.Equal(t, "mystring", Space("mystring").String())
}
