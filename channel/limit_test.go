package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimit0(t *testing.T) {

	done := make(Done)
	in := Limit(0, testChannel(), done)

	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}

func TestLimit1(t *testing.T) {

	done := make(Done)
	in := Limit(1, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}

func TestLimit2(t *testing.T) {

	done := make(Done)
	in := Limit(2, testChannel(), done)

	require.Equal(t, "Hello", <-in)
	require.Equal(t, "There", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
	require.Equal(t, "", <-in)
}
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

func TestLimit0_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(0, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 0, counter)
	require.True(t, Closed(done))
	require.True(t, Closed(in))
}

func TestLimit1_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(1, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 1, counter)

	require.True(t, Closed(done))
	require.True(t, Closed(in))
}

func TestLimit2_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(2, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 2, counter)

	require.True(t, Closed(done))
	require.True(t, Closed(in))
}

func TestLimit3_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(3, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 3, counter)

	require.True(t, Closed(done))
	require.True(t, Closed(in))
}

func TestLimit4_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(4, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 4, counter)

	require.True(t, Closed(done))
	require.True(t, Closed(in))
}

func TestLimit5_Closed(t *testing.T) {

	done := make(Done)
	in := Limit(5, testChannel(), done)
	counter := 0

	for range in {
		counter++
	}

	require.Equal(t, 4, counter)

	require.True(t, Closed(done))
	require.True(t, Closed(in))
}
