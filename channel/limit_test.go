package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimit3(t *testing.T) {

	done := make(Done)
	in := Limit(3, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}

func TestLimit4(t *testing.T) {

	done := make(Done)
	in := Limit(4, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.Equal(t, "Kenobi", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}

func TestLimit5(t *testing.T) {

	done := make(Done)
	in := Limit(5, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.Equal(t, "Kenobi", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}

func TestLimit3_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(3, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.True(t, Closed(in))
}

func TestLimit4_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(4, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.Equal(t, "Kenobi", <-in)
	require.True(t, Closed(in))
}

func TestLimit5_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(5, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "General", <-in)
	require.Equal(t, "Kenobi", <-in)
	require.True(t, Closed(in))
}
