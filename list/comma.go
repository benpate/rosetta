package list

import "strings"

// DelimiterComma is the delimiter used by the Comma list type.
const DelimiterComma = ','

// Comma is a List backed by a comma-delimited string.
type Comma string

// ByComma returns a List that joins/splits its items on commas.
func ByComma(value ...string) List {
	return Comma(strings.Join(value, string(DelimiterComma)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Comma) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Comma) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterComma)
}

// Head returns the first item in the list.
func (list Comma) Head() string {
	return Head(list, DelimiterComma)
}

// Tail returns a new list containing every item after the head.
func (list Comma) Tail() List {
	return Tail(list, DelimiterComma)
}

// First returns the first item in the list.
func (list Comma) First() string {
	return Head(list, DelimiterComma)
}

// Last returns the final item in the list.
func (list Comma) Last() string {
	return Last(list, DelimiterComma)
}

// RemoveLast returns a new list with the final item removed.
func (list Comma) RemoveLast() List {
	return RemoveLast(list, DelimiterComma)
}

// Split returns the head item and a list of the remaining items.
func (list Comma) Split() (string, List) {
	return Split(list, DelimiterComma)
}

// SplitTail returns a list of the leading items and the final item.
func (list Comma) SplitTail() (List, string) {
	return SplitTail(list, DelimiterComma)
}

// At returns the item at the given index, or "" if out of range.
func (list Comma) At(index int) string {
	return At(list, DelimiterComma, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Comma) PushHead(value string) List {
	return PushHead(list, value, DelimiterComma)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Comma) PushTail(value string) List {
	return PushTail(list, value, DelimiterComma)
}

// String returns the list as its underlying delimited string.
func (list Comma) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Comma) Bytes() []byte {
	return []byte(list)
}
