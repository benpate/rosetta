package null

import (
	"encoding/json"
	"reflect"

	"github.com/benpate/derp"
)

// Object provides a nullable wrapper around any type T
type Object[T any] struct {
	value   T
	present bool
}

// NewObject returns a fully populated, nullable Object
func NewObject[T any](value T) Object[T] {
	return Object[T]{
		value:   value,
		present: true,
	}
}

// Object returns the actual value of this object
func (x Object[T]) Object() T {
	return x.value
}

// Set applies a new value to the nullable item
func (x *Object[T]) Set(value T) {
	x.value = value
	x.present = true
}

// Unset removes the value from this item, and sets it to null
func (x *Object[T]) Unset() {
	var empty T
	x.value = empty
	x.present = false
}

// IsNull returns TRUE if this value is null
func (x Object[T]) IsNull() bool {
	return !x.present
}

// IsNil returns TRUE if this value is null.  It is an alias for IsNull
func (x Object[T]) IsNil() bool {
	return x.IsNull()
}

// IsZero returns TRUE if this value is null, or contains the zero value for its data type
func (x Object[T]) IsZero() bool {

	// A null value is always zero
	if x.IsNull() {
		return true
	}

	// An untyped nil (when T is an interface) has no reflected value at all
	value := reflect.ValueOf(x.value)

	if !value.IsValid() {
		return true
	}

	return value.IsZero()
}

// IsPresent returns TRUE if this value is present
func (x Object[T]) IsPresent() bool {
	return x.present
}

// Interface returns the underlying value (if present) or NIL
func (x Object[T]) Interface() any {

	if x.present {
		return x.value
	}

	return nil
}

// MarshalJSON implements the json.Marshaller interface
func (x Object[T]) MarshalJSON() ([]byte, error) {

	if !x.present {
		return []byte("null"), nil
	}

	result, err := json.Marshal(x.value)

	if err != nil {
		return nil, derp.Wrap(err, "null.Object.MarshalJSON", "Unable to marshal value")
	}

	return result, nil
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (x *Object[T]) UnmarshalJSON(value []byte) error {

	valueStr := string(value)

	// Allow null values to be null.  A present-but-nil pointer T also
	// marshals to "null", so it reads back as unset, not as a nil value.
	if (valueStr == "") || (valueStr == "null") {
		x.Unset()
		return nil
	}

	// Try to unmarshal the value into the underlying type
	var result T

	if err := json.Unmarshal(value, &result); err != nil {
		return derp.Wrap(err, "null.Object.UnmarshalJSON", "Invalid value", valueStr)
	}

	x.Set(result)

	// The Object abides.
	return nil
}
