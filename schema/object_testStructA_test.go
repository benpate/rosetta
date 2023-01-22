package schema

type testStructA struct {
	Name     string
	Active   bool
	Latitude float64
	Array    testArrayA
}

func newTestStructA() testStructA {
	return testStructA{
		Name:     "John Connor",
		Active:   true,
		Latitude: 45.123456,
		Array:    newTestArrayA(),
	}
}

func testStructA_Schema() Element {
	return Object{
		Properties: map[string]Element{
			"name":     String{},
			"active":   Boolean{},
			"latitude": Number{BitSize: 64},
			"array":    testArrayA_Schema(),
		},
	}
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (t testStructA) GetBool(path string) (bool, bool) {
	if path == "active" {
		return t.Active, true
	}
	return false, false
}

func (t testStructA) GetFloat(path string) (float64, bool) {
	if path == "latitude" {
		return t.Latitude, true
	}
	return 0, false
}

func (t *testStructA) GetObject(path string) (any, bool) {
	if path == "array" {
		return &t.Array, true
	}
	return nil, false
}

func (t testStructA) GetString(path string) (string, bool) {
	if path == "name" {
		return t.Name, true
	}
	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (t *testStructA) SetBool(path string, value bool) bool {
	if path == "active" {
		t.Active = value
		return true
	}
	return false
}

func (t *testStructA) SetFloat(path string, value float64) bool {
	if path == "latitude" {
		t.Latitude = value
		return true
	}
	return false
}

func (t *testStructA) SetString(path string, value string) bool {
	if path == "name" {
		t.Name = value
		return true
	}
	return false
}
