package list

func Second(l List) string {
	return l.Tail().Head()
}

func Third(l List) string {
	return l.Tail().Tail().Head()
}

func First2(l List) (string, string) {

	head, tail := l.Split()
	head2 := tail.Head()

	return head, head2
}

func First3(l List) (string, string, string) {

	head, tail := l.Split()
	head2, tail := tail.Split()
	head3 := tail.Head()

	return head, head2, head3
}

func First4(l List) (string, string, string, string) {

	head, tail := l.Split()
	head2, tail := tail.Split()
	head3, tail := tail.Split()
	head4 := tail.Head()

	return head, head2, head3, head4
}

func Last2(l List) (string, string) {

	tail, last := l.SplitTail()
	last2 := tail.Last()

	return last, last2
}

func Last3(l List) (string, string, string) {

	tail, last := l.SplitTail()
	tail, last2 := tail.SplitTail()
	last3 := tail.Last()

	return last, last2, last3
}

func Last4(l List) (string, string, string, string) {

	tail, last := l.SplitTail()
	tail, last2 := tail.SplitTail()
	tail, last3 := tail.SplitTail()
	last4 := tail.Last()

	return last, last2, last3, last4
}
