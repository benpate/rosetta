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

// testStruct is a type that Copy and Remove refuse
type testStruct struct{ Secret string }

// testCopier is a struct that copies itself, and counts how often it was asked
type testCopier struct {
	Name   string
	Copies *int
}

// Copy returns this value, and records the call
func (value testCopier) Copy() (any, error) {
	*value.Copies++
	return value, nil
}

// testWrongCopier is a Copier that returns a different type than its own
type testWrongCopier struct{}

// Copy returns a string, which no container of testWrongCopier can hold
func (testWrongCopier) Copy() (any, error) {
	return "not a testWrongCopier", nil
}

// testFailingCopier is a Copier and a Remover that always fails
type testFailingCopier struct{}

// Copy always fails
func (testFailingCopier) Copy() (any, error) {
	return "partial", errors.New("copy failed")
}

// Remove always fails
func (testFailingCopier) Remove(...string) error {
	return errors.New("remove failed")
}

// testRemover is a struct that records the names it was asked to remove
type testRemover struct{ Removed *[]string }

// Remove records the names
func (value testRemover) Remove(propertyNames ...string) error {
	*value.Removed = append(*value.Removed, propertyNames...)
	return nil
}
