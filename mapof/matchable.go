package mapof

import (
	"reflect"

	"github.com/benpate/exp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/maps"
)

// Matcher interface wraps the Match method, which allows objects to declare if they match a given predicate
type Matcher interface {
	Match(predicate exp.Predicate) bool
}

// Matchable is a map of objects that implement the Matcher interface and can be queried using the Match and MatchOne methods.
type Matchable[T Matcher] map[string]T

// NewMatchable returns a fully initialized Matchable object
func NewMatchable[T Matcher]() Matchable[T] {
	return make(map[string]T)
}

/******************************************
 * Query Interface
 ******************************************/

// Match returns a new map containing all elements that match the given expression.
func (m Matchable[T]) Match(expression exp.Expression) Matchable[T] {
	result := make(Matchable[T])

	for key, value := range m {
		if expression.Match(value.Match) {
			result[key] = value
		}
	}
	return result
}

// MatchOne returns the first element that matches the given expression, along with a boolean indicating if a match was found.
func (m Matchable[T]) MatchOne(expression exp.Expression) (T, bool) {
	for _, value := range m {
		if expression.Match(value.Match) {
			return value, true
		}
	}
	var empty T
	return empty, false
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (m Matchable[T]) Length() int {
	return len(m)
}

// Keys returns the map's keys in sorted order.
func (m Matchable[T]) Keys() []string {
	return maps.KeysSorted(m)
}

// Values returns all of the map's values (in unspecified order).
func (m Matchable[T]) Values() []T {
	result := make([]T, 0, len(m))
	for _, value := range m {
		result = append(result, value)
	}
	return result
}

// Equal returns TRUE if this map deeply equals the provided map.
func (m Matchable[T]) Equal(value map[string]any) bool {
	return reflect.DeepEqual(m, Any(value))
}

// NotEqual returns TRUE if this map does not deeply equal the provided map.
func (m Matchable[T]) NotEqual(value map[string]any) bool {
	return !reflect.DeepEqual(m, Any(value))
}

// IsEmpty returns TRUE if the map contains no elements.
func (m Matchable[T]) IsEmpty() bool {
	return len(m) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (m Matchable[T]) NotEmpty() bool {
	return len(m) > 0
}

/****************************************
 * Tree Traversal
 ****************************************/

// GetPointer returns the value for the key (implements the schema PointerGetter interface).
func (m Matchable[T]) GetPointer(key string) (any, bool) {
	result, ok := m[key]
	return result, ok
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (m *Matchable[T]) makeNotNil() {
	if *m == nil {
		*m = make(Matchable[T])
	}
}

// Remove deletes the key from the map.
func (m *Matchable[T]) Remove(key string) bool {
	m.makeNotNil()
	delete(*m, key)
	return true
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

// IsZeroValue returns TRUE if the named property is absent or holds a zero value.
func (m Matchable[T]) IsZeroValue(name string) bool {
	return compare.IsZero(m[name])
}

// MapOfAny implements the MapOfAny interface.
// It returns this value as a map[string]any
func (m Matchable[T]) MapOfAny() map[string]any {
	return convert.MapOfAny(m)
}

// MapOfString implements the MapOfString interface.
// It returns this value as a map[string]string
func (m Matchable[T]) MapOfString() map[string]string {
	return convert.MapOfString(m.MapOfAny())
}
