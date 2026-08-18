package list

// Stringlike constrains the underlying string-like types a list may be built from.
type Stringlike interface {
	~string | []byte
}

// List interface wraps all of the list manipulation methods implemented by standard lists in this library.
type List interface {

	// IsEmpty returns TRUE if the list contains no items
	IsEmpty() bool

	// IsEmptyTail returns TRUE if the list contains at most one item, so its tail is empty
	IsEmptyTail() bool

	// Head returns the first item in the list
	Head() string

	// Tail returns a list containing every item except the first
	Tail() List

	// First returns the first item in the list
	First() string

	// Last returns the last item in the list
	Last() string

	// RemoveLast returns a list containing every item except the last
	RemoveLast() List

	// Split returns the first item in the list, along with a list of the remaining items
	Split() (string, List)

	// SplitTail returns a list of every item except the last, along with the last item
	SplitTail() (List, string)

	// At returns the item at the specified index, or an empty string if the index is out of bounds
	At(index int) string

	// PushHead returns a new list with the value prepended to the front
	PushHead(value string) List

	// PushTail returns a new list with the value appended to the end
	PushTail(value string) List

	// String returns the delimited string representation of the list
	String() string

	// Bytes returns the delimited byte-slice representation of the list
	Bytes() []byte
}
