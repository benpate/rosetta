package list

import "strings"

// DelimiterSpace is the delimiter used by the Space list type.
const DelimiterSpace = ' '

// Space is a List backed by a space-delimited string.
type Space string

// BySpace returns a List that joins/splits its items on spaces.
func BySpace(value ...string) List {
	return Space(strings.Join(value, string(DelimiterSpace)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Space) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Space) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSpace)
}

// Head returns the first item in the list.
func (list Space) Head() string {
	return Head(list, DelimiterSpace)
}

// Tail returns a new list containing every item after the head.
func (list Space) Tail() List {
	return Tail(list, DelimiterSpace)
}

// First returns the first item in the list.
func (list Space) First() string {
	return Head(list, DelimiterSpace)
}

// Last returns the final item in the list.
func (list Space) Last() string {
	return Last(list, DelimiterSpace)
}

// RemoveLast returns a new list with the final item removed.
func (list Space) RemoveLast() List {
	return RemoveLast(list, DelimiterSpace)
}

// Split returns the head item and a list of the remaining items.
func (list Space) Split() (string, List) {
	return Split(list, DelimiterSpace)
}

// SplitTail returns a list of the leading items and the final item.
func (list Space) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSpace)
}

// At returns the item at the given index, or "" if out of range.
func (list Space) At(index int) string {
	return At(list, DelimiterSpace, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Space) PushHead(value string) List {
	return PushHead(list, value, DelimiterSpace)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Space) PushTail(value string) List {
	return PushTail(list, value, DelimiterSpace)
}

// String returns the list as its underlying delimited string.
func (list Space) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Space) Bytes() []byte {
	return []byte(list)
}
