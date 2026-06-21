package slice

// Grouper reports group boundaries within a slice that is already sorted by a named string field.
type Grouper[T stringOKGetter] struct {
	slice []T
	field string
}

// NewGrouper returns a Grouper that detects boundaries in the slice using the named field.
func NewGrouper[T stringOKGetter](slice []T, field string) Grouper[T] {
	return Grouper[T]{
		slice: slice,
		field: field,
	}
}

// IsHeader returns TRUE if the element at index begins a new group (its field differs from the previous element).
func (grouper Grouper[T]) IsHeader(index int) bool {

	if index < 1 {
		return true
	}

	if index >= len(grouper.slice) {
		return false
	}

	currentValue, _ := grouper.slice[index].GetStringOK(grouper.field)
	previousValue, _ := grouper.slice[index-1].GetStringOK(grouper.field)
	return currentValue != previousValue
}

// IsFooter returns TRUE if the element at index ends a group (its field differs from the next element).
func (grouper Grouper[T]) IsFooter(index int) bool {

	if index < 0 {
		return false
	}

	if index >= len(grouper.slice)-1 {
		return true
	}

	currentValue, _ := grouper.slice[index].GetStringOK(grouper.field)
	nextValue, _ := grouper.slice[index+1].GetStringOK(grouper.field)
	return currentValue != nextValue
}
