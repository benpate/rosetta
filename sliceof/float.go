package sliceof

import (
	"iter"
	"strconv"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Float []float64

func NewFloat(values ...float64) Float {

	if len(values) == 0 {
		return make(Float, 0)
	}

	return values
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

// IsZero returns TRUE if the slice contains no elements.
// This is an alias for IsEmpty, and implements the `Zeroer`
// interface used by many packages (including go/json)
func (x Float) IsZero() bool {
	return len(x) == 0
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

// Find returns the first element in the slice that satisfies the provided function.
func (x Float) Find(fn func(float64) bool) (float64, bool) {
	return slice.Find(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Float) At(index int) float64 {
	return slice.At(x, index)
}

// Reverse modifies the slice with the elements in reverse order
func (x Float) Reverse() Float {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Range returns an iterator that yields each value in this slice.
func (x Float) Range() iter.Seq2[int, float64] {
	return slice.Range(x)
}

// ContainsFloaterface returns TRUE if the provided generic value is contained in the slice.
func (x Float) ContainsFloaterface(value any) bool {

	// Convert the value to an float64
	if value, ok := convert.FloatOk(value, 0); ok {
		return slice.Contains(x, value)
	}

	// If we can't convert the value to a string, then it is not contained in the slice
	return false
}

// Contains returns TRUE if the slice contains the specified value
func (x Float) Contains(value float64) bool {
	return slice.Contains(x, value)
}

func (x Float) NotContains(value float64) bool {
	return !slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Float) ContainsAny(values ...float64) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Float) ContainsAll(values ...float64) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the "value" slice
func (x Float) Equal(value []float64) bool {
	return slice.Equal(x, value)
}

// NotEqual returns TRUE if the slice DOES NOT contain exactly the same elements as the "value" slice
func (x Float) NotEqual(value []float64) bool {
	return !slice.Equal(x, value)
}

// Append adds one or more elements to the end of the slice
func (x *Float) Append(values ...float64) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Float) Shuffle() Float {
	return slice.Shuffle(x)
}

// Keys returns a slice of float64s representing the indexes of this slice
func (x Float) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/****************************************
 * Getter Floaterfaces/Setters
 ****************************************/

func (x Float) GetAny(key string) any {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetAnyOK(key string) (any, bool) {
	return x.GetFloatOK(key)
}

func (x Float) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Float) GetFloatOK(key string) (float64, bool) {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		return x[index], true
	}

	if key == "last" {
		return x.Last(), true
	}

	return 0, false
}

func (s *Float) SetIndex(index int, value any) bool {
	growSlice(s, index)
	(*s)[index] = convert.Float(value)
	return true
}

func (s *Float) SetFloat(key string, value float64) bool {
	if index, ok := sliceStringIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	switch key {

	case "last":
		return s.SetFloat(strconv.Itoa(s.Length()-1), value)
	case "next":
		return s.SetFloat(strconv.Itoa(s.Length()), value)
	}

	return false
}

func (s *Float) SetValue(value any) error {
	*s = convert.SliceOfFloat(value)
	return nil
}

func (s *Float) Remove(key string) bool {

	if index, ok := sliceStringIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}

func (x *Float) RemoveAt(index int) bool {

	if index, ok := sliceIndex(index, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
