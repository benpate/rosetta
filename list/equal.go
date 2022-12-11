package list

import "strings"

const DelimiterEqual = '='

type Equal string

func ByEqual(value ...string) List {
	return Equal(strings.Join(value, string(DelimiterEqual)))
}

func (list Equal) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Equal) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterEqual)
}

func (list Equal) Head() string {
	return Head(list, DelimiterEqual)
}

func (list Equal) Tail() List {
	return Tail(list, DelimiterEqual)
}

func (list Equal) First() string {
	return Head(list, DelimiterEqual)
}

func (list Equal) Last() string {
	return Last(list, DelimiterEqual)
}

func (list Equal) RemoveLast() List {
	return RemoveLast(list, DelimiterEqual)
}

func (list Equal) Split() (string, List) {
	return Split(list, DelimiterEqual)
}

func (list Equal) SplitTail() (List, string) {
	return SplitTail(list, DelimiterEqual)
}

func (list Equal) At(index int) string {
	return At(list, DelimiterEqual, index)
}

func (list Equal) PushHead(value string) List {
	return PushHead(list, value, DelimiterEqual)
}

func (list Equal) PushTail(value string) List {
	return PushTail(list, value, DelimiterEqual)
}

func (list Equal) String() string {
	return string(list)
}
