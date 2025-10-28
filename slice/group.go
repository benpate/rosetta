package slice

type Grouper[T stringOKGetter] struct {
	slice []T
	field string
}

func NewGrouper[T stringOKGetter](slice []T, field string) Grouper[T] {
	return Grouper[T]{
		slice: slice,
		field: field,
	}
}

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
