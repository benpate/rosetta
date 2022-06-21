package list

const DelimiterSpace = ' '

type Space []byte

func (list Space) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Space) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSpace)
}

func (list Space) Head() string {
	return string(Head(list, DelimiterSpace))
}

func (list Space) Tail() Space {
	return Space(Tail(list, DelimiterSpace))
}

func (list Space) Last() string {
	return string(Last(list, DelimiterSpace))
}

func (list Space) RemoveLast() Space {
	return Space(RemoveLast(list, DelimiterSpace))
}

func (list Space) Split() (string, Space) {
	head, tail := Split(list, DelimiterSpace)
	return string(head), Space(tail)
}

func (list Space) SplitTail() (Space, string) {
	head, tail := SplitTail(list, DelimiterSpace)
	return Space(head), string(tail)
}

func (list Space) At(index int) string {
	return string(At(list, DelimiterSpace, index))
}

func (list Space) PushHead(value string) Space {
	return Space(PushHead(list, []byte(value), DelimiterSpace))
}

func (list Space) PushTail(value string) Space {
	return Space(PushTail(list, []byte(value), DelimiterSpace))
}

func (list Space) String() string {
	return string(list)
}
