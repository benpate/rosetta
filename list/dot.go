package list

import "strings"

const DelimiterDot = '.'

type Dot string

func ByDot(value ...string) List {
	return Dot(strings.Join(value, string(DelimiterDot)))
}

func (list Dot) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Dot) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterDot)
}

func (list Dot) Head() string {
	return Head(list, DelimiterDot)
}

func (list Dot) Tail() List {
	return Tail(list, DelimiterDot)
}

func (list Dot) Last() string {
	return Last(list, DelimiterDot)
}

func (list Dot) RemoveLast() List {
	return RemoveLast(list, DelimiterDot)
}

func (list Dot) Split() (string, List) {
	return Split(list, DelimiterDot)
}

func (list Dot) SplitTail() (List, string) {
	return SplitTail(list, DelimiterDot)
}

func (list Dot) At(index int) string {
	return At(list, DelimiterDot, index)
}

func (list Dot) PushHead(value string) List {
	return PushHead(list, value, DelimiterDot)
}

func (list Dot) PushTail(value string) List {
	return PushTail(list, value, DelimiterDot)
}

func (list Dot) String() string {
	return string(list)
}
