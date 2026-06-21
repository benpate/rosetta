package list

// Second returns the second item in the list.
func Second(l List) string {
	return l.Tail().Head()
}

// Third returns the third item in the list.
func Third(l List) string {
	return l.Tail().Tail().Head()
}

// First2 returns the first two items in the list.
func First2(l List) (string, string) {

	head, tail := l.Split()
	head2 := tail.Head()

	return head, head2
}

// First3 returns the first three items in the list.
func First3(l List) (string, string, string) {

	head, tail := l.Split()
	head2, tail := tail.Split()
	head3 := tail.Head()

	return head, head2, head3
}

// First4 returns the first four items in the list.
func First4(l List) (string, string, string, string) {

	head, tail := l.Split()
	head2, tail := tail.Split()
	head3, tail := tail.Split()
	head4 := tail.Head()

	return head, head2, head3, head4
}

// Last2 returns the final two items in the list (last, then second-to-last).
func Last2(l List) (string, string) {

	tail, last := l.SplitTail()
	last2 := tail.Last()

	return last, last2
}

// Last3 returns the final three items in the list (last first).
func Last3(l List) (string, string, string) {

	tail, last := l.SplitTail()
	tail, last2 := tail.SplitTail()
	last3 := tail.Last()

	return last, last2, last3
}

// Last4 returns the final four items in the list (last first).
func Last4(l List) (string, string, string, string) {

	tail, last := l.SplitTail()
	tail, last2 := tail.SplitTail()
	tail, last3 := tail.SplitTail()
	last4 := tail.Last()

	return last, last2, last3, last4
}
