package sliceof

// ObjectGrouper calculates group headers and footers for a sliceof.Object
type ObjectGrouper[T any] struct {
	slice Object[T]
	field string
}

// GroupBy returns an Grouper object for this slice
func (x Object[T]) GroupBy(field string) ObjectGrouper[T] {
	return ObjectGrouper[T]{
		slice: x,
		field: field,
	}
}

// IsHeader returns TRUE if the record at the given index is the FIRST item in a group.
func (grouper ObjectGrouper[T]) IsHeader(index int) bool {

	if index <= 0 {
		return true
	}

	if index >= grouper.slice.Length() {
		return false
	}

	if previousGetter, isGetter := any(grouper.slice[index-1]).(stringOKGetter); isGetter {

		if currentGetter, isGetter := any(grouper.slice[index]).(stringOKGetter); isGetter {

			previous, _ := previousGetter.GetStringOK(grouper.field)
			current, _ := currentGetter.GetStringOK(grouper.field)

			return previous != current
		}
	}

	return false
}

// IsFooter returns TRUE if the record at the given index is the LAST item in a group
func (grouper ObjectGrouper[T]) IsFooter(index int) bool {

	if index <= 0 {
		return false
	}

	if index >= grouper.slice.Length()-1 {
		return true
	}

	if currentGetter, isGetter := any(grouper.slice[index]).(stringOKGetter); isGetter {

		if nextGetter, isGetter := any(grouper.slice[index+1]).(stringOKGetter); isGetter {

			current, _ := currentGetter.GetStringOK(grouper.field)
			next, _ := nextGetter.GetStringOK(grouper.field)

			return current != next
		}
	}

	return false
}
