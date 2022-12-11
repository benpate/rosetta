package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDot_Empty(t *testing.T) {
	list := Dot("no delimiter here")
	require.Equal(t, "no delimiter here", list.Head())
	require.Equal(t, Dot(""), list.Tail())

	require.True(t, Dot("").IsEmpty())
	require.False(t, Dot("not-empty").IsEmpty())

	require.True(t, Dot("").IsEmptyTail())
	require.True(t, Dot("empty-tail").IsEmptyTail())
	require.True(t, Dot("empty-tail.").IsEmptyTail())
	require.False(t, Dot("empty.tail").IsEmptyTail())
	require.False(t, Dot("empty.tail.").IsEmptyTail())
}

func TestDot_Head(t *testing.T) {

	list := Dot("hello.there.general.kenobi")
	require.Equal(t, list.Head(), list.First())
	require.Equal(t, "hello", list.Head())

	empty := Dot("")
	require.Equal(t, "", empty.Head())
}

func TestDot_Tail(t *testing.T) {

	list := Dot("hello.there.general.kenobi")
	require.Equal(t, Dot("there.general.kenobi"), list.Tail())

	empty := Dot("")
	require.Equal(t, Dot(""), empty.Tail())
}

func TestDot_Last(t *testing.T) {

	list := Dot("hello.there.general.kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Dot("")
	require.Equal(t, "", empty.Last())
}

func TestDot_At(t *testing.T) {

	list := Dot("hello.there.general.kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestDot_Push(t *testing.T) {

	list := ByDot("")
	list = list.PushHead("")
	require.Equal(t, Dot(""), list)

	list = list.PushTail("")
	require.Equal(t, Dot(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Dot("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Dot("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Dot("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Dot("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Dot("THERE.GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Dot("HELLO.THERE.GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Dot("HELLO.THERE.GENERAL.KENOBI"), list)
}

func TestDot_Remove(t *testing.T) {

	list := Dot("hello.there.general.kenobi")

	require.Equal(t, Dot("there.general.kenobi"), list.Tail())
	require.Equal(t, Dot("hello.there.general"), list.RemoveLast())

	list2 := Dot("")
	require.Equal(t, Dot(""), list2.Tail())
	require.Equal(t, Dot(""), list2.RemoveLast())
}

func TestDot_Split(t *testing.T) {

	empty := Dot("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Dot([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Dot([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Dot("hello.there.general.kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Dot("there.general.kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Dot("hello.there.general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Dot("there.general.kenobi"), list.Tail())
}

func TestDot_String(t *testing.T) {
	require.Equal(t, "mystring", Dot("mystring").String())
}
