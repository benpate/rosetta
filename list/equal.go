package list

const DelimiterEqual = '='

type Equal []byte

func ByEqual(value string) List {
	return List(Equal(value))
}

func (list Equal) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Equal) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterEqual)
}

func (list Equal) Head() string {
	return string(Head(list, DelimiterEqual))
}

func (list Equal) Tail() List {
	return Equal(Tail(list, DelimiterEqual))
}

func (list Equal) Last() string {
	return string(Last(list, DelimiterEqual))
}

func (list Equal) RemoveLast() List {
	return Equal(RemoveLast(list, DelimiterEqual))
}

func (list Equal) Split() (string, List) {
	head, tail := Split(list, DelimiterEqual)
	return string(head), Equal(tail)
}

func (list Equal) SplitTail() (List, string) {
	head, tail := SplitTail(list, DelimiterEqual)
	return Equal(head), string(tail)
}

func (list Equal) At(index int) string {
	return string(At(list, DelimiterEqual, index))
}

func (list Equal) PushHead(value string) List {
	return Equal(PushHead(list, []byte(value), DelimiterEqual))
}

func (list Equal) PushTail(value string) List {
	return Equal(PushTail(list, []byte(value), DelimiterEqual))
}

func (list Equal) String() string {
	return string(list)
}
