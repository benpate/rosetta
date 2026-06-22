package schema

import (
	"math"

	"golang.org/x/exp/constraints"
)

// getLength returns the length of an object, if it is an ArrayGetter
func getLength(object any) (int, bool) {

	if getter, ok := object.(LengthGetter); ok {
		return getter.Length(), true
	}

	return 0, false
}

// getIndex returns the value at a specific index, if the object is an ArrayGetter
func getIndex(object any, index int) (any, bool) {

	if getter, ok := object.(ArrayGetter); ok {
		return getter.GetIndex(index)
	}

	return nil, false
}

// isMultipleOfInteger reports whether value is an exact integer multiple of
// multipleOf, using integer modulo so that large values are not corrupted by a
// detour through float64. A multipleOf of zero is treated as "no constraint"
// (and also avoids a divide-by-zero panic).
func isMultipleOfInteger[T constraints.Integer](value, multipleOf T) bool {
	if multipleOf == 0 {
		return true
	}
	return value%multipleOf == 0
}

// notMultipleOfInteger returns TRUE when the value is not an exact integer multiple of multipleOf.
func notMultipleOfInteger[T constraints.Integer](value, multipleOf T) bool {
	return !isMultipleOfInteger(value, multipleOf)
}

// isMultipleOfFloat reports whether value is an exact multiple of multipleOf.
// A multipleOf of zero is treated as "no constraint".
func isMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {
	if multipleOf == 0 {
		return true
	}
	return math.Mod(float64(value), float64(multipleOf)) == 0
}

// notMultipleOfFloat returns TRUE when the value is not an exact multiple of multipleOf.
func notMultipleOfFloat[T constraints.Float](value, multipleOf T) bool {
	return !isMultipleOfFloat(value, multipleOf)
}
