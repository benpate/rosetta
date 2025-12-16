package sliceof

import (
	"iter"
	"math/rand"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
)

type MapOfAny []mapof.Any

func NewMapOfAny(values ...mapof.Any) MapOfAny {

	if len(values) == 0 {
		return make(MapOfAny, 0)
	}

	return MapOfAny(values)
}

/******************************************
 * Slice Manipulations
 ******************************************/

// Length returns the number of elements in the slice
func (x MapOfAny) Length() int {
	return len(x)
}

// IsLength returns TRUE if the slice contains exactly "length" elements
func (x MapOfAny) IsLength(length int) bool {
	return len(x) == length
}

// IsZero returns TRUE if the slice contains no elements.
// This is an alias for IsEmpty, and implements the `Zeroer`
// interface used by many packages (including go/json)
func (x MapOfAny) IsZero() bool {
	return len(x) == 0
}

// IsEmpty returns TRUE if the slice contains no elements
func (x MapOfAny) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x MapOfAny) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x MapOfAny) First() mapof.Any {
	if len(x) > 0 {
		return x[0]
	}
	return mapof.NewAny()
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x MapOfAny) FirstN(n int) MapOfAny {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x MapOfAny) Last() mapof.Any {
	if len(x) > 0 {
		return x[len(x)-1]
	}

	return mapof.NewAny()
}

// Find returns the first element in the slice that satisfies the provided function.
func (x MapOfAny) Find(fn func(mapof.Any) bool) (mapof.Any, bool) {
	return slice.Find(x, fn)
}

// Filter returns all elements in the slice that satisfies the provided function.
func (x MapOfAny) Filter(fn func(mapof.Any) bool) MapOfAny {
	return slice.Filter(x, fn)
}

// At returns a bound-safe element from the slice.  If the index
// is out of bounds, then `At` returns the zero value for the slice type
func (x MapOfAny) At(index int) mapof.Any {
	return slice.At(x, index)
}

func (x MapOfAny) AtOK(index int) (mapof.Any, bool) {
	return slice.AtOK(x, index)
}

// Reverse returns a new slice with the elements in reverse order
func (x MapOfAny) Reverse() MapOfAny {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the "match" function returns TRUE for any element in the slice
func (x MapOfAny) Contains(match func(mapof.Any) bool) bool {
	for _, value := range x {
		if match(value) {
			return true
		}
	}

	return false
}

// Range returns an iterator that yields each value in this slice.
func (x MapOfAny) Range() iter.Seq2[int, mapof.Any] {
	return slice.Range(x)
}

// Append adds one or more elements to the end of the slice
func (x *MapOfAny) Append(values ...mapof.Any) {
	*x = append(*x, values...)
}

// Shuffle randomizes the order of the elements in the slice
func (x MapOfAny) Shuffle() MapOfAny {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return x
}

// Keys returns a slice of strings representing the indexes of this slice
func (x MapOfAny) Keys() []string {
	keys := make([]string, len(x))

	for i := range x {
		keys[i] = strconv.Itoa(i)
	}

	return keys
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (x MapOfAny) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x MapOfAny) GetAnyOK(key string) (any, bool) {
	if index, ok := sliceStringIndex(key, x.Length()); ok {
		return x[index], true
	}

	switch key {

	case "last":
		return x.AtOK(x.Length() - 1)
	case "next":
		return x.AtOK(x.Length())
	}

	return nil, false
}

func (x *MapOfAny) GetPointer(name string) (any, bool) {

	// Get a valid index for the slice
	if index, ok := sliceStringIndex(name); ok {
		growSlice(x, index)

		// Return result
		return &(*x)[index], true
	}

	switch name {

	case "last":
		return x.GetPointer(strconv.Itoa(x.Length() - 1))
	case "next":
		return x.GetPointer(strconv.Itoa(x.Length()))
	}

	// Failure!!
	return nil, false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (s *MapOfAny) SetIndex(index int, value any) bool {

	mapValue := convert.MapOfAny(value)

	growSlice(s, index)
	(*s)[index] = mapValue
	return true
}

func (s *MapOfAny) SetValue(value any) error {

	if typed, ok := value.(MapOfAny); ok {
		*s = typed
		return nil
	}

	return derp.InternalError("sliceof.Map.SetValue", "Unable to convert value to Map", value)
}

func (x *MapOfAny) Remove(key string) bool {

	if index, ok := sliceStringIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}

func (x *MapOfAny) RemoveAt(index int) bool {

	if index, ok := sliceIndex(index, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
