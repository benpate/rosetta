package sliceof

import (
	"iter"
	"strconv"

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

// Find returns the first element in the slice that satisfies the provided function.
func (x Int) Find(fn func(int) bool) (int, bool) {
	return slice.Find(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Int) At(index int) int {
	return slice.At(x, index)
}

// Reverse modifies the slice with the elements in reverse order
func (x Int) Reverse() Int {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Range returns an iterator that yields each value in this slice.
func (x Int) Range() iter.Seq2[int, int] {
	return slice.Range(x)
}

// ContainsInterface returns TRUE if the provided generic value is contained in the slice.
func (x Int) ContainsInterface(value any) bool {

	// Convert the value to an int
	if value, ok := convert.IntOk(value, 0); ok {
		return slice.Contains(x, value)
	}

	// If we can't convert the value to a string, then it is not contained in the slice
	return false
}

// Contains returns TRUE if the slice contains the specified value
func (x Int) Contains(value int) bool {
	return slice.Contains(x, value)
}

func (x Int) NotContains(value int) bool {
	return !slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Int) ContainsAny(values ...int) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Int) ContainsAll(values ...int) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the "value" slice
func (x Int) Equal(value []int) bool {
	return slice.Equal(x, value)
}

// NotEqual returns TRUE if the slice DOES NOT contain exactly the same elements as the "value" slice
func (x Int) NotEqual(value []int) bool {
	return !slice.Equal(x, value)
}

// Append adds one or more elements to the end of the slice
func (x *Int) Append(values ...int) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Int) Shuffle() Int {
	return slice.Shuffle(x)
}

// Keys returns a slice of ints representing the indexes of this slice
func (x Int) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/****************************************
 * Getter Interfaces/Setters
 ****************************************/

func (x Int) GetAny(key string) any {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetAnyOK(key string) (any, bool) {
	return x.GetIntOK(key)
}

func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetIntOK(key string) (int, bool) {

	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
	}

	if key == "last" {
		return x.Last(), true
	}

	return 0, false
}

func (s *Int) SetIndex(index int, value any) bool {
	growSlice(s, index)
	(*s)[index] = convert.Int(value)
	return true
}

func (s *Int) SetInt(key string, value int) bool {
	if index, ok := sliceIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	switch key {

	case "last":
		return s.SetInt(strconv.Itoa(s.Length()-1), value)
	case "next":
		return s.SetInt(strconv.Itoa(s.Length()), value)
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
