package list

import "strings"

// DelimiterSlash is the delimiter used by the Slash list type.
const DelimiterSlash = '/'

// Slash is a List backed by a slash-delimited string.
type Slash string

// BySlash returns a List that joins/splits its items on slashes.
func BySlash(value ...string) List {
	return Slash(strings.Join(value, string(DelimiterSlash)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Slash) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Slash) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSlash)
}

// Head returns the first item in the list.
func (list Slash) Head() string {
	return Head(list, DelimiterSlash)
}

// Tail returns a new list containing every item after the head.
func (list Slash) Tail() List {
	return Tail(list, DelimiterSlash)
}

// First returns the first item in the list.
func (list Slash) First() string {
	return Head(list, DelimiterSlash)
}

// Last returns the final item in the list.
func (list Slash) Last() string {
	return Last(list, DelimiterSlash)
}

// RemoveLast returns a new list with the final item removed.
func (list Slash) RemoveLast() List {
	return RemoveLast(list, DelimiterSlash)
}

// Split returns the head item and a list of the remaining items.
func (list Slash) Split() (string, List) {
	return Split(list, DelimiterSlash)
}

// SplitTail returns a list of the leading items and the final item.
func (list Slash) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSlash)
}

// At returns the item at the given index, or "" if out of range.
func (list Slash) At(index int) string {
	return At(list, DelimiterSlash, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Slash) PushHead(value string) List {
	return PushHead(list, value, DelimiterSlash)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Slash) PushTail(value string) List {
	return PushTail(list, value, DelimiterSlash)
}

// String returns the list as its underlying delimited string.
func (list Slash) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Slash) Bytes() []byte {
	return []byte(list)
}
