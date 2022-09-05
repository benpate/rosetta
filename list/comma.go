package list

import "strings"

const DelimiterComma = ','

type Comma []byte

func ByComma(value ...string) List {
	return List(Comma(strings.Join(value, string(DelimiterComma))))
}

func (list Comma) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Comma) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterComma)
}

func (list Comma) Head() string {
	return string(Head(list, DelimiterComma))
}

func (list Comma) Tail() List {
	return Comma(Tail(list, DelimiterComma))
}

func (list Comma) Last() string {
	return string(Last(list, DelimiterComma))
}

func (list Comma) RemoveLast() List {
	return Comma(RemoveLast(list, DelimiterComma))
}

func (list Comma) Split() (string, List) {
	head, tail := Split(list, DelimiterComma)
	return string(head), Comma(tail)
}

func (list Comma) SplitTail() (List, string) {
	head, tail := SplitTail(list, DelimiterComma)
	return Comma(head), string(tail)
}

func (list Comma) At(index int) string {
	return string(At(list, DelimiterComma, index))
}

func (list Comma) PushHead(value string) List {
	return Comma(PushHead(list, []byte(value), DelimiterComma))
}

func (list Comma) PushTail(value string) List {
	return Comma(PushTail(list, []byte(value), DelimiterComma))
}

func (list Comma) String() string {
	return string(list)
}
