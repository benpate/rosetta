package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual_Empty(t *testing.T) {
	list := Equal("no delimiter here")
	require.Equal(t, "no delimiter here", list.Head())
	require.Equal(t, Equal(""), list.Tail())

	require.True(t, Equal("").IsEmpty())
	require.False(t, Equal("not-empty").IsEmpty())

	require.True(t, Equal("").IsEmptyTail())
	require.True(t, Equal("empty-tail").IsEmptyTail())
	require.True(t, Equal("empty-tail=").IsEmptyTail())
	require.False(t, Equal("empty=tail").IsEmptyTail())
	require.False(t, Equal("empty=tail=").IsEmptyTail())
}

func TestEqual_Head(t *testing.T) {

	list := Equal("hello=there=general=kenobi")
	require.Equal(t, list.Head(), list.First())
	require.Equal(t, "hello", list.Head())

	empty := Equal("")
	require.Equal(t, "", empty.Head())
}

func TestEqual_Tail(t *testing.T) {

	list := Equal("hello=there=general=kenobi")
	require.Equal(t, Equal("there=general=kenobi"), list.Tail())

	empty := Equal("")
	require.Equal(t, Equal(""), empty.Tail())
}

func TestEqual_Last(t *testing.T) {

	list := Equal("hello=there=general=kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Equal("")
	require.Equal(t, "", empty.Last())
}

func TestEqual_At(t *testing.T) {

	list := Equal("hello=there=general=kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestEqual_Push(t *testing.T) {

	list := ByEqual("")
	list = list.PushHead("")
	require.Equal(t, Equal(""), list)

	list = list.PushTail("")
	require.Equal(t, Equal(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Equal("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Equal("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Equal("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Equal("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Equal("THERE=GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Equal("HELLO=THERE=GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Equal("HELLO=THERE=GENERAL=KENOBI"), list)
}

func TestEqual_Remove(t *testing.T) {

	list := Equal("hello=there=general=kenobi")

	require.Equal(t, Equal("there=general=kenobi"), list.Tail())
	require.Equal(t, Equal("hello=there=general"), list.RemoveLast())

	list2 := Equal("")
	require.Equal(t, Equal(""), list2.Tail())
	require.Equal(t, Equal(""), list2.RemoveLast())
}

func TestEqual_Split(t *testing.T) {

	empty := Equal("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Equal([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Equal([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Equal("hello=there=general=kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Equal("there=general=kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Equal("hello=there=general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Equal("there=general=kenobi"), list.Tail())
}

func TestEqual_String(t *testing.T) {
	require.Equal(t, "mystring", Equal("mystring").String())
}
