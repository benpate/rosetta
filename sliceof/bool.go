package sliceof

import (
	"math/rand"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Bool []bool

func NewBool() Bool {
	return make(Bool, 0)
}

/******************************************
 * Slice Manipulations
 ******************************************/

// Length returns the number of elements in the slice
func (x Bool) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Bool) IsLength(length int) bool {
	return len(x) == length
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Bool) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Bool) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Bool) First() bool {
	if len(x) > 0 {
		return x[0]
	}
	return false
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Bool) FirstN(n int) Bool {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Bool) Last() bool {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return false
}

// Reverse returns a new slice with the elements in reverse order
func (x Bool) Reverse() Bool {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the slice contains the specified value
func (x Bool) Contains(value bool) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Bool) ContainsAny(values ...bool) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Bool) ContainsAll(values ...bool) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the specified value
func (x Bool) Equal(value []bool) bool {
	return slice.Equal(x, value)
}

// Append adds the specified values to the end of the slice
func (x *Bool) Append(values ...bool) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Bool) Shuffle() Bool {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

/****************************************
 * Getter/Setter Interfaces
 ****************************************/

func (x Bool) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Bool) GetBoolOK(key string) (bool, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	return false, false
}

func (s *Bool) SetBool(key string, value bool) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	return false
}

func (s *Bool) SetValue(value any) error {
	*s = convert.SliceOfBool(value)
	return nil
}

func (s *Bool) Remove(key string) bool {

	if index, ok := sliceIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}
