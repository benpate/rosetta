package schema

type testArrayB []testStructA

func newTestArrayB() testArrayB {
	return testArrayB{
		newTestStructA(),
		newTestStructA(),
		newTestStructA(),
	}
}

func testArrayB_Schema() Element {
	return Array{
		Items:     testStructA_Schema(),
		MaxLength: 5,
	}
}

func (t testArrayB) Length() int {
	if isNil(t) {
		return 0
	}
	return len(t)
}

func (t testArrayB) GetObject(path string) (any, bool) {
	if index, ok := Index(path, len(t)); ok {
		return t[index], true
	}

	return "", false
}
