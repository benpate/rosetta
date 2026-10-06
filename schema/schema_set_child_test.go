package schema

import (
	"reflect"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// rejectingMap is a KeySetter that refuses every value that is not a string
type rejectingMap map[string]any

// SetKey implements the KeySetter interface
func (m *rejectingMap) SetKey(key string, value any) error {
	if _, ok := value.(string); !ok {
		const location = "schema.rejectingMap.SetKey"
		return derp.Internal(location, "Only strings are stored", typeName(value))
	}
	(*m)[key] = value
	return nil
}

// GetPointer implements the PointerGetter interface
func (m *rejectingMap) GetPointer(key string) (any, bool) {
	value, ok := (*m)[key]
	return value, ok
}

// anySetterMap is a map that declares itself through MapTyper and stores keys through SetAny,
// with no KeySetter
type anySetterMap map[string]any

// IsMap implements the MapTyper interface
func (anySetterMap) IsMap() bool {
	return true
}

// SetAny implements the AnySetter interface
func (m *anySetterMap) SetAny(key string, value any) bool {
	(*m)[key] = value
	return true
}

// emptyList is a list whose GetPointer always reports a missing item, and whose SetAny stores
// nothing when refuse is set
type emptyList struct {
	refuse bool
}

// Length implements the ArrayGetter interface
func (*emptyList) Length() int {
	return 0
}

// GetIndex implements the ArrayGetter interface
func (*emptyList) GetIndex(int) (any, bool) {
	return nil, false
}

// GetPointer implements the PointerGetter interface
func (*emptyList) GetPointer(string) (any, bool) {
	return nil, true
}

// SetAny implements the AnySetter interface
func (list *emptyList) SetAny(string, any) bool {
	return !list.refuse
}

// pointerOnlyList is a list whose GetPointer always reports a missing item, with no SetAny
type pointerOnlyList struct{}

// Length implements the ArrayGetter interface
func (*pointerOnlyList) Length() int {
	return 0
}

// GetIndex implements the ArrayGetter interface
func (*pointerOnlyList) GetIndex(int) (any, bool) {
	return nil, false
}

// GetPointer implements the PointerGetter interface
func (*pointerOnlyList) GetPointer(string) (any, bool) {
	return nil, true
}

// TestSetMapChild_SetKeyErrors requires that a KeySetter's refusal is reported, for a single key
// and for a child stored back after a nested write
func TestSetMapChild_SetKeyErrors(t *testing.T) {

	element := Object{Wildcard: Any{}}

	object := rejectingMap{}
	require.Error(t, SetProperty(element, &object, "count", 6))
	err := SetProperty(element, &object, "child.name", "v")
	require.Error(t, err, "a new mapof.Any child is refused")
	require.Equal(t, rejectingMap{}, object)

	require.NoError(t, SetProperty(element, &object, "name", "v"))
	require.Equal(t, rejectingMap{"name": "v"}, object)
}

// TestSetTypedMapKey covers each setter a map without KeySetter can store a key through
func TestSetTypedMapKey(t *testing.T) {

	withSetAny := anySetterMap{}
	require.True(t, setTypedMapKey(&withSetAny, "count", 6))
	require.Equal(t, anySetterMap{"count": 6}, withSetAny)

	withSetString := mapof.String{}
	require.True(t, setTypedMapKey(&withSetString, "name", "v"))
	require.False(t, setTypedMapKey(&withSetString, "count", 6), "a non-string cannot be stored")
	require.Equal(t, mapof.String{"name": "v"}, withSetString)

	require.False(t, setTypedMapKey(&struct{}{}, "name", "v"), "a value that is not a map")
}

// TestSetOtherChild_NoInterfaces requires an error for a target with no setter interfaces,
// beneath an Object and beneath an Any
func TestSetOtherChild_NoInterfaces(t *testing.T) {

	target := struct{ Name string }{}

	require.Error(t, SetProperty(Object{Wildcard: Any{}}, &target, "name", "v"))
	require.Error(t, SetProperty(Any{}, &target, "name.first", "v"))
	require.Empty(t, target.Name)
}

// TestSetListItem_WithoutSetAny pins a whole item written into a list with no SetAny, which
// falls back to the item's pointer
func TestSetListItem_WithoutSetAny(t *testing.T) {

	list := sliceof.Object[mapof.Any]{{"a": 1}}
	element := Array{Items: Object{Wildcard: Any{}}}

	require.False(t, setListItem(&list, "0", mapof.Any{"b": 2}))
	require.False(t, setListItem(&struct{}{}, "0", mapof.Any{"b": 2}), "a value that is not a list")
	require.NoError(t, SetProperty(element, &list, "0", mapof.Any{"b": 2}))
	require.Equal(t, sliceof.Object[mapof.Any]{{"b": 2}}, list)
}

// TestSetPointerChild_ItemOfWrongKind requires an error when a list item cannot hold the rest
// of the path
func TestSetPointerChild_ItemOfWrongKind(t *testing.T) {

	list := sliceof.Any{"scalar"}
	element := Array{Items: Object{Properties: ElementMap{"name": String{}}}}

	require.Error(t, SetProperty(element, &list, "0.name", "v"))
	require.Equal(t, sliceof.Any{"scalar"}, list)
}

// TestCreateListItem_Errors covers each way a missing list item can fail to be created
func TestCreateListItem_Errors(t *testing.T) {

	element := Array{Items: Object{Properties: ElementMap{"name": String{}}}}

	err := SetProperty(element, &pointerOnlyList{}, "0.name", "v")
	require.Error(t, err, "the list has no SetAny")

	err = SetProperty(element, &emptyList{refuse: true}, "0.name", "v")
	require.Error(t, err, "the list refuses the item")

	require.NoError(t, SetProperty(element, &emptyList{}, "0.name", "v"))

	list := sliceof.Any{}
	err = SetProperty(element, &list, "0.nope", "v")
	require.Error(t, err, "the item's schema has no such property")
	require.Equal(t, sliceof.Any{nil}, list, "GetPointer grew the list before the item failed")
}

// TestStoreListItem requires that a list item that is itself a list keeps its new length,
// although the parent's GetPointer handed back a copy
func TestStoreListItem(t *testing.T) {

	list := sliceof.Any{sliceof.Any{}}
	require.NoError(t, SetProperty(Any{}, &list, "0.0", "v"))
	require.Equal(t, sliceof.Any{sliceof.Any{"v"}}, list)
}

// TestCopyChild_CannotWrite requires an error for a child that has no setter interfaces and no
// rosetta type with its layout
func TestCopyChild_CannotWrite(t *testing.T) {

	_, err := copyChild([]string{"a"})
	require.Error(t, err)

	object := mapof.Any{"child": []string{"a"}}
	require.Error(t, SetProperty(Object{Wildcard: Any{}}, &object, "child.0", "v"))
	require.Equal(t, mapof.Any{"child": []string{"a"}}, object)
}

// TestNewTypedChild covers a typed parent whose children are pointers, and one whose children
// are single values
func TestNewTypedChild(t *testing.T) {

	child := Object{Properties: ElementMap{"name": String{}}}
	element := Object{Properties: ElementMap{"child": child}}

	pointers := mapof.Object[*mapof.Any]{}
	require.NoError(t, SetProperty(element, &pointers, "child.name", "v"))
	require.Equal(t, mapof.Any{"name": "v"}, *pointers["child"])

	strings := mapof.Object[string]{}
	err := SetProperty(element, &strings, "child.name", "v")
	require.Error(t, err)
	require.Equal(t, "Expected a map or a list", derp.RootMessage(err), "a missing string is \"\"")
}

// TestChildKind_String names every kind
func TestChildKind_String(t *testing.T) {
	require.Equal(t, "any value", kindAny.String())
	require.Equal(t, "a map", kindMap.String())
	require.Equal(t, "a list", kindList.String())
	require.Equal(t, "a single value", kindScalar.String())
}

// TestIsSlicePointer covers every shape of value
func TestIsSlicePointer(t *testing.T) {

	var nilList *sliceof.Any

	require.True(t, isSlicePointer(&sliceof.Any{}))
	require.False(t, isSlicePointer(sliceof.Any{}))
	require.False(t, isSlicePointer(nilList))
	require.False(t, isSlicePointer(&mapof.Any{}))
	require.False(t, isSlicePointer(nil))
}

// TestKindOf covers every kind a type can have
func TestKindOf(t *testing.T) {
	require.Equal(t, kindMap, kindOf(reflect.TypeFor[mapof.Any]()))
	require.Equal(t, kindMap, kindOf(reflect.TypeFor[*struct{}]()))
	require.Equal(t, kindList, kindOf(reflect.TypeFor[[2]int]()))
	require.Equal(t, kindScalar, kindOf(reflect.TypeFor[string]()))
}
