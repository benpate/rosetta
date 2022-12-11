package list

import "strings"

const DelimiterSpace = ' '

type Space string

func BySpace(value ...string) List {
	return Space(strings.Join(value, string(DelimiterSpace)))
}

func (list Space) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Space) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSpace)
}

func (list Space) Head() string {
	return Head(list, DelimiterSpace)
}

func (list Space) Tail() List {
	return Tail(list, DelimiterSpace)
}

func (list Space) Last() string {
	return Last(list, DelimiterSpace)
}

func (list Space) RemoveLast() List {
	return RemoveLast(list, DelimiterSpace)
}

func (list Space) Split() (string, List) {
	return Split(list, DelimiterSpace)
}

func (list Space) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSpace)
}

func (list Space) At(index int) string {
	return At(list, DelimiterSpace, index)
}

func (list Space) PushHead(value string) List {
	return PushHead(list, value, DelimiterSpace)
}

func (list Space) PushTail(value string) List {
	return PushTail(list, value, DelimiterSpace)
}

func (list Space) String() string {
	return string(list)
}
