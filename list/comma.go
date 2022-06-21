package list

const DelimiterComma = ','

type Comma []byte

func (list Comma) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Comma) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterComma)
}

func (list Comma) Head() string {
	return string(Head(list, DelimiterComma))
}

func (list Comma) Tail() Comma {
	return Comma(Tail(list, DelimiterComma))
}

func (list Comma) Last() string {
	return string(Last(list, DelimiterComma))
}

func (list Comma) RemoveLast() Comma {
	return Comma(RemoveLast(list, DelimiterComma))
}

func (list Comma) Split() (string, Comma) {
	head, tail := Split(list, DelimiterComma)
	return string(head), Comma(tail)
}

func (list Comma) SplitTail() (Comma, string) {
	head, tail := SplitTail(list, DelimiterComma)
	return Comma(head), string(tail)
}

func (list Comma) At(index int) string {
	return string(At(list, DelimiterComma, index))
}

func (list Comma) PushHead(value string) Comma {
	return Comma(PushHead(list, []byte(value), DelimiterComma))
}

func (list Comma) PushTail(value string) Comma {
	return Comma(PushTail(list, []byte(value), DelimiterComma))
}

func (list Comma) String() string {
	return string(list)
}
