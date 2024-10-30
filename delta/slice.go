package delta

import (
	"encoding/json"
	"slices"

	"github.com/benpate/derp"
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
