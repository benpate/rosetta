package list

import "strings"

// DelimiterSemicolon is the delimiter used by the Semicolon list type.
const DelimiterSemicolon = ';'

// Semicolon is a List backed by a semicolon-delimited string.
type Semicolon string

// BySemicolon returns a List that joins/splits its items on semicolons.
func BySemicolon(value ...string) List {
	return Semicolon(strings.Join(value, string(DelimiterSemicolon)))
}

// IsEmpty returns TRUE if the list contains no items.
func (list Semicolon) IsEmpty() bool {
	return IsEmpty(list)
}

// IsEmptyTail returns TRUE if the list has no items after the head.
func (list Semicolon) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSemicolon)
}

// Head returns the first item in the list.
func (list Semicolon) Head() string {
	return Head(list, DelimiterSemicolon)
}

// Tail returns a new list containing every item after the head.
func (list Semicolon) Tail() List {
	return Tail(list, DelimiterSemicolon)
}

// First returns the first item in the list.
func (list Semicolon) First() string {
	return Head(list, DelimiterSemicolon)
}

// Last returns the final item in the list.
func (list Semicolon) Last() string {
	return Last(list, DelimiterSemicolon)
}

// RemoveLast returns a new list with the final item removed.
func (list Semicolon) RemoveLast() List {
	return RemoveLast(list, DelimiterSemicolon)
}

// Split returns the head item and a list of the remaining items.
func (list Semicolon) Split() (string, List) {
	return Split(list, DelimiterSemicolon)
}

// SplitTail returns a list of the leading items and the final item.
func (list Semicolon) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSemicolon)
}

// At returns the item at the given index, or "" if out of range.
func (list Semicolon) At(index int) string {
	return At(list, DelimiterSemicolon, index)
}

// PushHead returns a new list with the value prepended as the new head.
func (list Semicolon) PushHead(value string) List {
	return PushHead(list, value, DelimiterSemicolon)
}

// PushTail returns a new list with the value appended as the new tail.
func (list Semicolon) PushTail(value string) List {
	return PushTail(list, value, DelimiterSemicolon)
}

// String returns the list as its underlying delimited string.
func (list Semicolon) String() string {
	return string(list)
}

// Bytes returns the list as its underlying delimited byte slice.
func (list Semicolon) Bytes() []byte {
	return []byte(list)
}
