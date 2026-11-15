package deep

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestCopy_KeepsTypeAndData proves that every supported value copies to an equal value of
// its own type
func TestCopy_KeepsTypeAndData(t *testing.T) {

	now := time.Now()

	values := map[string]any{
		"nil":                     nil,
		"bool":                    true,
		"float32":                 float32(1.5),
		"float64":                 2.5,
		"int":                     1,
		"int8":                    int8(2),
		"int16":                   int16(3),
		"int32":                   int32(4),
		"int64":                   int64(5),
		"uint":                    uint(6),
		"string":                  "seven",
		"time.Time":               now,
		"map[string]any":          map[string]any{"a": 1, "b": nil, "c": []any{nil, "d"}},
		"map[string]bool":         map[string]bool{"a": true},
		"map[string]int":          map[string]int{"a": 1},
		"map[string]int8":         map[string]int8{"a": 1},
		"map[string]int16":        map[string]int16{"a": 1},
		"map[string]int32":        map[string]int32{"a": 1},
		"map[string]int64":        map[string]int64{"a": 1},
		"map[string]float32":      map[string]float32{"a": 1},
		"map[string]float64":      map[string]float64{"a": 1},
		"map[string]string":       map[string]string{"a": "b"},
		"map[string]mapof.Any":    map[string]mapof.Any{"a": {"b": 1}},
		"map[string][]string":     map[string][]string{"a": {"b"}},
		"map[int]any":             map[int]any{1: []any{"a"}},
		"mapof.Any":               mapof.Any{"a": mapof.Any{"b": 1}},
		"mapof.Bool":              mapof.Bool{"a": true},
		"mapof.Float":             mapof.Float{"a": 1},
		"mapof.Int":               mapof.Int{"a": 1},
		"mapof.Int64":             mapof.Int64{"a": 1},
		"mapof.String":            mapof.String{"a": "b"},
		"[]any":                   []any{1, "a", nil, map[string]any{"b": 2}},
		"[]bool":                  []bool{true},
		"[]byte":                  []byte("abc"),
		"[]float32":               []float32{1},
		"[]float64":               []float64{1},
		"[]int":                   []int{1},
		"[]int8":                  []int8{1},
		"[]int16":                 []int16{1},
		"[]int32":                 []int32{1},
		"[]int64":                 []int64{1},
		"[]uint":                  []uint{1},
		"[]string":                []string{"a"},
		"[]map[string]any":        []map[string]any{{"a": 1}},
		"[]map[string]string":     []map[string]string{{"a": "b"}},
		"[]mapof.Any":             []mapof.Any{{"a": 1}},
		"[][]any":                 [][]any{{1}, nil},
		"sliceof.Any":             sliceof.Any{1, "a"},
		"sliceof.Float":           sliceof.Float{1},
		"sliceof.Int":             sliceof.Int{1},
		"sliceof.String":          sliceof.String{"a"},
		"sliceof.MapOfAny":        sliceof.MapOfAny{{"a": 1}},
		"sliceof.MapOfString":     sliceof.MapOfString{{"a": "b"}},
		"database array":          testArray{"a", testDocument{"b": testDateTime(1)}},
		"database document":       testDocument{"_id": testObjectID{1, 2, 3}, "list": testArray{1}},
		"database date":           testDateTime(1700000000000),
		"database ID":             testObjectID{1, 2, 3},
		"database decimal":        testDecimal{h: 1, l: 2},
		"database timestamp":      testTimestamp{T: 1, I: 2},
		"database regex":          testRegex{Pattern: "a", Options: "i"},
		"database null":           testNull{},
		"slice of simple structs": []testTimestamp{{T: 1}},
		"array of simple structs": [1]testRegex{{Pattern: "a"}},
		"map of simple structs":   map[string]testDecimal{"a": {h: 1}},
		"[]testObjectID":          []testObjectID{{1}, {2}},
		"*sliceof.Any":            &sliceof.Any{1, map[string]any{"a": nil}},
		"nil *sliceof.Any":        (*sliceof.Any)(nil),
		"[2]any":                  [2]any{map[string]any{"a": 1}, nil},
		"empty map":               map[string]any{},
		"empty slice":             []any{},
		"nil map[string]any":      map[string]any(nil),
		"nil mapof.String":        mapof.String(nil),
		"nil []any":               []any(nil),
		"nil []string":            []string(nil),
		"nil database array":      testArray(nil),
		"nil database doc":        testDocument(nil),
	}

	for name, original := range values {
		t.Run(name, func(t *testing.T) {

			result, err := Copy(original)

			// ObjectsAreEqual compares types as well as data, and tells nil from empty
			require.NoError(t, err)
			require.Equal(t, original, result)
		})
	}
}

// TestCopy_SharesNothing proves that changing a copy never changes the original
func TestCopy_SharesNothing(t *testing.T) {

	t.Run("plain containers", func(t *testing.T) {

		original := map[string]any{"to": []string{"a"}, "object": map[string]any{"tag": []any{"x"}}}

		result, err := Copy(original)
		require.NoError(t, err)

		result["to"].([]string)[0] = "changed"
		result["object"].(map[string]any)["tag"].([]any)[0] = "changed"
		result["object"].(map[string]any)["new"] = true

		expected := map[string]any{"to": []string{"a"}, "object": map[string]any{"tag": []any{"x"}}}
		require.Equal(t, expected, original)
	})

	t.Run("named containers", func(t *testing.T) {

		original := sliceof.MapOfAny{{"a": mapof.String{"b": "c"}}}

		result, err := Copy(original)
		require.NoError(t, err)

		result[0]["a"].(mapof.String)["b"] = "changed"

		require.Equal(t, sliceof.MapOfAny{{"a": mapof.String{"b": "c"}}}, original)
	})

	t.Run("reflected containers", func(t *testing.T) {

		original := testDocument{"list": testArray{testDocument{"a": 1}}, "bytes": []uint16{1}}

		result, err := Copy(original)
		require.NoError(t, err)

		result["list"].(testArray)[0].(testDocument)["a"] = 2
		result["bytes"].([]uint16)[0] = 2

		expected := testDocument{"list": testArray{testDocument{"a": 1}}, "bytes": []uint16{1}}
		require.Equal(t, expected, original)
	})

	t.Run("array of maps", func(t *testing.T) {

		original := [1]map[string]any{{"a": 1}}

		result, err := Copy(original)
		require.NoError(t, err)

		result[0]["a"] = 2

		require.Equal(t, 1, original[0]["a"])
	})
}

// TestCopy_PointerToSlice proves that a *sliceof.Any is copied to a new pointer, which
// shares nothing with the original
func TestCopy_PointerToSlice(t *testing.T) {

	original := map[string]any{"list": &sliceof.Any{map[string]any{"a": 1}}}

	result, err := Copy(original)
	require.NoError(t, err)

	copied := result["list"].(*sliceof.Any)
	require.NotSame(t, original["list"], copied)

	(*copied)[0].(map[string]any)["a"] = 2
	*copied = append(*copied, "new")

	require.Equal(t, &sliceof.Any{map[string]any{"a": 1}}, original["list"])
}

// TestCopy_MaxDepth proves that a walk stops at the maximum depth, which is what stops a
// value that contains itself
func TestCopy_MaxDepth(t *testing.T) {

	t.Run("map that contains itself", func(t *testing.T) {

		original := map[string]any{}
		original["self"] = original

		result, err := Copy(original)

		require.Error(t, err)
		require.Nil(t, result)
		require.True(t, derp.IsInternalServerError(err))
	})

	t.Run("slice that contains itself", func(t *testing.T) {

		original := []any{nil}
		original[0] = original

		_, err := Copy(original)
		require.Error(t, err)
	})

	t.Run("reflected containers that contain themselves", func(t *testing.T) {

		document := testDocument{}
		document["list"] = testArray{document}

		_, err := Copy(document)
		require.Error(t, err)

		array := [1]any{}
		slice := testArray{nil}
		array[0] = slice
		slice[0] = array

		_, err = Copy(slice)
		require.Error(t, err)
	})

	t.Run("pointer that leads to itself", func(t *testing.T) {

		original := &sliceof.Any{nil}
		(*original)[0] = original

		_, err := Copy(original)
		require.Error(t, err)
	})

	t.Run("the limit allows what a JSON decoder allows", func(t *testing.T) {

		var original any = "deepest"

		for range maxDepth {
			original = []any{original}
		}

		result, err := Copy(original)
		require.NoError(t, err)
		require.Equal(t, original, result)

		_, err = Copy([]any{original})
		require.Error(t, err)
	})

	t.Run("the error does not grow with the depth", func(t *testing.T) {

		original := map[string]any{}
		original["self"] = original

		_, err := Copy(original)

		require.Error(t, err)
		require.Less(t, len(errorText(t, err)), 2_000)
	})
}

// TestCopy_UnsupportedTypes proves that a value Copy cannot copy is an error, at any depth,
// and that no partial copy comes back with it
func TestCopy_UnsupportedTypes(t *testing.T) {

	number := 1

	values := map[string]any{
		"struct":                     testStruct{},
		"pointer":                    &number,
		"pointer to slice":           &[]any{1},
		"pointer in sliceof pointer": &sliceof.Any{&number},
		"struct with a pointer":      struct{ P *int }{},
		"struct with an interface":   struct{ V any }{},
		"nil pointer":                (*int)(nil),
		"function":                   func() {},
		"channel":                    make(chan int),
		"struct in map":              map[string]any{"a": 1, "b": testStruct{}},
		"struct in slice":            []any{1, testStruct{}},
		"struct in named map":        mapof.Any{"a": testStruct{}},
		"struct in typed map":        map[string]mapof.Any{"a": {"b": testStruct{}}},
		"struct in slice of maps":    sliceof.MapOfAny{{"a": testStruct{}}},
		"struct in reflected map":    testDocument{"a": testStruct{}},
		"struct in reflected slice":  testArray{testStruct{}},
		"struct in array":            [1]any{testStruct{}},
		"slice of structs":           []testStruct{{}},
		"array of structs":           [1]testStruct{},
		"map of structs":             map[string]testStruct{"a": {}},
	}

	for name, original := range values {
		t.Run(name, func(t *testing.T) {

			result, err := Copy(original)

			require.Error(t, err)
			require.Nil(t, result)
			require.True(t, derp.IsInternalServerError(err), "an unsupported type is an Internal error")
		})
	}
}

// TestCopy_ErrorNamesTypeNotValue proves that an unsupported value is never written into
// the error
func TestCopy_ErrorNamesTypeNotValue(t *testing.T) {

	_, err := Copy(map[string]any{"bcc": testStruct{Secret: "hunter2-recipient", Hidden: map[string]any{"bcc": 1}}})

	require.Error(t, err)
	require.NotContains(t, errorText(t, err), "hunter2-recipient")
	require.Contains(t, errorText(t, err), "deep.testStruct")
}

// TestCopy_Copier proves that a value which copies itself is asked to, wherever it is held
func TestCopy_Copier(t *testing.T) {

	t.Run("called at every depth", func(t *testing.T) {

		copies := 0
		copier := testCopier{Name: "a", Copies: &copies}

		original := map[string]any{"direct": copier, "listed": []any{copier}, "reflected": testArray{copier}}

		result, err := Copy(original)

		require.NoError(t, err)
		require.Equal(t, original, result)
		require.Equal(t, 3, copies)
	})

	t.Run("typed containers", func(t *testing.T) {

		copies := 0
		original := map[string]testCopier{"a": {Name: "a", Copies: &copies}}

		result, err := Copy(original)

		require.NoError(t, err)
		require.Equal(t, original, result)
		require.Equal(t, 1, copies)
	})

	t.Run("is told its depth, and continues from it", func(t *testing.T) {

		box := &testBox{Child: map[string]any{"a": []any{1}}}
		original := map[string]any{"list": []any{box}}

		result, err := Copy(original)
		require.NoError(t, err)
		require.Equal(t, 2, box.Depth)

		// The copy is a new box, and shares nothing with the original
		copied := result["list"].([]any)[0].(*testBox)
		require.NotSame(t, box, copied)

		copied.Child.(map[string]any)["a"].([]any)[0] = 2
		require.Equal(t, map[string]any{"a": []any{1}}, box.Child)
	})

	t.Run("a cycle through a Copier ends in an error", func(t *testing.T) {

		box := &testBox{}
		box.Child = map[string]any{"box": box}

		result, err := Copy(box)

		require.Error(t, err)
		require.Nil(t, result)
		require.Less(t, len(errorText(t, err)), 2_000)
	})

	t.Run("error returns nil", func(t *testing.T) {

		var original any = testFailingCopier{}

		result, err := Copy(original)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("error returns the zero value of a concrete type", func(t *testing.T) {

		result, err := Copy(testFailingCopier{})

		require.Error(t, err)
		require.Equal(t, testFailingCopier{}, result)
	})

	t.Run("wrong type in a typed slice", func(t *testing.T) {

		result, err := Copy([]testWrongCopier{{}})

		require.Error(t, err)
		require.Nil(t, result)
		require.True(t, derp.IsInternalServerError(err))
	})

	t.Run("wrong type in a typed map", func(t *testing.T) {

		result, err := Copy(map[string]testWrongCopier{"a": {}})

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("wrong type in an array", func(t *testing.T) {

		result, err := Copy([1]testWrongCopier{})

		require.Error(t, err)
		require.Equal(t, [1]testWrongCopier{}, result)
	})

	t.Run("any type fits a container of any", func(t *testing.T) {

		result, err := Copy([]any{testWrongCopier{}})

		require.NoError(t, err)
		require.Equal(t, []any{"not a testWrongCopier"}, result)
	})
}

// TestCopy_ReturnsCallersType proves that Copy returns the copy as the type it was given
func TestCopy_ReturnsCallersType(t *testing.T) {

	t.Run("named map", func(t *testing.T) {

		original := mapof.Any{"a": []any{1}}

		result, err := Copy(original)

		require.NoError(t, err)
		require.Equal(t, original, result)
	})

	t.Run("simple value", func(t *testing.T) {

		result, err := Copy("a")

		require.NoError(t, err)
		require.Equal(t, "a", result)
	})

	t.Run("nil interface", func(t *testing.T) {

		var original any

		result, err := Copy(original)

		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("unsupported type returns the zero value", func(t *testing.T) {

		result, err := Copy(testStruct{Secret: "a"})

		require.Error(t, err)
		require.Empty(t, result)
	})

	t.Run("wrong type returns the zero value", func(t *testing.T) {

		result, err := Copy(testWrongCopier{})

		require.Error(t, err)
		require.True(t, derp.IsInternalServerError(err))
		require.Equal(t, testWrongCopier{}, result)
	})
}

// FuzzCopy_JSON proves that any decoded JSON document copies to an equal document, and that
// Remove then leaves the original untouched
func FuzzCopy_JSON(f *testing.F) {

	f.Add(`{"to":["a"],"bcc":"b","object":{"bcc":["c"],"tag":[{"bcc":null}]}}`)
	f.Add(`[1,2.5,"three",true,null,{"bcc":{}}]`)
	f.Add(`null`)
	f.Add(`"bcc"`)
	f.Add(`{"":{"":[]}}`)
	f.Add(`[[[[[[{"bcc":1}]]]]]]`)

	f.Fuzz(func(t *testing.T, input string) {

		var original any

		if err := json.Unmarshal([]byte(input), &original); err != nil {
			t.Skip()
		}

		// The copy holds the same data as the original
		result, err := Copy(original)
		require.NoError(t, err)
		require.Equal(t, original, result)

		// Removing from the copy removes the name at every depth, and changes only the copy
		require.NoError(t, Remove(result, "bcc"))
		require.False(t, hasKey(result, "bcc"))

		var untouched any
		require.NoError(t, json.Unmarshal([]byte(input), &untouched))
		require.Equal(t, untouched, original)
	})
}

// hasKey returns TRUE if a decoded JSON value holds the named key at any depth
func hasKey(value any, name string) bool {

	switch typed := value.(type) {

	case map[string]any:

		if _, found := typed[name]; found {
			return true
		}

		for _, child := range typed {
			if hasKey(child, name) {
				return true
			}
		}

	case []any:

		for _, child := range typed {
			if hasKey(child, name) {
				return true
			}
		}
	}

	return false
}

// errorText returns an error and everything it wraps as JSON, details included
func errorText(t *testing.T, err error) string {

	t.Helper()

	result, marshalErr := json.Marshal(err)
	require.NoError(t, marshalErr)

	return string(result)
}
