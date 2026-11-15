package deep

import "errors"

// These types have the same layout as the MongoDB driver's primitive.A, M, DateTime, and
// ObjectID. Rosetta does not import the driver, so Copy and Remove reach them by reflection.
type (
	testArray    []any
	testDocument map[string]any
	testDateTime int64
	testObjectID [12]byte
)

// These types have the same layout as the driver's Decimal128, Timestamp, Regex, and Null:
// structs that hold no map, slice, or pointer
type (
	testDecimal   struct{ h, l uint64 }
	testTimestamp struct{ T, I uint32 }
	testRegex     struct{ Pattern, Options string }
	testNull      struct{}
)

// testStruct is a type that Copy and Remove refuse, because it holds a map
type testStruct struct {
	Secret string
	Hidden map[string]any
}

// testCopier is a struct that copies itself, and counts how often it was asked
type testCopier struct {
	Name   string
	Copies *int
}

// Copy returns this value, and records the call
func (value testCopier) Copy(int) (any, error) {
	*value.Copies++
	return value, nil
}

// testWrongCopier is a Copier that returns a different type than its own
type testWrongCopier struct{}

// Copy returns a string, which no container of testWrongCopier can hold
func (testWrongCopier) Copy(int) (any, error) {
	return "not a testWrongCopier", nil
}

// testFailingCopier is a Copier and a Remover that always fails
type testFailingCopier struct{}

// Copy always fails
func (testFailingCopier) Copy(int) (any, error) {
	return "partial", errors.New("copy failed")
}

// Remove always fails
func (testFailingCopier) Remove(int, Matcher) error {
	return errors.New("remove failed")
}

// testRemover is a struct that records which of its names the Matcher selected
type testRemover struct {
	Names   []string
	Removed *[]string
}

// Remove records the names that match
func (value testRemover) Remove(_ int, match Matcher) error {

	for _, name := range value.Names {
		if match(name) {
			*value.Removed = append(*value.Removed, name)
		}
	}

	return nil
}

// testBox is a pointer type that copies and searches the value it holds, one level deeper
type testBox struct {
	Child any
	Depth int // The depth that Copy or Remove was last called with
}

// Copy returns a new box that holds a deep copy of the child
func (box *testBox) Copy(depth int) (any, error) {

	box.Depth = depth

	child, err := CopyWithDepth(box.Child, depth+1)

	if err != nil {
		return nil, err
	}

	return &testBox{Child: child}, nil
}

// Remove deletes the named properties from the child
func (box *testBox) Remove(depth int, match Matcher) error {
	box.Depth = depth
	return RemoveWithDepth(box.Child, depth+1, match)
}
