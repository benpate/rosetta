package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// declaredMap is a map that explicitly declares itself via MapTyper. Absent keys named by the
// schema must be treated as legitimately-empty optional properties.
type declaredMap map[string]string

func (m declaredMap) GetStringOK(key string) (string, bool) {
	value, ok := m[key]
	return value, ok
}

func (m declaredMap) SetString(key string, value string) bool {
	m[key] = value
	return true
}

func (m declaredMap) IsMap() bool { return true }

// TestValidateObject_DeclaredMap_SkipsAbsentKeys confirms that a MapTyper whose IsMap() is true
// validates cleanly even when the schema names keys the map does not contain.
func TestValidateObject_DeclaredMap_SkipsAbsentKeys(t *testing.T) {

	s := New(Object{Properties: ElementMap{
		"present": String{},
		"absent":  String{},
	}})

	value := declaredMap{"present": "here"}

	_, err := s.Validate(value)
	require.Nil(t, err)
}

// brokenStruct is a struct-shaped object that forgot to wire up a PointerGetter (and is NOT a
// MapTyper). Validation must FAIL loudly rather than silently skip the unreadable property.
type brokenStruct struct {
	Name string
}

// (Intentionally no GetStringOK / GetPointer / IsMap.)

// TestValidateObject_StructMissingAccessor_Errors is the footgun guard: an object that is neither a
// declared map nor a working getter must produce an error, not pass validation by accident.
func TestValidateObject_StructMissingAccessor_Errors(t *testing.T) {

	s := New(Object{Properties: ElementMap{
		"name": String{},
	}})

	_, err := s.Validate(brokenStruct{Name: "Sarah Connor"})
	require.NotNil(t, err)
}
