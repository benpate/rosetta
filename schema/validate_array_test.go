package schema

import (
	"reflect"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// bsonArray stands in for the mongo driver's primitive.A, a named []any
type bsonArray []any

// Subtest names for the two kinds of plain list
const (
	namedList = "named []any"
	plainList = "plain []any"
)

// tagsSchema describes a document whose data.tags holds at most 3 strings of at most 5 runes
func tagsSchema() Schema {
	return New(Object{
		Properties: map[string]Element{
			"data": Object{
				Properties: map[string]Element{
					"tags": Array{MaxLength: 3, Items: String{MaxLength: 5}},
				},
			},
		},
	})
}

// tagsDocument returns a document holding this list under data.tags
func tagsDocument(tags any) mapof.Any {
	return mapof.Any{"data": mapof.Any{"tags": tags}}
}

func TestValidateArray_ListsHeldByValue(t *testing.T) {

	// Normalize a document whose data.tags holds the given list
	normalize := func(tags any) (mapof.Any, []string, error) {
		document := tagsDocument(tags)
		rewrites, err := tagsSchema().Normalize(&document)
		return document, rewrites, err
	}

	// A valid list passes unchanged, and keeps its own type
	unchanged := func(name string, tags any) {
		t.Run("unchanged: "+name, func(t *testing.T) {
			document, rewrites, err := normalize(tags)
			require.NoError(t, err)
			require.Empty(t, rewrites)
			require.Equal(t, tagsDocument(tags), document)
		})
	}

	// An over-long item is truncated, and the list is stored back as its own type
	rewritten := func(name string, tags any, expected any) {
		t.Run("rewritten: "+name, func(t *testing.T) {
			document, rewrites, err := normalize(tags)
			require.NoError(t, err)
			require.Equal(t, []string{"data.tags.1"}, rewrites)
			require.Equal(t, tagsDocument(expected), document)
		})
	}

	// A list longer than MaxLength is a validation error, not an internal one
	tooLong := func(name string, tags any) {
		t.Run("too long: "+name, func(t *testing.T) {
			_, _, err := normalize(tags)
			require.Error(t, err)
			code := derp.ErrorCode(err)
			require.True(t, derp.IsValidationError(err), "expected 422, got %d", code)
		})
	}

	unchanged(namedList, bsonArray{"a", "b"})
	unchanged("named []any, empty", bsonArray{})
	unchanged(plainList, []any{"a", "b"})
	unchanged("plain []any, empty", []any{})
	unchanged("sliceof.String", sliceof.String{"a", "b"})
	unchanged("sliceof.Any", sliceof.Any{"a", "b"})

	rewritten(namedList, bsonArray{"ok", "toolong"}, bsonArray{"ok", "toolo"})
	rewritten(plainList, []any{"ok", "toolong"}, []any{"ok", "toolo"})
	rewritten("sliceof.String", sliceof.String{"ok", "toolong"}, sliceof.String{"ok", "toolo"})

	tooLong(namedList, bsonArray{"a", "b", "c", "d"})
	tooLong(plainList, []any{"a", "b", "c", "d"})
}

func TestValidateArray_ListOfMapsHeldByValue(t *testing.T) {

	schema := New(Object{
		Properties: map[string]Element{
			"people": Array{Items: Object{
				Properties: map[string]Element{"name": String{MaxLength: 4}},
			}},
		},
	})

	// A []mapof.Any is validated through a sliceof.Object[mapof.Any], and stored back as itself
	document := mapof.Any{"people": []mapof.Any{{"name": "Sarah"}, {"name": "Kyle"}}}
	rewrites, err := schema.Normalize(&document)

	require.NoError(t, err)
	require.Equal(t, []string{"people.0.name"}, rewrites)
	require.Equal(t, mapof.Any{"people": []mapof.Any{{"name": "Sara"}, {"name": "Kyle"}}}, document)
}

func TestValidateArray_UnconvertibleListHeldByValue(t *testing.T) {

	// A []string converts to no rosetta container, so it still fails as before (BUG-298 K5)
	document := tagsDocument([]string{"a"})
	_, err := tagsSchema().Normalize(&document)

	require.Error(t, err)
	require.True(t, derp.IsInternalServerError(err))
	require.Equal(t, tagsDocument([]string{"a"}), document)
}

func TestAddressableCopy(t *testing.T) {

	// The name of the type addressableCopy returns for a value held by value
	copyType := func(value any) string {
		return reflect.TypeOf(addressableCopy(reflect.ValueOf(value))).String()
	}

	// Plain lists and maps become the rosetta type with the same layout
	require.Equal(t, reflect.TypeFor[*sliceof.Any]().String(), copyType([]any{"a"}))
	require.Equal(t, reflect.TypeFor[*sliceof.Any]().String(), copyType(bsonArray{}))
	listOfMaps := reflect.TypeFor[*sliceof.Object[mapof.Any]]()
	require.Equal(t, listOfMaps.String(), copyType([]mapof.Any{}))
	require.Equal(t, reflect.TypeFor[*mapof.Any]().String(), copyType(map[string]any{}))

	// Types with their own accessors keep their type
	require.Equal(t, reflect.TypeFor[*sliceof.String]().String(), copyType(sliceof.String{"a"}))
	require.Equal(t, reflect.TypeFor[*sliceof.Any]().String(), copyType(sliceof.Any{"a"}))
	require.Equal(t, reflect.TypeFor[*mapof.Any]().String(), copyType(mapof.Any{}))
	require.Equal(t, reflect.TypeFor[*testArrayA]().String(), copyType(testArrayA{"a"}))

	// Anything else is a pointer to a copy of its own type
	require.Equal(t, reflect.TypeFor[*[]string]().String(), copyType([]string{"a"}))
	require.Equal(t, reflect.TypeFor[*map[string]string]().String(), copyType(map[string]string{}))
	structType := reflect.TypeFor[*pointerGetterStruct]()
	require.Equal(t, structType.String(), copyType(pointerGetterStruct{}))
	require.Equal(t, reflect.TypeFor[*string]().String(), copyType("a"))

	// The copy is shallow, but never the original's own storage
	original := []any{"a"}
	copied, isList := addressableCopy(reflect.ValueOf(original)).(*sliceof.Any)
	require.True(t, isList)
	(*copied)[0] = "b"
	require.Equal(t, []any{"b"}, original, "a slice copy shares its backing array")
	*copied = append(*copied, "c")
	require.Equal(t, []any{"b"}, original, "growing the copy leaves the original's length alone")
}

func TestIsValidatable(t *testing.T) {
	require.True(t, isValidatable(&mapof.Any{}))
	require.True(t, isValidatable(&sliceof.Any{}))
	require.True(t, isValidatable(&sliceof.String{}))
	require.True(t, isValidatable(&pointerGetterStruct{}))
	require.False(t, isValidatable(nil))
	require.False(t, isValidatable("a"))
	require.False(t, isValidatable(&[]any{}))
	require.False(t, isValidatable(&[]string{}))
	require.False(t, isValidatable(map[string]any{}))
}
