package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// nestedElements are the three element kinds whose nested paths schema walks into a map,
// each with the path that reaches a String leaf named "name" beneath "child"
var nestedElements = map[string]struct {
	Element schema.Element
	Path    string
}{
	"any":    {Element: schema.Any{}, Path: "child.name"},
	"object": {Element: schema.Object{Properties: schema.ElementMap{"name": schema.String{}}}, Path: "child.name"},
	"array":  {Element: schema.Array{Items: schema.Object{Properties: schema.ElementMap{"name": schema.String{}}}}, Path: "child.0.name"},
}

// nestedSchema returns a schema whose "child" property is the named element
func nestedSchema(kind string) schema.Schema {
	return schema.New(schema.Object{Properties: schema.ElementMap{"child": nestedElements[kind].Element}})
}

// TestSet_Nested_MapofAny pins what Schema.Set does to each kind of existing child when it
// writes a nested path into a mapof.Any
func TestSet_Nested_MapofAny(t *testing.T) {

	// BUG-234: lists stay lists, other map types keep their keys, and a wrong kind is an error

	// set writes "v" at the element's path into a mapof.Any that holds the given child
	set := func(t *testing.T, kind string, child any) (mapof.Any, error) {
		t.Helper()
		object := mapof.Any{"sibling": "s"}
		if child != nil {
			object["child"] = child
		}
		err := nestedSchema(kind).Set(&object, nestedElements[kind].Path, "v")
		return object, err
	}

	t.Run("absent child becomes a mapof.Any", func(t *testing.T) {
		for _, kind := range []string{"any", "object"} {
			object, err := set(t, kind, nil)
			require.NoError(t, err, kind)
			require.Equal(t, mapof.Any{"sibling": "s", "child": mapof.Any{"name": "v"}}, object, kind)
		}
	})

	t.Run("absent child under an Array element becomes a list, held by pointer", func(t *testing.T) {
		object, err := set(t, "array", nil)
		require.NoError(t, err)
		require.Equal(t, mapof.Any{"sibling": "s", "child": &sliceof.Any{mapof.Any{"name": "v"}}}, object)
	})

	t.Run("existing mapof.Any child is reused, keeping its other keys", func(t *testing.T) {
		for _, kind := range []string{"any", "object"} {
			object, err := set(t, kind, mapof.Any{"keep": "k"})
			require.NoError(t, err, kind)
			require.Equal(t, mapof.Any{"sibling": "s", "child": mapof.Any{"keep": "k", "name": "v"}}, object, kind)
		}
	})

	t.Run("existing map of another type is written into, keeping its type and other keys", func(t *testing.T) {
		for _, kind := range []string{"any", "object"} {
			object, err := set(t, kind, map[string]any{"keep": "k"})
			require.NoError(t, err, kind)
			require.Equal(t, mapof.Any{"sibling": "s", "child": map[string]any{"keep": "k", "name": "v"}}, object, kind)

			object, err = set(t, kind, mapof.String{"keep": "k"})
			require.NoError(t, err, kind)
			require.Equal(t, mapof.Any{"sibling": "s", "child": mapof.String{"keep": "k", "name": "v"}}, object, kind)

			pointer := &mapof.Any{"keep": "k"}
			object, err = set(t, kind, pointer)
			require.NoError(t, err, kind)
			require.Same(t, pointer, object["child"], kind)
			require.Equal(t, mapof.Any{"keep": "k", "name": "v"}, *pointer, kind)
		}
	})

	t.Run("existing list is written into, keeping its type and other items", func(t *testing.T) {
		// This is Bandwagon's "artists.0.id", which used to turn the list into {"0": {...}}
		pointer := &sliceof.Any{mapof.Any{"type": "Artist"}}
		object, err := set(t, "array", pointer)
		require.NoError(t, err)
		require.Same(t, pointer, object["child"])
		require.Equal(t, sliceof.Any{mapof.Any{"type": "Artist", "name": "v"}}, *pointer)

		for _, child := range []any{
			sliceof.Any{mapof.Any{"type": "Artist"}},
			[]any{mapof.Any{"type": "Artist"}},
			sliceof.Object[mapof.Any]{{"type": "Artist"}},
		} {
			object, err := set(t, "array", child)
			require.NoError(t, err, "%T", child)
			require.IsType(t, child, object["child"])
			require.Equal(t, `[{"name":"v","type":"Artist"}]`, toJSON(t, object["child"]), "%T", child)
		}
	})

	t.Run("a child of the wrong kind is an error, and nothing changes", func(t *testing.T) {
		cases := []struct {
			kind  string
			child any
		}{
			{kind: "any", child: "scalar"},
			{kind: "object", child: "scalar"},
			{kind: "object", child: &sliceof.Any{}},
			{kind: "object", child: []any{}},
			{kind: "array", child: "scalar"},
			{kind: "array", child: mapof.Any{"keep": "k"}},
			{kind: "array", child: map[string]any{"keep": "k"}},
		}

		for _, c := range cases {
			object, err := set(t, c.kind, c.child)
			require.Error(t, err, "%s %T", c.kind, c.child)
			require.Equal(t, mapof.Any{"sibling": "s", "child": c.child}, object, "%s %T", c.kind, c.child)
		}
	})
}

// TestSet_Nested_AnyIndex pins how a missing child beneath an Any element is chosen from the
// next path segment, and how far a list beneath Any may grow in one write
func TestSet_Nested_AnyIndex(t *testing.T) {

	s := nestedSchema("any")

	t.Run("an index creates a list, held by pointer", func(t *testing.T) {
		object := mapof.Any{}
		require.NoError(t, s.Set(&object, "child.0.id", "v"))
		require.Equal(t, mapof.Any{"child": &sliceof.Any{mapof.Any{"id": "v"}}}, object)

		require.NoError(t, s.Set(&object, "child.1.id", "w"), "the next item may be appended")
		require.NoError(t, s.Set(&object, "child.0.name", "x"), "an existing item is written into")
		require.Equal(t, mapof.Any{"child": &sliceof.Any{mapof.Any{"id": "v", "name": "x"}, mapof.Any{"id": "w"}}}, object)
	})

	t.Run("an index past the end of the list is an error", func(t *testing.T) {
		object := mapof.Any{}
		require.Error(t, s.Set(&object, "child.1.id", "v"))
		require.Error(t, s.Set(&object, "child.99999999.id", "v"))
		require.Equal(t, mapof.Any{}, object)

		object = mapof.Any{"child": &sliceof.Any{"a"}}
		require.Error(t, s.Set(&object, "child.2", "v"))
		require.NoError(t, s.Set(&object, "child.1", "b"))
		require.Equal(t, mapof.Any{"child": &sliceof.Any{"a", "b"}}, object)
	})

	t.Run("a segment that is not plain digits creates a map", func(t *testing.T) {
		for _, segment := range []string{"007", "+1", "-1", "1e3", "x"} {
			object := mapof.Any{}
			require.NoError(t, s.Set(&object, "child."+segment+".id", "v"), segment)
			require.Equal(t, mapof.Any{"child": mapof.Any{segment: mapof.Any{"id": "v"}}}, object, segment)
		}
	})

	t.Run("a word is not a list index", func(t *testing.T) {
		object := mapof.Any{"child": &sliceof.Any{}}
		require.Error(t, s.Set(&object, "child.last.id", "v"))
	})
}

// TestSet_Nested_MapofTemplate pins the same matrix for a mapof.Template, whose child maps are
// plain mapof.Any values
func TestSet_Nested_MapofTemplate(t *testing.T) {

	set := func(t *testing.T, kind string, child any) (mapof.Template, error) {
		t.Helper()
		object := mapof.Template{"sibling": "s"}
		if child != nil {
			object["child"] = child
		}
		err := nestedSchema(kind).Set(&object, nestedElements[kind].Path, "v")
		return object, err
	}

	t.Run("absent child becomes a mapof.Any", func(t *testing.T) {
		object, err := set(t, "object", nil)
		require.NoError(t, err)
		require.Equal(t, mapof.Template{"sibling": "s", "child": mapof.Any{"name": "v"}}, object)
	})

	t.Run("absent child under an Array element becomes a list", func(t *testing.T) {
		object, err := set(t, "array", nil)
		require.NoError(t, err)
		require.Equal(t, mapof.Template{"sibling": "s", "child": &sliceof.Any{mapof.Any{"name": "v"}}}, object)
	})

	t.Run("existing mapof.Any child is reused", func(t *testing.T) {
		object, err := set(t, "object", mapof.Any{"keep": "k"})
		require.NoError(t, err)
		require.Equal(t, mapof.Template{"sibling": "s", "child": mapof.Any{"keep": "k", "name": "v"}}, object)
	})

	t.Run("existing map of another type is written into", func(t *testing.T) {
		object, err := set(t, "object", map[string]any{"keep": "k"})
		require.NoError(t, err)
		require.Equal(t, mapof.Template{"sibling": "s", "child": map[string]any{"keep": "k", "name": "v"}}, object)
	})

	t.Run("existing list is written into", func(t *testing.T) {
		object, err := set(t, "array", &sliceof.Any{mapof.Any{"type": "Artist"}})
		require.NoError(t, err)
		require.Equal(t, mapof.Template{"sibling": "s", "child": &sliceof.Any{mapof.Any{"type": "Artist", "name": "v"}}}, object)
	})
}

// TestSet_Nested_MapofObject pins nested writes into a mapof.Object[mapof.Any], whose children
// can only be mapof.Any values
func TestSet_Nested_MapofObject(t *testing.T) {

	set := func(t *testing.T, kind string, object mapof.Object[mapof.Any]) (mapof.Object[mapof.Any], error) {
		t.Helper()
		err := nestedSchema(kind).Set(&object, nestedElements[kind].Path, "v")
		return object, err
	}

	t.Run("absent child is created", func(t *testing.T) {
		object, err := set(t, "object", mapof.Object[mapof.Any]{})
		require.NoError(t, err)
		require.Equal(t, mapof.Object[mapof.Any]{"child": mapof.Any{"name": "v"}}, object)
	})

	t.Run("existing child is reused", func(t *testing.T) {
		object, err := set(t, "object", mapof.Object[mapof.Any]{"child": {"keep": "k"}})
		require.NoError(t, err)
		require.Equal(t, mapof.Object[mapof.Any]{"child": {"keep": "k", "name": "v"}}, object)
	})

	t.Run("an Array element cannot be held by a map of maps", func(t *testing.T) {
		object, err := set(t, "array", mapof.Object[mapof.Any]{})
		require.Error(t, err)
		require.Equal(t, mapof.Object[mapof.Any]{}, object)
	})
}

// nestedRecord is a struct that holds a map, the way Emissary's Stream holds "data"
type nestedRecord struct {
	Data mapof.Any
}

// GetPointer implements the schema PointerGetter interface
func (record *nestedRecord) GetPointer(name string) (any, bool) {
	if name == "data" {
		return &record.Data, true
	}
	return nil, false
}

// TestSet_Nested_StructField pins a nested write that reaches a map through a struct field,
// as Emissary's "data.links.SPOTIFY" form fields do
func TestSet_Nested_StructField(t *testing.T) {

	s := schema.New(schema.Object{Properties: schema.ElementMap{
		"data": schema.Object{Wildcard: schema.Any{}, Properties: schema.ElementMap{
			"links": schema.Object{Properties: schema.ElementMap{"SPOTIFY": schema.String{}, "BANDCAMP": schema.String{}}},
		}},
	}})

	t.Run("absent nested map is created", func(t *testing.T) {
		record := nestedRecord{}
		require.NoError(t, s.Set(&record, "data.links.SPOTIFY", "https://spotify.example/a"))
		require.Equal(t, mapof.Any{"links": mapof.Any{"SPOTIFY": "https://spotify.example/a"}}, record.Data)
	})

	t.Run("existing nested map keeps its other keys", func(t *testing.T) {
		record := nestedRecord{Data: mapof.Any{"links": mapof.Any{"BANDCAMP": "https://bandcamp.example/a"}}}
		require.NoError(t, s.Set(&record, "data.links.SPOTIFY", "https://spotify.example/a"))
		require.Equal(t, mapof.Any{"links": mapof.Any{"BANDCAMP": "https://bandcamp.example/a", "SPOTIFY": "https://spotify.example/a"}}, record.Data)
	})

	t.Run("a nested map decoded from JSON keeps its type and other keys", func(t *testing.T) {
		record := nestedRecord{Data: mapof.Any{"links": map[string]any{"BANDCAMP": "https://bandcamp.example/a"}}}
		require.NoError(t, s.Set(&record, "data.links.SPOTIFY", "https://spotify.example/a"))
		require.Equal(t, mapof.Any{"links": map[string]any{"BANDCAMP": "https://bandcamp.example/a", "SPOTIFY": "https://spotify.example/a"}}, record.Data)
	})
}

// toJSON encodes a value, so that containers of different Go types can be compared by content
func toJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return string(encoded)
}

// TestSet_Nested_MapofObjectOfString pins nested writes into a mapof.Object whose children are
// mapof.String values: a missing child is created as the parent's own child type
func TestSet_Nested_MapofObjectOfString(t *testing.T) {

	value := mapof.NewObject[mapof.String]()

	s := schema.New(schema.Object{
		Properties: schema.ElementMap{
			"key1": schema.Object{Properties: schema.ElementMap{"subkey1": schema.String{}, "subkey2": schema.String{}}},
			"key2": schema.Object{Properties: schema.ElementMap{"subkey1": schema.String{}, "subkey2": schema.String{}}},
		},
	})

	require.NoError(t, s.Set(&value, "key1.subkey1", "subvalue.1.1"))
	require.NoError(t, s.Set(&value, "key1.subkey2", "subvalue.1.2"))
	require.NoError(t, s.Set(&value, "key2.subkey1", "subvalue.2.1"))
	require.NoError(t, s.Set(&value, "key2.subkey2", "subvalue.2.2"))

	require.Equal(t, mapof.Object[mapof.String]{
		"key1": {"subkey1": "subvalue.1.1", "subkey2": "subvalue.1.2"},
		"key2": {"subkey1": "subvalue.2.1", "subkey2": "subvalue.2.2"},
	}, value)

	result, err := s.Get(&value, "key2.subkey1")
	require.NoError(t, err)
	require.Equal(t, "subvalue.2.1", result)
}

// TestSet_Nested_MapofTemplateChild pins that a mapof.Template's child map is a plain mapof.Any,
// reused by later writes, and that its strings are not compiled as templates
func TestSet_Nested_MapofTemplateChild(t *testing.T) {

	s := schema.New(schema.Object{Properties: schema.ElementMap{
		"validator": schema.String{},
		"child":     schema.Object{Properties: schema.ElementMap{"name": schema.String{}, "title": schema.String{}}},
	}})

	var value mapof.Template

	require.NoError(t, s.Set(&value, "validator", "/v?id={{.ID}}"))
	require.NoError(t, s.Set(&value, "child.name", "{{.Name}}"))
	require.NoError(t, s.Set(&value, "child.title", "Captain"))
	require.Equal(t, mapof.Any{"name": "{{.Name}}", "title": "Captain"}, value["child"])

	require.Error(t, s.Set(&value, "missing.name", "x"))
	require.Error(t, s.Set(&value, "child.nope", "x"))
}

// TestSet_Nested_AnyIndexOverflow pins that an index too large for an int is refused before
// anything is allocated
func TestSet_Nested_AnyIndexOverflow(t *testing.T) {

	s := nestedSchema("any")
	object := mapof.Any{"child": &sliceof.Any{}}

	require.Error(t, s.Set(&object, "child.99999999999999999999999", "v"))
	require.Error(t, s.Set(&object, "child.99999999999999999999999.id", "v"))
	require.Equal(t, mapof.Any{"child": &sliceof.Any{}}, object)
}
