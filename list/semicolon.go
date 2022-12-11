package list

import "strings"

const DelimiterSemicolon = ';'

type Semicolon string

func BySemicolon(value ...string) List {
	return Semicolon(strings.Join(value, string(DelimiterSemicolon)))
}

func (list Semicolon) IsEmpty() bool {
	return IsEmpty(list)
}

func (list Semicolon) IsEmptyTail() bool {
	return IsEmptyTail(list, DelimiterSemicolon)
}

func (list Semicolon) Head() string {
	return Head(list, DelimiterSemicolon)
}

func (list Semicolon) Tail() List {
	return Tail(list, DelimiterSemicolon)
}

func (list Semicolon) Last() string {
	return Last(list, DelimiterSemicolon)
}

func (list Semicolon) RemoveLast() List {
	return RemoveLast(list, DelimiterSemicolon)
}

func (list Semicolon) Split() (string, List) {
	return Split(list, DelimiterSemicolon)
}

func (list Semicolon) SplitTail() (List, string) {
	return SplitTail(list, DelimiterSemicolon)
}

func (list Semicolon) At(index int) string {
	return At(list, DelimiterSemicolon, index)
}

func (list Semicolon) PushHead(value string) List {
	return PushHead(list, value, DelimiterSemicolon)
}

func (list Semicolon) PushTail(value string) List {
	return PushTail(list, value, DelimiterSemicolon)
}

func (list Semicolon) String() string {
	return string(list)
}
