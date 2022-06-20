package convert

type testStringer []string

func (t testStringer) String() string {
	return t[0]
}

func getTestStringer() testStringer {
	return testStringer{"hello", "there", "general", "grievous"}
}
