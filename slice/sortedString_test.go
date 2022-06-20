package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRemove1(t *testing.T) {

	s := []string{}

	s = AddUnique(s, "hello")
	require.Equal(t, []string{"hello"}, s)

	s = AddUnique(s, "hello")
	require.Equal(t, []string{"hello"}, s)

	s = AddUnique(s, "there")
	require.Equal(t, []string{"hello", "there"}, s)

	s = AddUnique(s, "there")
	require.Equal(t, []string{"hello", "there"}, s)

	s = AddUnique(s, "general")
	require.Equal(t, []string{"general", "hello", "there"}, s)

	s = AddUnique(s, "general")
	require.Equal(t, []string{"general", "hello", "there"}, s)

	s = AddUnique(s, "kenobi")
	require.Equal(t, []string{"general", "hello", "kenobi", "there"}, s)

	s = AddUnique(s, "kenobi")
	require.Equal(t, []string{"general", "hello", "kenobi", "there"}, s)

	s = Remove(s, "general")
	require.Equal(t, []string{"hello", "kenobi", "there"}, s)

	s = Remove(s, "there")
	require.Equal(t, []string{"hello", "kenobi"}, s)
}

func TestAddRemove2(t *testing.T) {

	s := []string{}

	s = AddUnique(s, "first")
	require.Equal(t, []string{"first"}, s)

	s = AddUnique(s, "second")
	require.Equal(t, []string{"first", "second"}, s)

	s = AddUnique(s, "third")
	require.Equal(t, []string{"first", "second", "third"}, s)

	s = AddUnique(s, "fourth")
	require.Equal(t, []string{"first", "fourth", "second", "third"}, s)

	s = AddUnique(s, "another first")
	require.Equal(t, []string{"another first", "first", "fourth", "second", "third"}, s)

	s = AddUnique(s, "very last")
	require.Equal(t, []string{"another first", "first", "fourth", "second", "third", "very last"}, s)

	s = AddUnique(s, "middle")
	require.Equal(t, []string{"another first", "first", "fourth", "middle", "second", "third", "very last"}, s)

	s = Remove(s, "missing element")
	require.Equal(t, []string{"another first", "first", "fourth", "middle", "second", "third", "very last"}, s)

	s = Remove(s, "another first")
	require.Equal(t, []string{"first", "fourth", "middle", "second", "third", "very last"}, s)

	s = Remove(s, "very last")
	require.Equal(t, []string{"first", "fourth", "middle", "second", "third"}, s)

	s = Remove(s, "middle")
	require.Equal(t, []string{"first", "fourth", "second", "third"}, s)

	s = Remove(s, "not in there")
	require.Equal(t, []string{"first", "fourth", "second", "third"}, s)

	s = Remove(s, "fourth")
	require.Equal(t, []string{"first", "second", "third"}, s)

	s = Remove(s, "third")
	require.Equal(t, []string{"first", "second"}, s)

	s = Remove(s, "first")
	require.Equal(t, []string{"second"}, s)

	s = Remove(s, "missing element")
	require.Equal(t, []string{"second"}, s)

	s = Remove(s, "second")
	require.Equal(t, []string{}, s)

	s = Remove(s, "first")
	require.Equal(t, []string{}, s)
}

func TestAddRemove3(t *testing.T) {

	s := []string{}

	s = AddUnique(s, "Hello", "There", "General", "Kenobi")
	require.Equal(t, []string{"General", "Hello", "Kenobi", "There"}, s)
}

func TestIdentical(t *testing.T) {

	{
		a := []string{"one", "two", "three"}
		b := []string{"one", "two", "three"}
		c := []string{"one", "two"}
		d := []string{"one", "three", "two"}
		e := []string{"one", "two", "three", "four"}

		require.True(t, Identical(a, b))
		require.False(t, Identical(a, c))
		require.False(t, Identical(a, d))
		require.False(t, Identical(a, e))
	}
}
