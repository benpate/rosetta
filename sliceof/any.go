package sliceof

import (
	"math/rand"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
)

type Any []any

func NewAny() Any {
	return make(Any, 0)
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

// IsEmpty returns TRUE if the slice contains no elements
func (x Any) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the slice contains at least one element
func (x Any) NotEmpty() bool {
	return len(x) > 0
}

// First returns the first element in the slice, or nil if the slice is empty
func (x Any) First() any {
	if len(x) > 0 {
		return x[0]
	}
	return nil
}

// FirstN returns the first "n" elements in the slice, or all elements if "n" is greater than the length of the slice
func (x Any) FirstN(n int) Any {
	if n > len(x) {
		n = len(x)
	}

	return x[:n]
}

// Last returns the last element in the slice, or nil if the slice is empty
func (x Any) Last() any {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return nil
}

// Reverse returns a new slice with the elements in reverse order
func (x Any) Reverse() Any {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	return x
}

// Contains returns TRUE if the slice contains the specified value
func (x Any) Contains(value any) bool {
	return slice.Contains(x, value)
}

// ContainsAny returns TRUE if the slice contains any of the specified values
func (x Any) ContainsAny(values ...any) bool {
	return slice.ContainsAny(x, values...)
}

// ContainsAll returns TRUE if the slice contains all of the specified values
func (x Any) ContainsAll(values ...any) bool {
	return slice.ContainsAll(x, values...)
}

// Equal returns TRUE if the slice contains exactly the same elements as the specified value
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

/******************************************
 * Getter Interfaces
 ******************************************/

func (x Any) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Any) GetAnyOK(key string) (any, bool) {
	if index, ok := sliceIndex(key, x.Length()); ok {
		return x[index], true
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
	if index, ok := sliceIndex(key); ok {
		growSlice(x, index)
		return &(*x)[index], true
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

func (x *Any) SetAny(key string, value any) bool {

	if index, ok := sliceIndex(key); ok {
		growSlice(x, index)
		(*x)[index] = value
		return true
	}

	return false
}

func (x *Any) SetBool(key string, value bool) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetInt(key string, value int) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetInt64(key string, value int64) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetFloat(key string, value float64) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetString(key string, value string) bool {
	return x.SetAny(key, value)
}

func (x *Any) SetValue(value any) error {
	*x = convert.SliceOfAny(value)
	return nil
}

func (x *Any) Remove(key string) bool {

	if index, ok := sliceIndex(key, x.Length()); ok {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}

	return false
}
