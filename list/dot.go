package list

import "strings"

const DelimiterDot = '.'

type Dot []byte

func ByDot(value ...string) List {
	return List(Dot(strings.Join(value, string(DelimiterDot))))
}

func (list Dot) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Dot) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterDot)
}

func (list Dot) Head() string {
	return string(Head(list, DelimiterDot))
}

func (list Dot) Tail() List {
	return Dot(Tail(list, DelimiterDot))
}

func (list Dot) Last() string {
	return string(Last(list, DelimiterDot))
}

func (list Dot) RemoveLast() List {
	return Dot(RemoveLast(list, DelimiterDot))
}

func (list Dot) Split() (string, List) {
	head, tail := Split(list, DelimiterDot)
	return string(head), Dot(tail)
}

func (list Dot) SplitTail() (List, string) {
	head, tail := SplitTail(list, DelimiterDot)
	return Dot(head), string(tail)
}

func (list Dot) At(index int) string {
	return string(At(list, DelimiterDot, index))
}

func (list Dot) PushHead(value string) List {
	return Dot(PushHead(list, []byte(value), DelimiterDot))
}

func (list Dot) PushTail(value string) List {
	return Dot(PushTail(list, []byte(value), DelimiterDot))
}

func (list Dot) String() string {
	return string(list)
}
