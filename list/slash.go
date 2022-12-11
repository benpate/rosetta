package list

import "strings"

const DelimiterSlash = '/'

type Slash string

func BySlash(value ...string) List {
	return Slash(strings.Join(value, string(DelimiterSlash)))
}

func (list Slash) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Slash) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSlash)
}

func (list Slash) Head() string {
	return Head(list, DelimiterSlash)
}

func (list Slash) Tail() List {
	return Tail(list, DelimiterSlash)
}

func (list Slash) First() string {
	return Head(list, DelimiterSlash)
}

func (list Slash) Last() string {
	return Last(list, DelimiterSlash)
}

func (list Slash) RemoveLast() List {
	return RemoveLast(list, DelimiterSlash)
}

func (list Slash) Split() (string, List) {
	return Split(list, DelimiterSlash)
}

func (list Slash) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSlash)
}

func (list Slash) At(index int) string {
	return At(list, DelimiterSlash, index)
}

func (list Slash) PushHead(value string) List {
	return PushHead(list, value, DelimiterSlash)
}

func (list Slash) PushTail(value string) List {
	return PushTail(list, value, DelimiterSlash)
}

func (list Slash) String() string {
	return string(list)
}
