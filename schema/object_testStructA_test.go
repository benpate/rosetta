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

func (t testStructA) GetBoolOK(path string) (bool, bool) {
	if path == "active" {
		return t.Active, true
	}
	return false, false
}

func (t testStructA) GetFloatOK(path string) (float64, bool) {
	if path == "latitude" {
		return t.Latitude, true
	}
	return 0, false
}

func (t *testStructA) GetObjectOK(path string) (any, bool) {
	if path == "array" {
		return &t.Array, true
	}
	return nil, false
}

func (t testStructA) GetStringOK(path string) (string, bool) {
	if path == "name" {
		return t.Name, true
	}
	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (t *testStructA) SetBoolOK(path string, value bool) bool {
	if path == "active" {
		t.Active = value
		return true
	}
	return false
}

func (t *testStructA) SetFloatOK(path string, value float64) bool {
	if path == "latitude" {
		t.Latitude = value
		return true
	}
	return false
}

func (t *testStructA) SetStringOK(path string, value string) bool {
	if path == "name" {
		t.Name = value
		return true
	}
	return false
}
