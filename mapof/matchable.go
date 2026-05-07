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

func (m Matchable[T]) Keys() []string {
	return maps.KeysSorted(m)
}

func (m Matchable[T]) Values() []T {
	result := make([]T, 0, len(m))
	for _, value := range m {
		result = append(result, value)
	}
	return result
}

func (m Matchable[T]) Equal(value map[string]any) bool {
	return reflect.DeepEqual(m, Any(value))
}

func (m Matchable[T]) NotEqual(value map[string]any) bool {
	return !reflect.DeepEqual(m, Any(value))
}

func (m Matchable[T]) IsEmpty() bool {
	return len(m) == 0
}

func (m Matchable[T]) NotEmpty() bool {
	return len(m) > 0
}

/****************************************
 * Tree Traversal
 ****************************************/

func (m Matchable[T]) GetPointer(key string) (any, bool) {
	result, ok := m[key]
	return result, ok
}

func (m *Matchable[T]) makeNotNil() {
	if *m == nil {
		*m = make(Matchable[T])
	}
}

func (m *Matchable[T]) Remove(key string) bool {
	m.makeNotNil()
	delete(*m, key)
	return true
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

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
