package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHead(t *testing.T) {

	require.Equal(t, "Hello", Head("Hello There", " "))

	require.Equal(t, "", Head(" Hello", " "))

	require.Equal(t, "Hello", Head("Hello ", " "))

	require.Equal(t, "Hello", Head("Hello", " "))

	require.Equal(t, "This is a very long string but still not a list.", Head("This is a very long string but still not a list.", ","))

	require.Equal(t, "One,Two,Three", Head("One,Two,Three", " "))
}

func TestTail(t *testing.T) {

	require.Equal(t, "There", Tail("Hello There", " "))

	require.Equal(t, "Bananas,Pears", Tail("Apples,Bananas,Pears", ","))

	require.Equal(t, "", Tail("One,Two,Three", " "))

	require.Equal(t, "Hello", Tail(" Hello", " "))

	require.Equal(t, "", Tail("Hello ", " "))

}

func TestLast(t *testing.T) {

	require.Equal(t, "Dude", Last("Hello There Dude", " "))

	require.Equal(t, "Hello", Last(" Hello", " "))

	require.Equal(t, "", Last("Hello ", " "))

	require.Equal(t, "Hello", Last("Hello", " "))

	require.Equal(t, "This is a very long string but still not a list.", Last("This is a very long string but still not a list.", ","))

	require.Equal(t, "One,Two,Three", Last("One,Two,Three", " "))

	require.Equal(t, "There", Last("Hello There", " "))

	require.Equal(t, "Pears", Last("Apples,Bananas,Pears", ","))

	require.Equal(t, "One,Two,Three", Last("One,Two,Three", " "))

	require.Equal(t, "Hello", Last(" Hello", " "))

	require.Equal(t, "", Last("Hello ", " "))

	require.Equal(t, "three", Last("one++two++three", "++"))
}

func TestRemoveLast(t *testing.T) {

	require.Equal(t, "Hello There", RemoveLast("Hello There Dude", " "))

	require.Equal(t, "", RemoveLast(" Hello", " "))

	require.Equal(t, "Hello", RemoveLast("Hello ", " "))

	require.Equal(t, "", RemoveLast("Hello", " "))

	require.Equal(t, "", RemoveLast("This is a very long string but still not a list.", ","))

	require.Equal(t, "", RemoveLast("One,Two,Three", " "))

	require.Equal(t, "Hello", RemoveLast("Hello There", " "))

	require.Equal(t, "Apples,Bananas", RemoveLast("Apples,Bananas,Pears", ","))

	require.Equal(t, "", RemoveLast("One,Two,Three", " "))

	require.Equal(t, "", RemoveLast(" Hello", " "))

	require.Equal(t, "Hello", RemoveLast("Hello ", " "))

}

func TestSplit(t *testing.T) {

	{
		a, b := Split("one,two,three", ",")
		require.Equal(t, "one", a)
		require.Equal(t, "two,three", b)
	}

	{
		a, b := Split("one,two,three", "!")
		require.Equal(t, "one,two,three", a)
		require.Equal(t, "", b)
	}

	{
		head, tail := Split("one++two++three", "++")
		require.Equal(t, "one", head)
		require.Equal(t, "two++three", tail)
	}
}

func TestSplitTail(t *testing.T) {

	{
		a, b := SplitTail("one,two,three", ",")
		require.Equal(t, "one,two", a)
		require.Equal(t, "three", b)
	}

	{
		a, b := SplitTail("one,two,three", "!")
		require.Equal(t, "one,two,three", a)
		require.Equal(t, "", b)
	}

	{
		head, tail := SplitTail("one+two+three", "+")
		require.Equal(t, "one+two", head)
		require.Equal(t, "three", tail)
	}

	{
		head, tail := SplitTail("one++two++three", "++")
		require.Equal(t, "one++two", head)
		require.Equal(t, "three", tail)
	}
}

func TestAt(t *testing.T) {

	{
		require.Equal(t, "one", At("one,two,three", ",", 0))
		require.Equal(t, "two", At("one,two,three", ",", 1))
		require.Equal(t, "three", At("one,two,three", ",", 2))
		require.Equal(t, "", At("one,two,three", ",", 3))
	}
}
