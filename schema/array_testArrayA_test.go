package schema

type testArrayA []string

func newTestArrayA() testArrayA {
	return testArrayA{
		"one",
		"two",
		"three",
	}
}

// testArrayA_Schema defines the schema for testArrayA
func testArrayA_Schema() Element {
	return Array{
		Items:     String{},
		MaxLength: 5,
	}
}

// Length returns the length of the array
func (t testArrayA) Length() int {
	if isNil(t) {
		return 0
	}
	return len(t)
}

func (t testArrayA) GetIndex(index int) (any, bool) {
	if index >= 0 && index < len(t) {
		return t[index], true
	}

	return nil, false
}

func (t *testArrayA) SetIndex(index int, value any) bool {
	if index >= 0 && index < len(*t) {
		if typedValue, ok := value.(string); ok {
			(*t)[index] = typedValue
			return true
		}
	}

	return false
}

// GetPointer implements the Getter interface
func (t testArrayA) GetStringOK(path string) (string, bool) {
	if index, ok := Index(path, len(t)); ok {
		return t[index], true
	}

	return "", false
}

// SetString implements the Setter interface
func (t *testArrayA) SetString(path string, value string) bool {

	if index, ok := Index(path); ok {

		for index >= len(*t) {
			*t = append(*t, "")
		}

		(*t)[index] = value
		return true
	}

	return false
}

// Remove implements the Remover interface
func (t *testArrayA) Remove(path string) bool {

	if index, ok := Index(path); ok {

		if index < len(*t) {
			*t = append((*t)[:index], (*t)[index+1:]...)
			return true
		}
	}

	return false
}
