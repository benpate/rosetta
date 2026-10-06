package schema_test

import (
	"encoding/json"
	"reflect"
	"runtime"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

/******************************************
 * Follow-up defects in Get and Set (BUG-234 §11)
 *
 * Each TestFollowup_* test describes the behavior a fix produced, and failed before it. Each
 * TestFollowupGuard_* test pins neighbouring behavior that the fix had to keep.
 ******************************************/

// followupRecord is a struct whose fields are reached through GetPointer, as Emissary's models are
type followupRecord struct {
	Name string
	Rank int
	Data mapof.Any
	Tags sliceof.String
}

// GetPointer implements the schema PointerGetter interface
func (record *followupRecord) GetPointer(name string) (any, bool) {

	switch name {
	case "name":
		return &record.Name, true
	case "rank":
		return &record.Rank, true
	case "data":
		return &record.Data, true
	case "tags":
		return &record.Tags, true
	}

	return nil, false
}

// followupSchemas are the schemas the follow-up tests read and write through
var (
	followupWildcard = schema.New(schema.Object{Wildcard: schema.Any{}})
	followupArrays   = schema.New(schema.Object{Properties: schema.ElementMap{
		"tags": schema.Array{Items: schema.String{}},
		"objs": schema.Array{Items: schema.Object{Wildcard: schema.Any{}}},
	}})
	followupStrings = schema.New(schema.Object{
		Wildcard: schema.Any{},
		Properties: schema.ElementMap{
			"name":  schema.String{},
			"child": schema.Object{Properties: schema.ElementMap{"name": schema.String{}}},
		},
	})
)

/******************************************
 * Defect 1: Get could not read a list held by value
 ******************************************/

// TestFollowup_GetReadsListsHeldByValue requires Get to read an item of a list stored by value,
// as it already does for a list stored by pointer
func TestFollowup_GetReadsListsHeldByValue(t *testing.T) {

	// BUG-234 follow-up 1: a value list's GetPointer has a pointer receiver, so Get found none
	value, err := followupWildcard.Get(&mapof.Any{"list": sliceof.Any{"a", "b"}}, "list.1")
	require.NoError(t, err)
	require.Equal(t, "b", value)

	objects := mapof.Any{"objs": sliceof.Object[mapof.Any]{{"x": 1}}}
	value, err = followupArrays.Get(&objects, "objs.0.x")
	require.NoError(t, err)
	require.Equal(t, 1, value)

	value, err = followupArrays.Get(&mapof.Any{"tags": []any{"a", "b"}}, "tags.1")
	require.NoError(t, err)
	require.Equal(t, "b", value)
}

// TestFollowupGuard_GetLists pins the list reads that worked before the fix, and the errors kept
func TestFollowupGuard_GetLists(t *testing.T) {

	t.Run("a list held by pointer", func(t *testing.T) {
		value, err := followupArrays.Get(&mapof.Any{"tags": &sliceof.Any{"a", "b"}}, "tags.1")
		require.NoError(t, err)
		require.Equal(t, "b", value)

		objects := mapof.Any{"objs": &sliceof.Object[mapof.Any]{{"x": 1}}}
		value, err = followupArrays.Get(&objects, "objs.0.x")
		require.NoError(t, err)
		require.Equal(t, 1, value)
	})

	t.Run("a list of strings held by value, read through its typed getter", func(t *testing.T) {
		for _, list := range []any{sliceof.Any{"a", "b"}, sliceof.String{"a", "b"}} {
			value, err := followupArrays.Get(&mapof.Any{"tags": list}, "tags.1")
			require.NoError(t, err, "%T", list)
			require.Equal(t, "b", value, "%T", list)
		}
	})

	t.Run("a whole list is returned as it is stored", func(t *testing.T) {
		stored := sliceof.Any{"a"}
		value, err := followupWildcard.Get(&mapof.Any{"list": stored}, "list")
		require.NoError(t, err)
		require.Equal(t, stored, value)
	})

	t.Run("a value with no getter interfaces is still an error", func(t *testing.T) {
		_, err := followupWildcard.Get(&struct{}{}, "x")
		require.Error(t, err)

		_, err = followupWildcard.Get(&mapof.Any{"list": 7}, "list.0")
		require.Error(t, err)
	})
}

/******************************************
 * Defect 2: Get returned a scalar list item as a pointer
 ******************************************/

// TestFollowup_GetReturnsScalarListItems requires Get to return a scalar item of a list beneath
// Any as its value, the way it returns a scalar from a map
func TestFollowup_GetReturnsScalarListItems(t *testing.T) {

	// BUG-234 follow-up 2: sliceof.Any.GetPointer returns a pointer to a copy, which Get returned
	object := mapof.Any{"list": &sliceof.Any{"a", 2, true}}

	for path, expected := range map[string]any{"list.0": "a", "list.1": 2, "list.2": true} {
		value, err := followupWildcard.Get(&object, path)
		require.NoError(t, err, path)
		require.Equal(t, expected, value, path)
	}
}

// TestFollowupGuard_GetBeneathAny pins what Get returns beneath Any, other than scalar list items.
// A change to any of these changes what callers receive.
func TestFollowupGuard_GetBeneathAny(t *testing.T) {

	t.Run("a scalar in a map is returned as its value", func(t *testing.T) {
		object := mapof.Any{"child": mapof.Any{"name": "v"}}
		value, err := followupWildcard.Get(&object, "child.name")
		require.NoError(t, err)
		require.Equal(t, "v", value)
	})

	t.Run("a map is returned as it is stored", func(t *testing.T) {
		value, err := followupWildcard.Get(&mapof.Any{"child": mapof.Any{"name": "v"}}, "child")
		require.NoError(t, err)
		require.Equal(t, mapof.Any{"name": "v"}, value)
	})

	t.Run("a map in a list is returned by pointer, and can be read through", func(t *testing.T) {
		object := mapof.Any{"list": &sliceof.Any{mapof.Any{"k": "v"}}}

		value, err := followupWildcard.Get(&object, "list.0")
		require.NoError(t, err)
		require.Equal(t, &mapof.Any{"k": "v"}, value)

		value, err = followupWildcard.Get(&object, "list.0.k")
		require.NoError(t, err)
		require.Equal(t, "v", value)
	})

	t.Run("a struct's fields are returned by pointer, and can be read through", func(t *testing.T) {
		record := followupRecord{Name: "n", Rank: 3}
		record.Data = mapof.Any{"k": "v"}
		record.Tags = sliceof.String{"t"}

		cases := map[string]any{
			"name": &record.Name,
			"rank": &record.Rank,
			"data": &record.Data,
			"tags": &record.Tags,
		}

		for path, expected := range cases {
			value, err := followupWildcard.Get(&record, path)
			require.NoError(t, err, path)
			require.Same(t, expected, value, path)
		}

		value, err := followupWildcard.Get(&record, "data.k")
		require.NoError(t, err)
		require.Equal(t, "v", value)
	})

	t.Run("every value encodes to the same JSON as the value it points to", func(t *testing.T) {
		record := followupRecord{Name: "n", Rank: 3}
		object := mapof.Any{"list": &sliceof.Any{"a", 2}}

		for _, read := range []struct {
			object any
			path   string
			JSON   string
		}{
			{object: &record, path: "name", JSON: `"n"`},
			{object: &record, path: "rank", JSON: `3`},
			{object: &object, path: "list.0", JSON: `"a"`},
			{object: &object, path: "list.1", JSON: `2`},
		} {
			value, err := followupWildcard.Get(read.object, read.path)
			require.NoError(t, err, read.path)
			encoded, err := json.Marshal(value)
			require.NoError(t, err, read.path)
			require.Equal(t, read.JSON, string(encoded), read.path)
		}
	})
}

/******************************************
 * Defect 3: an empty path segment was stored as a literal key
 ******************************************/

// TestFollowup_SetRejectsEmptySegments requires Set to refuse a path with an empty segment, and
// to store nothing, instead of writing a key that Get cannot read back
func TestFollowup_SetRejectsEmptySegments(t *testing.T) {

	// BUG-234 follow-up 3: "name." stored the key "name.", and ".name" stored under the key ""
	for _, path := range []string{"name.", ".name", "child.name.", ".", "child..name"} {
		object := mapof.Any{}
		require.Error(t, followupStrings.Set(&object, path, "v"), path)
		require.Equal(t, mapof.Any{}, object, path)
	}
}

// TestFollowupGuard_SetPaths pins the paths Set accepts, and the empty path, which it refused
// before the fix as well
func TestFollowupGuard_SetPaths(t *testing.T) {

	t.Run("ordinary, unusual, and nested keys are stored", func(t *testing.T) {
		cases := map[string]mapof.Any{
			"name":       {"name": "v"},
			"a b":        {"a b": "v"},
			"ü":          {"ü": "v"},
			"0":          {"0": "v"},
			"x.y":        {"x": mapof.Any{"y": "v"}},
			"child.name": {"child": mapof.Any{"name": "v"}},
		}

		for path, expected := range cases {
			object := mapof.Any{}
			require.NoError(t, followupStrings.Set(&object, path, "v"), path)
			require.Equal(t, expected, object, path)
		}
	})

	t.Run("an empty path is refused", func(t *testing.T) {
		object := mapof.Any{}
		require.Error(t, followupStrings.Set(&object, "", "v"))
		require.Equal(t, mapof.Any{}, object)
	})

	t.Run("a whole child is still written back through an empty path", func(t *testing.T) {
		child := mapof.Any{"name": "old"}
		element := schema.Object{Properties: schema.ElementMap{"name": schema.String{}}}
		require.NoError(t, schema.SetProperty(element, &child, "", mapof.Any{"name": "new"}))
		require.Equal(t, mapof.Any{"name": "new"}, child)
	})
}

/******************************************
 * Defect 4: Get of a missing list index grew the list
 ******************************************/

// followupGrowthReads are reads past the end of a one-item list, through every kind of element
var followupGrowthReads = []struct {
	path string
	list func() any
}{
	{path: "tags.5", list: func() any { return &sliceof.Any{"a"} }},
	{path: "objs.5", list: func() any { return &sliceof.Object[mapof.Any]{{"x": 1}} }},
	{path: "objs.5.x", list: func() any { return &sliceof.Object[mapof.Any]{{"x": 1}} }},
	{path: "other.5", list: func() any { return &sliceof.Any{"a"} }},
	{path: "other.5.x", list: func() any { return &sliceof.Any{"a"} }},
	{path: "other.next", list: func() any { return &sliceof.Object[mapof.Any]{{"x": 1}} }},
}

// followupGrowthSchema reads lists through a typed Array, an object Array, and Any
var followupGrowthSchema = schema.New(schema.Object{
	Wildcard: schema.Any{},
	Properties: schema.ElementMap{
		"tags": schema.Array{Items: schema.String{}},
		"objs": schema.Array{Items: schema.Object{Wildcard: schema.Any{}}},
	},
})

// TestFollowup_GetDoesNotGrowLists requires that reading past the end of a list leaves the list
// unchanged
func TestFollowup_GetDoesNotGrowLists(t *testing.T) {

	// BUG-234 follow-up 4: sliceof's GetPointer grows the slice to reach an index, even to read
	for _, read := range followupGrowthReads {
		list := read.list()
		object := mapof.Any{"tags": list, "objs": list, "other": list}
		_, _ = followupGrowthSchema.Get(&object, read.path) // Only the list's length matters here
		require.Equal(t, 1, reflect.ValueOf(list).Elem().Len(), read.path)
	}
}

// TestFollowup_GetPastTheEndReturnsZero requires a read past the end of a list to return the
// zero value of the element at the end of the path, with no error and no allocation
func TestFollowup_GetPastTheEndReturnsZero(t *testing.T) {

	// BUG-234 follow-up 4: these used to return nil or an error, depending on the element
	expected := map[string]any{
		"tags.5":     "",
		"objs.5":     nil,
		"objs.5.x":   nil,
		"other.5":    nil,
		"other.5.x":  nil,
		"other.next": nil,
	}

	for _, read := range followupGrowthReads {
		object := mapof.Any{"tags": read.list(), "objs": read.list(), "other": read.list()}
		value, err := followupGrowthSchema.Get(&object, read.path)
		require.NoError(t, err, read.path)
		require.Equal(t, expected[read.path], value, read.path)
	}

	t.Run("a far index allocates no more than a read inside the list", func(t *testing.T) {
		object := mapof.Any{"other": &sliceof.Any{"a"}}
		near := allocatedBytes(func() { _, _ = followupGrowthSchema.Get(&object, "other.0") })
		far := allocatedBytes(func() { _, _ = followupGrowthSchema.Get(&object, "other.10000000") })
		require.LessOrEqual(t, far, near+1024, "growing the list allocated %d bytes", far)
	})

	t.Run("a path the schema lacks beneath the end is still an error", func(t *testing.T) {
		strict := schema.New(schema.Object{Properties: schema.ElementMap{
			"tags": schema.Array{Items: schema.String{}},
		}})
		_, err := strict.Get(&mapof.Any{"tags": &sliceof.Any{}}, "tags.5.nope")
		require.Error(t, err)
	})
}

// TestFollowupGuard_GetInsideTheList pins the reads that a past-the-end check must not catch
func TestFollowupGuard_GetInsideTheList(t *testing.T) {

	t.Run("an index inside the list still reads its item", func(t *testing.T) {
		object := mapof.Any{"tags": &sliceof.Any{"a", "b"}}
		value, err := followupGrowthSchema.Get(&object, "tags.1")
		require.NoError(t, err)
		require.Equal(t, "b", value)
	})

	t.Run("last still reads the final item", func(t *testing.T) {
		object := mapof.Any{"other": &sliceof.Any{"a", "b"}}
		value, err := followupGrowthSchema.Get(&object, "other.last")
		require.NoError(t, err)
		require.Equal(t, "b", value)
	})

	t.Run("a missing map key is still an error", func(t *testing.T) {
		_, err := followupGrowthSchema.Get(&mapof.Any{}, "missing")
		require.Error(t, err)
	})
}

// allocatedBytes returns the bytes allocated while running a function
func allocatedBytes(run func()) uint64 {

	var before, after runtime.MemStats

	runtime.ReadMemStats(&before)
	run()
	runtime.ReadMemStats(&after)

	return after.TotalAlloc - before.TotalAlloc
}
