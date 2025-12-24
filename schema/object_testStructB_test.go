package schema

type testStructB struct {
	Name      string
	Age       int
	Published int64
	Array     testArrayB
}

func newTestStructB() testStructB {
	return testStructB{
		Name:      "John Connor",
		Age:       42,
		Published: 1234567890,
		Array:     newTestArrayB(),
	}
}

func testStructB_Schema() Element {
	return Object{
		Properties: map[string]Element{
			"name":      String{},
			"age":       Integer{BitSize: 32},
			"published": Integer{BitSize: 64},
			"array":     testArrayB_Schema(),
		},
	}
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (t testStructB) GetIntOK(path string) (int, bool) {
	if path == "age" {
		return t.Age, true
	}
	return 0, false
}

func (t testStructB) GetInt64OK(path string) (int64, bool) {
	if path == "published" {
		return t.Published, true
	}
	return 0, false
}

func (t *testStructB) GetPointer(path string) (any, bool) {
	if path == "array" {
		return &t.Array, true
	}
	return nil, false
}

func (t testStructB) GetStringOK(path string) (string, bool) {
	if path == "name" {
		return t.Name, true
	}
	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (t *testStructB) SetInt(path string, value int) bool {
	if path == "age" {
		t.Age = value
		return true
	}
	return false
}

func (t *testStructB) SetInt64(path string, value int64) bool {
	if path == "published" {
		t.Published = value
		return true
	}
	return false
}

func (t *testStructB) SetString(path string, value string) bool {
	if path == "name" {
		t.Name = value
		return true
	}
	return false
}
