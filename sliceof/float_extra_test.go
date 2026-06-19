package sliceof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat_Manipulations(t *testing.T) {

	x := NewFloat(1.0, 2.0, 3.0)

	require.Equal(t, 3, x.Length())
	require.True(t, x.IsLength(3))
	require.False(t, x.IsZero())
	require.False(t, x.IsEmpty())
	require.True(t, x.NotEmpty())

	require.Equal(t, 1.0, x.First())
	require.Equal(t, 3.0, x.Last())
	require.Equal(t, Float{1.0, 2.0}, x.FirstN(2))
	require.Equal(t, 2.0, x.At(1))
	require.Equal(t, 0.0, x.At(99))
}

func TestFloat_EmptyAccessors(t *testing.T) {
	x := NewFloat()
	require.True(t, x.IsZero())
	require.Equal(t, 0.0, x.First())
	require.Equal(t, 0.0, x.Last())
}

func TestFloat_FindFilterReverseRange(t *testing.T) {

	x := NewFloat(1.0, 2.0, 3.0, 4.0)

	found, ok := x.Find(func(v float64) bool { return v > 2 })
	require.True(t, ok)
	require.Equal(t, 3.0, found)

	_, ok = x.Find(func(v float64) bool { return v > 99 })
	require.False(t, ok)

	require.Equal(t, Float{2.0, 4.0}, x.Filter(func(v float64) bool { return int(v)%2 == 0 }))

	require.Equal(t, Float{4.0, 3.0, 2.0, 1.0}, x.Reverse())

	collected := make([]float64, 0)
	for _, value := range NewFloat(5.0, 6.0).Range() {
		collected = append(collected, value)
	}
	require.Equal(t, []float64{5.0, 6.0}, collected)
}

func TestFloat_Contains(t *testing.T) {

	x := NewFloat(1.0, 2.0, 3.0)

	require.True(t, x.Contains(2.0))
	require.False(t, x.Contains(9.0))
	require.True(t, x.NotContains(9.0))

	// Note: the method is named ContainsFloaterface in the source (a typo)
	require.True(t, x.ContainsFloaterface(2.0))
	require.True(t, x.ContainsFloaterface("2")) // coerced
	require.False(t, x.ContainsFloaterface(9.0))

	require.True(t, x.ContainsAny(9.0, 2.0))
	require.False(t, x.ContainsAny(8.0, 9.0))
	require.True(t, x.ContainsAll(1.0, 2.0))
	require.False(t, x.ContainsAll(1.0, 9.0))
}

func TestFloat_EqualAppendShuffleKeys(t *testing.T) {

	x := NewFloat(1.0, 2.0)
	require.True(t, x.Equal([]float64{1.0, 2.0}))
	require.False(t, x.NotEqual([]float64{1.0, 2.0}))
	require.True(t, x.NotEqual([]float64{9.0}))

	x.Append(3.0)
	require.Equal(t, Float{1.0, 2.0, 3.0}, x)

	shuffled := NewFloat(1.0, 2.0, 3.0, 4.0).Shuffle()
	require.Equal(t, 4, shuffled.Length())

	require.Equal(t, []string{"0", "1", "2"}, x.Keys())
}

func TestFloat_GettersSetters(t *testing.T) {

	x := NewFloat(10.0, 20.0, 30.0)

	require.Equal(t, 20.0, x.GetFloat("1"))
	require.Equal(t, 30.0, x.GetFloat("last"))
	require.Equal(t, any(10.0), x.GetAny("0"))

	value, ok := x.GetFloatOK("2")
	require.Equal(t, 30.0, value)
	require.True(t, ok)

	_, ok = x.GetFloatOK("bogus")
	require.False(t, ok)

	anyValue, ok := x.GetAnyOK("0")
	require.True(t, ok)
	require.Equal(t, 10.0, anyValue)

	indexValue, ok := x.GetIndex(1)
	require.True(t, ok)
	require.Equal(t, 20.0, indexValue)

	_, ok = x.GetIndex(99)
	require.False(t, ok)

	y := NewFloat()
	require.True(t, y.SetFloat("0", 1.5))
	require.True(t, y.SetFloat("next", 2.5))
	require.True(t, y.SetFloat("last", 9.9))
	require.False(t, y.SetFloat("bogus", 1))
	require.Equal(t, Float{1.5, 9.9}, y)

	require.True(t, y.SetIndex(4, 4.4))
	require.Equal(t, 4.4, y.At(4))
}

func TestFloat_SetValueRemove(t *testing.T) {

	x := NewFloat()
	require.NoError(t, x.SetValue([]float64{4.0, 5.0, 6.0}))
	require.Equal(t, Float{4.0, 5.0, 6.0}, x)

	require.True(t, x.Remove("1"))
	require.Equal(t, Float{4.0, 6.0}, x)
	require.False(t, x.Remove("bogus"))

	require.True(t, x.RemoveAt(0))
	require.Equal(t, Float{6.0}, x)
	require.False(t, x.RemoveAt(99))
}
