package list

import "strings"

const DelimiterComma = ','

type Comma string

func ByComma(value ...string) List {
	return Comma(strings.Join(value, string(DelimiterComma)))
}

func (list Comma) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Comma) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterComma)
}

func (list Comma) Head() string {
	return Head(list, DelimiterComma)
}

func (list Comma) Tail() List {
	return Tail(list, DelimiterComma)
}

func (list Comma) Last() string {
	return Last(list, DelimiterComma)
}

func (list Comma) RemoveLast() List {
	return RemoveLast(list, DelimiterComma)
}

func (list Comma) Split() (string, List) {
	return Split(list, DelimiterComma)
}

func (list Comma) SplitTail() (List, string) {
	return SplitTail(list, DelimiterComma)
}

func (list Comma) At(index int) string {
	return At(list, DelimiterComma, index)
}

func (list Comma) PushHead(value string) List {
	return PushHead(list, value, DelimiterComma)
}

func (list Comma) PushTail(value string) List {
	return PushTail(list, value, DelimiterComma)
}

func (list Comma) String() string {
	return string(list)
}
