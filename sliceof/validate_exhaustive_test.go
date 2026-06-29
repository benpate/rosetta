package sliceof_test

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestSliceof_SchemaValidate_AllScalarTypes runs every concrete scalar sliceof type through
// schema.Validate as an Array property. This is the exhaustive guard the array-interface fix calls
// for: if a new sliceof type is added (or an existing one loses GetIndex/SetIndex), validation here
// fails with "Value must implement ArrayGetterSetter interface" instead of failing silently in a
// downstream consumer like Emissary.
func TestSliceof_SchemaValidate_AllScalarTypes(t *testing.T) {

	t.Run("String", func(t *testing.T) {
		value := sliceof.String{"alpha", "bravo"}
		s := schema.New(schema.Array{Items: schema.String{}})
		_, err := s.Validate(&value)
		require.Nil(t, err)
	})

	t.Run("Int", func(t *testing.T) {
		value := sliceof.Int{1, 2, 3}
		s := schema.New(schema.Array{Items: schema.Integer{}})
		_, err := s.Validate(&value)
		require.Nil(t, err)
	})

	t.Run("Float", func(t *testing.T) {
		value := sliceof.Float{1.5, 2.5}
		s := schema.New(schema.Array{Items: schema.Number{}})
		_, err := s.Validate(&value)
		require.Nil(t, err)
	})

	t.Run("Any", func(t *testing.T) {
		value := sliceof.Any{"alpha", "bravo"}
		s := schema.New(schema.Array{Items: schema.String{}})
		_, err := s.Validate(&value)
		require.Nil(t, err)
	})
}

// objectItem models an array element that exposes its scalar fields through VALUE receivers (the
// counterpart to the pointer-receiver case in the schema package's own tests). Both shapes must
// validate, so this confirms the array-element fix did not regress the value-receiver path.
type objectItem struct {
	Name string
}

func objectItem_Schema() schema.Element {
	return schema.Object{Properties: schema.ElementMap{"name": schema.String{MaxLength: 16}}}
}

func (item objectItem) GetStringOK(name string) (string, bool) {
	if name == "name" {
		return item.Name, true
	}
	return "", false
}

func (item *objectItem) SetString(name string, value string) bool {
	if name == "name" {
		item.Name = value
		return true
	}
	return false
}

// TestSliceof_Object_ValueReceiverItems validates a sliceof.Object whose element type uses value
// receivers for its getters. This must succeed via the StringGetter path.
func TestSliceof_Object_ValueReceiverItems(t *testing.T) {

	value := sliceof.Object[objectItem]{{Name: "Alice"}, {Name: "Bob"}}
	s := schema.New(schema.Array{Items: objectItem_Schema()})

	changed, err := s.Validate(&value)
	require.Nil(t, err)
	require.False(t, changed)
}

// TestSliceof_NestedArray validates an array whose items are themselves arrays
// (sliceof.Object[sliceof.String]). This exercises the Array branch of validate_Array's
// addressable-item handling, not just the Object branch.
func TestSliceof_NestedArray(t *testing.T) {

	value := sliceof.Object[sliceof.String]{
		{"alpha", "bravo"},
		{"charlie"},
	}
	s := schema.New(schema.Array{Items: schema.Array{Items: schema.String{}}})

	changed, err := s.Validate(&value)
	require.Nil(t, err)
	require.False(t, changed)
}

// TestSliceof_NestedArray_InnerRewrite confirms a rewrite inside a nested inner array is detected.
// The inner String has MaxLength:3, so "toolong" must be flagged as needing a rewrite (which the
// top-level Validate reports as not-valid-as-given).
func TestSliceof_NestedArray_InnerRewrite(t *testing.T) {

	value := sliceof.Object[sliceof.String]{
		{"toolong"},
	}
	s := schema.New(schema.Array{Items: schema.Array{Items: schema.String{MaxLength: 3}}})

	_, err := s.Validate(&value)
	require.NotNil(t, err)
}
