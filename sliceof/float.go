package sliceof

import (
	"math/rand"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Float []float64

func NewFloat() Float {
	return make(Float, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

// Length returns the number of elements in the slice
func (x Float) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Float) IsLength(length int) bool {
	return len(x) == length
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Float) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Float) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Float) First() float64 {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Float) FirstN(n int) Float {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Float) Last() float64 {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

// Reverse returns a new slice with the elements in reverse order
func (x Float) Reverse() Float {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the slice contains the specified value
func (x Float) Contains(value float64) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Float) ContainsAny(values ...float64) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Float) ContainsAll(values ...float64) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the specified value
func (x Float) Equal(value []float64) bool {
	return slice.Equal(x, value)
}

// Append adds one or more elements to the end of the slice
func (x *Float) Append(values ...float64) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Float) Shuffle() Float {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetFloatOK(key string) (float64, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Float) SetFloat(key string, value float64) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Float) SetValue(value any) error {
	*s = convert.SliceOfFloat(value)
	return nil
}

func (s *Float) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
