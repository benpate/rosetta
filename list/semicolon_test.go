package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSemicolon_Empty(t *testing.T) {
	list := Semicolon("no delimiter here")
	require.Equal(t, "no delimiter here", list.Head())
	require.Equal(t, Semicolon(""), list.Tail())

	require.True(t, Semicolon("").IsEmpty())
	require.False(t, Semicolon("not-empty").IsEmpty())

	require.True(t, Semicolon("").IsEmptyTail())
	require.True(t, Semicolon("empty-tail").IsEmptyTail())
	require.True(t, Semicolon("empty-tail;").IsEmptyTail())
	require.False(t, Semicolon("empty;tail").IsEmptyTail())
	require.False(t, Semicolon("empty;tail;").IsEmptyTail())
}

func TestSemicolon_Head(t *testing.T) {

	list := Semicolon("hello;there;general;kenobi")
	require.Equal(t, list.Head(), list.First())
	require.Equal(t, "hello", list.Head())

	empty := Semicolon("")
	require.Equal(t, "", empty.Head())
}

func TestSemicolon_Tail(t *testing.T) {

	list := Semicolon("hello;there;general;kenobi")
	require.Equal(t, Semicolon("there;general;kenobi"), list.Tail())

	empty := Semicolon("")
	require.Equal(t, Semicolon(""), empty.Tail())
}

func TestSemicolon_Last(t *testing.T) {

	list := Semicolon("hello;there;general;kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Semicolon("")
	require.Equal(t, "", empty.Last())
}

func TestSemicolon_At(t *testing.T) {

	list := Semicolon("hello;there;general;kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestSemicolon_Push(t *testing.T) {

	list := BySemicolon("")
	list = list.PushHead("")
	require.Equal(t, Semicolon(""), list)

	list = list.PushTail("")
	require.Equal(t, Semicolon(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Semicolon("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Semicolon("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Semicolon("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Semicolon("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Semicolon("THERE;GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Semicolon("HELLO;THERE;GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Semicolon("HELLO;THERE;GENERAL;KENOBI"), list)
}

func TestSemicolon_Remove(t *testing.T) {

	list := Semicolon("hello;there;general;kenobi")

	require.Equal(t, Semicolon("there;general;kenobi"), list.Tail())
	require.Equal(t, Semicolon("hello;there;general"), list.RemoveLast())

	list2 := Semicolon("")
	require.Equal(t, Semicolon(""), list2.Tail())
	require.Equal(t, Semicolon(""), list2.RemoveLast())
}

func TestSemicolon_Split(t *testing.T) {

	empty := Semicolon("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Semicolon([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Semicolon([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Semicolon("hello;there;general;kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Semicolon("there;general;kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Semicolon("hello;there;general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Semicolon("there;general;kenobi"), list.Tail())
}

func TestSemicolon_String(t *testing.T) {
	require.Equal(t, "mystring", Semicolon("mystring").String())
}
