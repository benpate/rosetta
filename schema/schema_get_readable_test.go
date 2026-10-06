package schema

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestScalarValue covers each kind of value a list's GetPointer can return
func TestScalarValue(t *testing.T) {

	text, number := "a", 2
	child := mapof.Any{"k": "v"}
	var nilText *string

	require.Equal(t, "a", scalarValue(&text))
	require.Equal(t, 2, scalarValue(&number))
	require.Same(t, &child, scalarValue(&child), "a pointer to a map is kept")
	require.Nil(t, scalarValue(nilText))
	require.Equal(t, "a", scalarValue("a"))
	require.Nil(t, scalarValue(nil))
}

// TestHasEmptySegment covers every place an empty segment can sit, and the empty path itself
func TestHasEmptySegment(t *testing.T) {

	for _, path := range []string{"name.", ".name", "a..b", "."} {
		require.True(t, hasEmptySegment(path), path)
	}

	for _, path := range []string{"", "name", "a.b", "a b", "0"} {
		require.False(t, hasEmptySegment(path), path)
	}
}

// TestIsPastEnd covers indexes inside, at, and past the end of a list, and keys that are not
// indexes
func TestIsPastEnd(t *testing.T) {

	require.False(t, isPastEnd("0", 1))
	require.False(t, isPastEnd("+0", 1), "strconv.Atoi accepts a sign, as sliceof does")
	require.True(t, isPastEnd("1", 1))
	require.True(t, isPastEnd("+1", 1))
	require.True(t, isPastEnd("next", 1))
	require.True(t, isPastEnd("99999999999999999999999", 1), "too large for an int")
	require.False(t, isPastEnd("-1", 1), "sliceof refuses a negative index without growing")
	require.False(t, isPastEnd("last", 1))
	require.False(t, isPastEnd("", 1))
}

// TestZeroValueAt covers the zero value of every element type, and a path the schema lacks
func TestZeroValueAt(t *testing.T) {

	element := Object{Properties: ElementMap{
		"text":  String{Default: "ignored"},
		"flag":  Boolean{},
		"count": Integer{},
		"big":   Integer{BitSize: 64},
		"ratio": Number{},
		"any":   Any{},
		"list":  Array{Items: String{}},
		"child": Object{Properties: ElementMap{"name": String{}}},
	}}

	cases := map[string]any{
		"text":       "",
		"flag":       false,
		"count":      0,
		"big":        int64(0),
		"ratio":      float64(0),
		"any":        nil,
		"list":       nil,
		"child":      nil,
		"child.name": "",
	}

	for path, expected := range cases {
		value, err := zeroValueAt(element, path)
		require.NoError(t, err, path)
		require.Equal(t, expected, value, path)
	}

	_, err := zeroValueAt(element, "missing")
	require.Error(t, err)
}

// TestIsDigits covers digits, signs, and the empty string
func TestIsDigits(t *testing.T) {
	require.True(t, isDigits("0"))
	require.True(t, isDigits("007"))
	require.False(t, isDigits(""))
	require.False(t, isDigits("+1"))
	require.False(t, isDigits("1a"))
}

// TestPointerToCopy covers each shape of value a read can start from
func TestPointerToCopy(t *testing.T) {

	list := sliceof.Any{"a"}
	require.Equal(t, &list, pointerToCopy(list), "a value list is read through a pointer to a copy")
	plain := []any{"a"}
	require.Equal(t, &sliceof.Any{"a"}, pointerToCopy(plain), "a plain slice reads as sliceof.Any")
	require.Equal(t, &mapof.Any{"k": "v"}, pointerToCopy(map[string]any{"k": "v"}))
	require.Same(t, &list, pointerToCopy(&list))
	require.Equal(t, 7, pointerToCopy(7))
	held := followupStruct{}
	require.Equal(t, held, pointerToCopy(held), "a struct held by value is kept")
	require.Nil(t, pointerToCopy(nil))
}

// followupStruct is a struct whose getter has a pointer receiver
type followupStruct struct {
	Name string
}

// GetPointer implements the PointerGetter interface
func (value *followupStruct) GetPointer(name string) (any, bool) {
	if name == "name" {
		return &value.Name, true
	}
	return nil, false
}

// TestFollowupGuard_StructHeldByValue pins that Set still refuses a whole struct held by value
// whose schema it must validate, as Emissary's newStartupStream relies on for model.Content
func TestFollowupGuard_StructHeldByValue(t *testing.T) {

	element := Object{Properties: ElementMap{
		"child": Object{Properties: ElementMap{"name": String{Required: true}}},
	}}

	object := mapof.Any{}
	require.Error(t, New(element).Set(&object, "child", followupStruct{Name: "v"}))

	_, err := New(element).Get(&mapof.Any{"child": followupStruct{Name: "v"}}, "child.name")
	require.Error(t, err, "a struct held by value has no PointerGetter")
}
