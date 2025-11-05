package sliceof

import (
	"iter"
	"strconv"
	"strings"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type String []string

func NewString() String {
	return make(String, 0)
}

/****************************************
 * Slice Manipulations
 ****************************************/

// Length returns the number of elements in the slice
func (x String) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x String) IsLength(length int) bool {
	return len(x) == length
}

// IsZero returns TRUE if the slice contains no elements.
// This is an alias for IsEmpty, and implements the `Zeroer`
// interface used by many packages (including go/json)
func (x String) IsZero() bool {
	return len(x) == 0
}

// IsEmpty returns TRUE if the slice contains no elements
func (x String) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x String) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x String) First() string {
	if len(x) > 0 {
		return x[0]
	}
	return ""
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x String) FirstN(n int) String {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x String) Last() string {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return ""
}

// Find returns the first element in the slice that satisfies the provided function.
func (x String) Find(fn func(string) bool) (string, bool) {
	return slice.Find(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x String) At(index int) string {
	return slice.At(x, index)
}

// Reverse modifies the slice with the elements in reverse order
func (x String) Reverse() String {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Range returns an iterator that yields each value in this slice.
func (x String) Range() iter.Seq2[int, string] {
	return slice.Range(x)
}

// ContainsInterface returns TRUE if the provided generic value is contained in the slice.
func (x String) ContainsInterface(value any) bool {

	// Convert the value to a string
	if value, ok := convert.StringOk(value, ""); ok {
		return slice.Contains(x, value)
	}

	// If we can't convert the value to a string, then it is not contained in the slice
	return false
}

// Contains returns TRUE if the slice contains the specified value
func (x String) Contains(value string) bool {
	return slice.Contains(x, value)
}

func (x String) NotContains(value string) bool {
	return !slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x String) ContainsAny(values ...string) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x String) ContainsAll(values ...string) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the "value" slice
func (x String) Equal(value []string) bool {
	return slice.Equal(x, value)
}

// NotEqual returns TRUE if the slice DOES NOT contain exactly the same elements as the "value" slice
func (x String) NotEqual(value []string) bool {
	return !slice.Equal(x, value)
}

// Join concatenates all elements of the slice into a single string, separated by the specified delimiter
func (x String) Join(delimiter string) string {
	return strings.Join(x, delimiter)
}

// Append adds one or more elements to the end of the slice
func (x *String) Append(values ...string) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x String) Shuffle() String {
	return slice.Shuffle(x)
}

// Keys returns a slice of strings representing the indexes of this slice
func (x String) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/****************************************
 * Getter Interfaces/Setters
 ****************************************/

func (x String) GetAny(key string) any {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetAnyOK(key string) (any, bool) {
	return x.GetStringOK(key)
}

func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetStringOK(key string) (string, bool) {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		return x[index], true
	}

	if key == "last" {
		return x.Last(), true
	}

	return "", false
}

func (s *String) SetIndex(index int, value any) bool {
	growSlice(s, index)
	(*s)[index] = convert.String(value)
	return true
}

func (s *String) SetString(key string, value string) bool {
	if index, ok := sliceStringIndex(key); ok {
		growSlice(s, index)
		(*s)[index] = value
		return true
	}

	switch key {

	case "last":
		return s.SetString(strconv.Itoa(s.Length()-1), value)
	case "next":
		return s.SetString(strconv.Itoa(s.Length()), value)
	}

	return false
}

func (s *String) SetValue(value any) error {
	*s = convert.SliceOfString(value)
	return nil
}

func (s *String) Remove(key string) bool {

	if index, ok := sliceStringIndex(key, s.Length()); ok {
		*s = append((*s)[:index], (*s)[index+1:]...)
		return true
	}

	return false
}

func (x *String) RemoveAt(index int) bool {

	if index, ok := sliceIndex(index, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
