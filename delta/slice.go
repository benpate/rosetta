// Package delta provides collection types that track their own additions and deletions
package delta

import (
	"encoding/json"
	"slices"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
)

// Slice tracks changes to a slice of generic values.  As items are added or removed from the slice,
// they are tracked in the "Added" and "Deleted" lists.
type Slice[T comparable] struct {
	Values  []T `json:"values" bson:"values"`
	Added   []T `json:"-"      bson:"-"`
	Deleted []T `json:"-"      bson:"-"`
}

// NewSlice returns a fully initialized Slice object
func NewSlice[T comparable](values ...T) Slice[T] {

	if values == nil {
		values = make([]T, 0)
	}

	return Slice[T]{
		Values:  values,
		Added:   make([]T, 0),
		Deleted: make([]T, 0),
	}
}

// Length returns the number of keys in the map
func (s Slice[T]) Length() int {
	return len(s.Values)
}

// IsChanged returns TRUE if the slice has been modified
func (s Slice[T]) IsChanged() bool {

	if len(s.Added) > 0 {
		return true
	}

	if len(s.Deleted) > 0 {
		return true
	}

	return false
}

// Unchanged returns the values that have not been added or removed
func (s Slice[T]) Unchanged() []T {
	return slice.Difference(s.Values, s.Added)
}

// Reset resets all of the added/deleted values
func (s *Slice[T]) Reset() {
	s.Added = make([]T, 0)
	s.Deleted = make([]T, 0)
}

/******************************************
 * Schema Interfaces
 ******************************************/

// GetValue implements schema.ValueGetter and returns the current value
func (s Slice[T]) GetValue() any {
	return s.Values
}

// GetIndex implements schema.ArrayGetter, returning the value at the provided index
// and TRUE, or (nil, false) if the index is out of range.
func (s Slice[T]) GetIndex(index int) (any, bool) {

	if (index < 0) || (index >= len(s.Values)) {
		return nil, false
	}

	return s.Values[index], true
}

// SetIndex implements schema.ArraySetter, storing a value at the provided index (growing the
// slice to fit if necessary) and tracking any newly-added value. Returns FALSE if the value
// is not assignable to the slice's element type.
func (s *Slice[T]) SetIndex(index int, value any) bool {

	typed, ok := value.(T)

	if !ok {
		return false
	}

	if index < 0 {
		return false
	}

	// Grow the slice to fit, recording each appended slot as a new value.
	var zero T
	for index >= len(s.Values) {
		s.Values = append(s.Values, zero)
	}

	// Record the new value in the "added" list unless it was already present.
	if slices.Index(s.Values, typed) == -1 {
		s.Added = append(s.Added, typed)
	}

	s.Values[index] = typed
	return true
}

// SetValue implements schema.ValueSetter and updates the current value,
// tracking changes to the added and deleted lists.
func (s *Slice[T]) SetValue(value any) error {

	newValues, ok := value.([]T)

	if !ok {
		newValues = make([]T, 0)
	}

	// Reset added/deleted lists
	s.Added = make([]T, 0)
	s.Deleted = s.Values
	s.Values = newValues

	// Find all values that are in the new list, but not in the existing list
	for _, newValue := range s.Values {

		// If the new value is found in the list of "old" values,
		// then remove it from the "deleted" list
		if index := slices.Index(s.Deleted, newValue); index > -1 {
			s.Deleted = append(s.Deleted[:index], s.Deleted[index+1:]...)
			continue
		}

		// Otherise, add the new value to the "added" list
		s.Added = append(s.Added, newValue)
	}

	return nil
}

/******************************************
 * JSON Serialization
 ******************************************/

// MarshalJSON implements the json.Marshaler interface
// and serializes the Slice into a JSON array
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

// UnmarshalJSON implements the json.Unmarshaler interface
// and deserializes the Slice from a JSON array
func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.Values); err != nil {
		return derp.Wrap(err, "delta.Slice.UnmarshalJSON", "Error unmarshalling JSON data", data)
	}

	s.Reset()
	return nil
}
