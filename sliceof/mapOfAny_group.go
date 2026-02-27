package sliceof

type MapOfAnyGrouper struct {
	slice MapOfAny
	field string
}

func (x MapOfAny) GroupBy(field string) MapOfAnyGrouper {
	return MapOfAnyGrouper{
		slice: x,
		field: field,
	}
}

func (grouper MapOfAnyGrouper) IsHeader(index int) bool {

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

func (grouper MapOfAnyGrouper) IsFooter(index int) bool {

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
