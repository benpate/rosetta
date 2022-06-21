package list

const DelimiterSlash = '/'

type Slash []byte

func (list Slash) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Slash) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSlash)
}

func (list Slash) Head() string {
	return string(Head(list, DelimiterSlash))
}

func (list Slash) Tail() Slash {
	return Slash(Tail(list, DelimiterSlash))
}

func (list Slash) Last() string {
	return string(Last(list, DelimiterSlash))
}

func (list Slash) RemoveLast() Slash {
	return Slash(RemoveLast(list, DelimiterSlash))
}

func (list Slash) Split() (string, Slash) {
	head, tail := Split(list, DelimiterSlash)
	return string(head), Slash(tail)
}

func (list Slash) SplitTail() (Slash, string) {
	head, tail := SplitTail(list, DelimiterSlash)
	return Slash(head), string(tail)
}

func (list Slash) At(index int) string {
	return string(At(list, DelimiterSlash, index))
}

func (list Slash) PushHead(value string) Slash {
	return Slash(PushHead(list, []byte(value), DelimiterSlash))
}

func (list Slash) PushTail(value string) Slash {
	return Slash(PushTail(list, []byte(value), DelimiterSlash))
}

func (list Slash) String() string {
	return string(list)
}
