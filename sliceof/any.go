package sliceof

import (
	"iter"
	"math/rand"
	"strconv"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/pointer"
	"github.com/benpate/rosetta/slice"
)

type Any []any

func NewAny(values ...any) Any {

	if len(values) == 0 {
		return make(Any, 0)
	}

	return values
}

/****************************************
 * Slice Manipulations
 ****************************************/

// Length returns the number of elements in the slice
func (x Any) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x Any) IsLength(length int) bool {
	return len(x) == length
}

// IsZero returns TRUE if the slice contains no elements.
// This is an alias for IsEmpty, and implements the `Zeroer`
// interface used by many packages (including go/json)
func (x Any) IsZero() bool {
	return len(x) == 0
}

// IsEmpty returns TRUE if the slice contains no elements
func (x Any) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Any) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice,
// or nil if the slice is empty
func (x Any) First() any {
	if len(x) > 0 {
		return x[0]
	}
	return nil
}

// FirstN returns the first "n" elements in the slice,
// or all elements if "n" is greater than the length of the slice
func (x Any) FirstN(n int) Any {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice,
// or nil if the slice is empty
func (x Any) Last() any {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return nil
}

// Find returns the first element in the slice that satisfies the provided function.
func (x Any) Find(fn func(any) bool) (any, bool) {
	return slice.Find(x, fn)
}

// Filter returns all elements in the slice that satisfies the provided function.
func (x Any) Filter(fn func(any) bool) Any {
	return slice.Filter(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x Any) At(index int) any {
	return slice.At(x, index)
}

// Reverse returns a new slice with the elements in reverse order
func (x Any) Reverse() Any {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Range returns an iterator that yields each value in this slice.
func (x Any) Range() iter.Seq2[int, any] {
	return slice.Range(x)
}

// ContainsInterface returns TRUE if the provided generic value is contained in the slice.
func (x Any) ContainsInterface(value any) bool {
	return slice.Contains(x, value)
}

// Contains returns TRUE if the slice contains the specified value
func (x Any) Contains(value any) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains
// any of the specified values
func (x Any) ContainsAny(values ...any) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains
// all of the specified values
func (x Any) ContainsAll(values ...any) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains
// exactly the same elements as the specified value
func (x Any) Equal(value []any) bool {
	return slice.Equal(x, value)
}

// Append adds one or more elements to the end of the slice
func (x *Any) Append(values ...any) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x Any) Shuffle() Any {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

// Keys returns a slice of strings representing the indexes of this slice
func (x Any) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (x Any) GetIndex(index int) (any, bool) {
	return slice.AtOK(x, index)
}

func (x Any) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Any) GetAnyOK(key string) (any, bool) {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		return x[index], true
	}

	switch key {

	case "last":
		return x.GetIndex(x.Length() - 1)

	case "next":
		return x.GetIndex(x.Length())
	}

	return nil, false
}

func (x Any) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Any) GetBoolOK(key string) (bool, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.BoolOk(value, false)
	}
	return false, false
}

func (x Any) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Any) GetIntOK(key string) (int, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.IntOk(value, 0)
	}

	return 0, false
}

func (x Any) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Any) GetInt64OK(key string) (int64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.Int64Ok(value, 0)
	}
	return 0, false
}

func (x Any) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Any) GetFloatOK(key string) (float64, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.FloatOk(value, 0)
	}
	return 0, false
}

func (x *Any) GetPointer(key string) (any, bool) {
	if index, ok := sliceStringIndex(key); ok {
		growSlice(x, index)
		return pointer.To((*x)[index]), true
	}

	if key == "last" {
		return x.GetPointer(strconv.Itoa(x.Length() - 1))
	}

	return nil, false
}

func (x Any) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x Any) GetStringOK(key string) (string, bool) {
	if value, ok := x.GetAnyOK(key); ok {
		return convert.StringOk(value, "")
	}
	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (x *Any) SetIndex(index int, value any) bool {
	growSlice(x, index)
	(*x)[index] = value
	return true
}

func (x *Any) SetAny(key string, value any) bool {

	if index, ok := sliceStringIndex(key); ok {
		growSlice(x, index)
		(*x)[index] = value
		return true
	}

	return false
}

func (x *Any) SetBool(key string, value bool) bool {
	return x.SetAny(key, value)
}

// SetInt sets a property with an int value
func (x *Any) SetInt(key string, value int) bool {
	return x.SetAny(key, value)
}

// SetInt64 sets a property with an int64 value
func (x *Any) SetInt64(key string, value int64) bool {
	return x.SetAny(key, value)
}

// SetFloat sets a property with a float64 value
func (x *Any) SetFloat(key string, value float64) bool {
	return x.SetAny(key, value)
}

// SetString sets a property with an string value
func (x *Any) SetString(key string, value string) bool {
	return x.SetAny(key, value)
}

// SetValue sets the entire value of this slice to the provided value
func (x *Any) SetValue(value any) error {
	*x = convert.SliceOfAny(value)
	return nil
}

// Remove removes the element with the specified key
func (x *Any) Remove(key string) bool {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}

// RemoveAt removes the element at the specified index
func (x *Any) RemoveAt(index int) bool {

	if index, ok := sliceIndex(index, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
