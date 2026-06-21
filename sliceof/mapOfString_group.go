package sliceof

// MapOfStringGrouper reports group boundaries within a MapOfString slice, grouped by a named field.
type MapOfStringGrouper struct {
	slice MapOfString
	field string
}

// GroupBy returns a grouper that detects boundaries in this slice using the named field.
func (x MapOfString) GroupBy(field string) MapOfStringGrouper {
	return MapOfStringGrouper{
		slice: x,
		field: field,
	}
}

// IsHeader returns TRUE if the element at index begins a new group.
func (grouper MapOfStringGrouper) IsHeader(index int) bool {

	if index <= 0 {
		return true
	}

	if index >= grouper.slice.Length() {
		return false
	}

	previous, _ := grouper.slice[index-1].GetStringOK(grouper.field)
	current, _ := grouper.slice[index].GetStringOK(grouper.field)

	return previous != current
}

// IsFooter returns TRUE if the element at index ends a group.
func (grouper MapOfStringGrouper) IsFooter(index int) bool {

	if index <= 0 {
		return false
	}

	if index >= grouper.slice.Length()-1 {
		return true
	}

	current, _ := grouper.slice[index].GetStringOK(grouper.field)
	next, _ := grouper.slice[index+1].GetStringOK(grouper.field)

	return current != next
}
