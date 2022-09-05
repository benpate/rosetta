package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComma_Empty(t *testing.T) {
	list := Comma("no delimiter here")
	require.Equal(t, "no delimiter here", list.Head())
	require.Equal(t, Comma(""), list.Tail())

	require.True(t, Comma("").IsEmpty())
	require.False(t, Comma("not-empty").IsEmpty())

	require.True(t, Comma("").IsEmptyTail())
	require.True(t, Comma("empty-tail").IsEmptyTail())
	require.True(t, Comma("empty-tail,").IsEmptyTail())
	require.False(t, Comma("empty,tail").IsEmptyTail())
	require.False(t, Comma("empty,tail,").IsEmptyTail())
}

func TestComma_Head(t *testing.T) {

	list := Comma("hello,there,general,kenobi")
	require.Equal(t, "hello", list.Head())

	empty := Comma("")
	require.Equal(t, "", empty.Head())
}

func TestComma_Tail(t *testing.T) {

	list := Comma("hello,there,general,kenobi")
	require.Equal(t, Comma("there,general,kenobi"), list.Tail())

	empty := Comma("")
	require.Equal(t, Comma(""), empty.Tail())
}

func TestComma_Last(t *testing.T) {

	list := Comma("hello,there,general,kenobi")
	require.Equal(t, "kenobi", list.Last())

	empty := Comma("")
	require.Equal(t, "", empty.Last())
}

func TestComma_At(t *testing.T) {

	list := Comma("hello,there,general,kenobi")
	require.Equal(t, "there", list.At(1))

	require.Equal(t, "", list.At(-10))
	require.Equal(t, "", list.At(5))
}

func TestComma_Push(t *testing.T) {

	list := ByComma("")
	list = list.PushHead("")
	require.Equal(t, Comma(""), list)

	list = list.PushTail("")
	require.Equal(t, Comma(""), list)

	list2 := list.PushTail("THERE")
	require.Equal(t, Comma("THERE"), list2)

	list = list.PushHead("THERE")
	require.Equal(t, Comma("THERE"), list)

	list = list.PushHead("")
	require.Equal(t, Comma("THERE"), list)

	list = list.PushTail("")
	require.Equal(t, Comma("THERE"), list)

	list = list.PushTail("GENERAL")
	require.Equal(t, Comma("THERE,GENERAL"), list)

	list = list.PushHead("HELLO")
	require.Equal(t, Comma("HELLO,THERE,GENERAL"), list)

	list = list.PushTail("KENOBI")
	require.Equal(t, Comma("HELLO,THERE,GENERAL,KENOBI"), list)
}

func TestComma_Remove(t *testing.T) {

	list := Comma("hello,there,general,kenobi")

	require.Equal(t, Comma("there,general,kenobi"), list.Tail())
	require.Equal(t, Comma("hello,there,general"), list.RemoveLast())

	list2 := Comma("")
	require.Equal(t, Comma(""), list2.Tail())
	require.Equal(t, Comma(""), list2.RemoveLast())
}

func TestComma_Split(t *testing.T) {

	empty := Comma("")

	{
		head, tail := empty.Split()
		require.Equal(t, "", head)
		require.Equal(t, Comma([]byte{}), tail)
	}

	{
		head, tail := empty.SplitTail()
		require.Equal(t, Comma([]byte{}), head)
		require.Equal(t, "", tail)
	}

	list := Comma("hello,there,general,kenobi")

	{
		head, tail := list.Split()
		require.Equal(t, "hello", head)
		require.Equal(t, Comma("there,general,kenobi"), tail)
	}

	{
		head, tail := list.SplitTail()
		require.Equal(t, Comma("hello,there,general"), head)
		require.Equal(t, "kenobi", tail)
	}

	require.Equal(t, Comma("there,general,kenobi"), list.Tail())
}

func TestComma_String(t *testing.T) {
	require.Equal(t, "mystring", Comma("mystring").String())
}
