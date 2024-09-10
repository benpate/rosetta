package sliceof

import (
	"math/rand"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Int []int

func NewInt() Int {
	return make(Int, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

// Length returns the number of elements in the slice
func (x Int) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Int) IsLength(length int) bool {
	return len(x) == length
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Int) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Int) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Int) First() int {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Int) FirstN(n int) Int {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Int) Last() int {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Int) At(index int) int {
	return slice.At(x, index)
}

// Reverse returns the elements of the slice in reverse order
func (x Int) Reverse() Int {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the slice contains the specified value
func (x Int) Contains(value int) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Int) ContainsAny(values ...int) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Int) ContainsAll(values ...int) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the specified value
func (x Int) Equal(value []int) bool {
	return slice.Equal(x, value)
}

// Append adds the specified values to the end of the slice
func (x *Int) Append(values ...int) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Int) Shuffle() Int {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetIntOK(key string) (int, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Int) SetInt(key string, value int) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Int) SetValue(value any) error {
	*s = convert.SliceOfInt(value)
	return nil
}

func (s *Int) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
