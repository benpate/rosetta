package list

import "strings"

// DelimiterEqual is the delimiter used by the Equal list type.
const DelimiterEqual = '='

// Equal is a List backed by an equals-delimited string.
type Equal string

// ByEqual returns a List that joins/splits its items on equals signs.
func ByEqual(value ...string) List {
	return Equal(strings.Join(value, string(DelimiterEqual)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Equal) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Equal) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterEqual)
}

// Head returns the first item in the list.
func (list Equal) Head() string {
	return Head(list, DelimiterEqual)
}

// Tail returns a new list containing every item after the head.
func (list Equal) Tail() List {
	return Tail(list, DelimiterEqual)
}

// First returns the first item in the list.
func (list Equal) First() string {
	return Head(list, DelimiterEqual)
}

// Last returns the final item in the list.
func (list Equal) Last() string {
	return Last(list, DelimiterEqual)
}

// RemoveLast returns a new list with the final item removed.
func (list Equal) RemoveLast() List {
	return RemoveLast(list, DelimiterEqual)
}

// Split returns the head item and a list of the remaining items.
func (list Equal) Split() (string, List) {
	return Split(list, DelimiterEqual)
}

// SplitTail returns a list of the leading items and the final item.
func (list Equal) SplitTail() (List, string) {
	return SplitTail(list, DelimiterEqual)
}

// At returns the item at the given index, or "" if out of range.
func (list Equal) At(index int) string {
	return At(list, DelimiterEqual, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Equal) PushHead(value string) List {
	return PushHead(list, value, DelimiterEqual)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Equal) PushTail(value string) List {
	return PushTail(list, value, DelimiterEqual)
}

// String returns the list as its underlying delimited string.
func (list Equal) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Equal) Bytes() []byte {
	return []byte(list)
}
