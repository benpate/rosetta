package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlash_Empty(t *testing.T) {
	list := Slash("no delimiter here")
	require.Equal(t, "no delimiter here", list.Head())
	require.Equal(t, Slash(""), list.Tail())

	require.True(t, Slash("").IsEmpty())
	require.False(t, Slash("not-empty").IsEmpty())

	require.True(t, Slash("").IsEmptyTail())
	require.True(t, Slash("empty-tail").IsEmptyTail())
	require.True(t, Slash("empty-tail/").IsEmptyTail())
	require.False(t, Slash("empty/tail").IsEmptyTail())
	require.False(t, Slash("empty/tail/").IsEmptyTail())
}

func TestSlash_Head(t *testing.T) {

	list := Slash("hello/there/general/kenobi")
	require.Equal(t, "hello", list.Head())

	empty := Slash("")
	require.Equal(t, "", empty.Head())
}

func TestSlash_Tail(t *testing.T) {

	list := Slash("hello/there/general/kenobi")
	require.Equal(t, Slash("there/general/kenobi"), list.Tail())

	empty := Slash("")
	require.Equal(t, Slash(""), empty.Tail())
}

func TestSlash_Last(t *testing.T) {

	list := Slash("hello/there/general/kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Slash("")
	require.Equal(t, "", empty.Last())
}

func TestSlash_At(t *testing.T) {

	list := Slash("hello/there/general/kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestSlash_Push(t *testing.T) {

	list := BySlash("")
	list = list.PushHead("")
	require.Equal(t, Slash(""), list)

	list = list.PushTail("")
	require.Equal(t, Slash(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Slash("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Slash("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Slash("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Slash("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Slash("THERE/GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Slash("HELLO/THERE/GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Slash("HELLO/THERE/GENERAL/KENOBI"), list)
}

func TestSlash_Remove(t *testing.T) {

	list := Slash("hello/there/general/kenobi")

	require.Equal(t, Slash("there/general/kenobi"), list.Tail())
	require.Equal(t, Slash("hello/there/general"), list.RemoveLast())

	list2 := Slash("")
	require.Equal(t, Slash(""), list2.Tail())
	require.Equal(t, Slash(""), list2.RemoveLast())
}

func TestSlash_Split(t *testing.T) {

	empty := Slash("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Slash([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Slash([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Slash("hello/there/general/kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Slash("there/general/kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Slash("hello/there/general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Slash("there/general/kenobi"), list.Tail())
}

func TestSlash_String(t *testing.T) {
	require.Equal(t, "mystring", Slash("mystring").String())
}
