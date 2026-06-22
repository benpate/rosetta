package list

import "strings"

// DelimiterDot is the delimiter used by the Dot list type.
const DelimiterDot = '.'

// Dot is a List backed by a dot-delimited string.
type Dot string

// ByDot returns a List that joins/splits its items on dots.
func ByDot(value ...string) List {
	return Dot(strings.Join(value, string(DelimiterDot)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Dot) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Dot) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterDot)
}

// Head returns the first item in the list.
func (list Dot) Head() string {
	return Head(list, DelimiterDot)
}

// Tail returns a new list containing every item after the head.
func (list Dot) Tail() List {
	return Tail(list, DelimiterDot)
}

// First returns the first item in the list.
func (list Dot) First() string {
	return Head(list, DelimiterDot)
}

// Last returns the final item in the list.
func (list Dot) Last() string {
	return Last(list, DelimiterDot)
}

// RemoveLast returns a new list with the final item removed.
func (list Dot) RemoveLast() List {
	return RemoveLast(list, DelimiterDot)
}

// Split returns the head item and a list of the remaining items.
func (list Dot) Split() (string, List) {
	return Split(list, DelimiterDot)
}

// SplitTail returns a list of the leading items and the final item.
func (list Dot) SplitTail() (List, string) {
	return SplitTail(list, DelimiterDot)
}

// At returns the item at the given index, or "" if out of range.
func (list Dot) At(index int) string {
	return At(list, DelimiterDot, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Dot) PushHead(value string) List {
	return PushHead(list, value, DelimiterDot)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Dot) PushTail(value string) List {
	return PushTail(list, value, DelimiterDot)
}

// String returns the list as its underlying delimited string.
func (list Dot) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Dot) Bytes() []byte {
	return []byte(list)
}
