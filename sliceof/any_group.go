package sliceof

// AnyGrouper reports group boundaries within an Any slice, grouped by a named field.
type AnyGrouper struct {
	slice Any
	field string
}

// GroupBy returns a grouper that detects boundaries in this slice using the named field.
func (x Any) GroupBy(field string) AnyGrouper {
	return AnyGrouper{
		slice: x,
		field: field,
	}
}

// IsHeader returns TRUE if the element at index begins a new group.
func (grouper AnyGrouper) IsHeader(index int) bool {

	if index <= 0 {
		return true
	}

	if index >= grouper.slice.Length() {
		return false
	}

	if previousGetter, isGetter := grouper.slice[index-1].(stringOKGetter); isGetter {

		if currentGetter, isGetter := grouper.slice[index].(stringOKGetter); isGetter {

			previous, _ := previousGetter.GetStringOK(grouper.field)
			current, _ := currentGetter.GetStringOK(grouper.field)

			return previous != current
		}
	}

	return false
}

// IsFooter returns TRUE if the element at index ends a group.
func (grouper AnyGrouper) IsFooter(index int) bool {

	if index <= 0 {
		return false
	}

	if index >= grouper.slice.Length()-1 {
		return true
	}

	if currentGetter, isGetter := grouper.slice[index].(stringOKGetter); isGetter {

		if nextGetter, isGetter := grouper.slice[index+1].(stringOKGetter); isGetter {

			current, _ := currentGetter.GetStringOK(grouper.field)
			next, _ := nextGetter.GetStringOK(grouper.field)

			return current != next
		}
	}

	return false
}
