package sliceof

import (
	"math/rand"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Int64 []int64

func NewInt64() Int64 {
	return make(Int64, 0)
}

/****************************************
 * Accessors
 ****************************************/

// Length returns the number of elements in the slice
func (x Int64) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Int64) IsLength(length int) bool {
	return len(x) == length
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Int64) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Int64) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Int64) First() int64 {
	if len(x) > 0 {
		return x[0]
	}
	return 0
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Int64) FirstN(n int) Int64 {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Int64) Last() int64 {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return 0
}

// Reverse returns a new slice with all elements in the opposite order
func (x Int64) Reverse() Int64 {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Copy returns a new slice with all elements copied from the original slice
func (x Int64) Contains(value int64) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Int64) ContainsAny(values ...int64) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Int64) ContainsAll(values ...int64) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the specified slice
func (x Int64) Equal(value []int64) bool {
	return slice.Equal(x, value)
}

// Append adds the specified values to the end of the slice
func (x *Int64) Append(values ...int64) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Int64) Shuffle() Int64 {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Int64) GetInt64OK(key string) (int64, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return 0, false
}

func (s *Int64) SetInt64(key string, value int64) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Int64) SetValue(value any) error {
	*s = convert.SliceOfInt64(value)
	return nil
}

func (s *Int64) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
