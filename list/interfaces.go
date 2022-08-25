package list

type Stringlike interface {
	~string | []byte
}

// List interface wraps all of the list manipulation methods implemented by standard lists in this library.
type List interface {
	IsEmpty() bool
	IsEmptyTail() bool
	Head() string
	Tail() List
	Last() string
	RemoveLast() List
	Split() (string, List)
	SplitTail() (List, string)
	At(index int) string
	PushHead(value string) List
	PushTail(value string) List
	String() string
}
