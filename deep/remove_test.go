package deep

import (
	"testing"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestRemove_EveryContainer proves that a name is removed from every kind of container, and
// that every other property stays
func TestRemove_EveryContainer(t *testing.T) {

	// Each closure returns a fresh value that holds "bcc", and the value expected afterwards
	values := map[string]func() (any, any){
		"map[string]any": func() (any, any) {
			return map[string]any{"bcc": 1, "to": 2}, map[string]any{"to": 2}
		},
		"nested maps": func() (any, any) {
			return map[string]any{"object": map[string]any{"bcc": 1, "tag": []any{map[string]any{"bcc": 1}}}},
				map[string]any{"object": map[string]any{"tag": []any{map[string]any{}}}}
		},
		"map[string]bool":    func() (any, any) { return map[string]bool{"bcc": true, "to": true}, map[string]bool{"to": true} },
		"map[string]int":     func() (any, any) { return map[string]int{"bcc": 1, "to": 2}, map[string]int{"to": 2} },
		"map[string]int8":    func() (any, any) { return map[string]int8{"bcc": 1}, map[string]int8{} },
		"map[string]int16":   func() (any, any) { return map[string]int16{"bcc": 1}, map[string]int16{} },
		"map[string]int32":   func() (any, any) { return map[string]int32{"bcc": 1}, map[string]int32{} },
		"map[string]int64":   func() (any, any) { return map[string]int64{"bcc": 1}, map[string]int64{} },
		"map[string]float32": func() (any, any) { return map[string]float32{"bcc": 1}, map[string]float32{} },
		"map[string]float64": func() (any, any) { return map[string]float64{"bcc": 1}, map[string]float64{} },
		"map[string]string":  func() (any, any) { return map[string]string{"bcc": "a", "to": "b"}, map[string]string{"to": "b"} },
		"map[string]mapof.Any": func() (any, any) {
			return map[string]mapof.Any{"bcc": {}, "a": {"bcc": 1, "to": 2}}, map[string]mapof.Any{"a": {"to": 2}}
		},
		"mapof.Any":    func() (any, any) { return mapof.Any{"bcc": 1, "a": mapof.Any{"bcc": 1}}, mapof.Any{"a": mapof.Any{}} },
		"mapof.Bool":   func() (any, any) { return mapof.Bool{"bcc": true}, mapof.Bool{} },
		"mapof.Float":  func() (any, any) { return mapof.Float{"bcc": 1}, mapof.Float{} },
		"mapof.Int":    func() (any, any) { return mapof.Int{"bcc": 1}, mapof.Int{} },
		"mapof.Int64":  func() (any, any) { return mapof.Int64{"bcc": 1}, mapof.Int64{} },
		"mapof.String": func() (any, any) { return mapof.String{"bcc": "a", "to": "b"}, mapof.String{"to": "b"} },
		"[]any": func() (any, any) {
			return []any{"bcc", map[string]any{"bcc": 1}, nil}, []any{"bcc", map[string]any{}, nil}
		},
		"[]map[string]any":    func() (any, any) { return []map[string]any{{"bcc": 1}}, []map[string]any{{}} },
		"[]map[string]string": func() (any, any) { return []map[string]string{{"bcc": "a"}}, []map[string]string{{}} },
		"[]mapof.Any":         func() (any, any) { return []mapof.Any{{"bcc": 1}}, []mapof.Any{{}} },
		"sliceof.Any":         func() (any, any) { return sliceof.Any{mapof.Any{"bcc": 1}}, sliceof.Any{mapof.Any{}} },
		"sliceof.MapOfAny":    func() (any, any) { return sliceof.MapOfAny{{"bcc": 1}}, sliceof.MapOfAny{{}} },
		"sliceof.MapOfString": func() (any, any) { return sliceof.MapOfString{{"bcc": "a"}}, sliceof.MapOfString{{}} },
		"database document": func() (any, any) {
			return testDocument{"bcc": 1, "list": testArray{testDocument{"bcc": 1, "id": testObjectID{1}}}},
				testDocument{"list": testArray{testDocument{"id": testObjectID{1}}}}
		},
		"database array": func() (any, any) { return testArray{testDocument{"bcc": 1}}, testArray{testDocument{}} },
		"map[string][]mapof.Any": func() (any, any) {
			return map[string][]mapof.Any{"a": {{"bcc": 1}}}, map[string][]mapof.Any{"a": {{}}}
		},
		"map[any]any": func() (any, any) {
			return map[any]any{"bcc": 1, 2: map[any]any{"bcc": 1}}, map[any]any{2: map[any]any{}}
		},
		"map[int]any": func() (any, any) {
			return map[int]any{1: map[string]any{"bcc": 1}}, map[int]any{1: map[string]any{}}
		},
		"array of maps": func() (any, any) {
			return [1]map[string]any{{"bcc": 1}}, [1]map[string]any{{}}
		},
		"*sliceof.Any": func() (any, any) {
			return &sliceof.Any{map[string]any{"bcc": 1, "to": 2}}, &sliceof.Any{map[string]any{"to": 2}}
		},
		"*sliceof.Any in a map": func() (any, any) {
			return map[string]any{"list": &sliceof.Any{mapof.Any{"bcc": 1}}},
				map[string]any{"list": &sliceof.Any{mapof.Any{}}}
		},
		"slice of slices": func() (any, any) {
			return [][]any{{map[string]any{"bcc": 1}}}, [][]any{{map[string]any{}}}
		},
	}

	for name, build := range values {
		t.Run(name, func(t *testing.T) {

			value, expected := build()

			require.NoError(t, Remove(value, "bcc"))
			require.Equal(t, expected, value)
		})
	}
}

// TestRemove_NothingToRemove proves that values with no properties pass through unchanged
func TestRemove_NothingToRemove(t *testing.T) {

	values := map[string]any{
		"nil":                nil,
		"bool":               true,
		"float32":            float32(1),
		"float64":            1.5,
		"int":                1,
		"int8":               int8(1),
		"int16":              int16(1),
		"int32":              int32(1),
		"int64":              int64(1),
		"uint":               uint(1),
		"string":             "bcc",
		"time.Time":          time.Now(),
		"[]bool":             []bool{true},
		"[]float32":          []float32{1},
		"[]float64":          []float64{1},
		"[]int":              []int{1},
		"[]int8":             []int8{1},
		"[]int16":            []int16{1},
		"[]int32":            []int32{1},
		"[]int64":            []int64{1},
		"[]string":           []string{"bcc"},
		"[]byte":             []byte("bcc"),
		"sliceof.Float":      sliceof.Float{1},
		"sliceof.Int":        sliceof.Int{1},
		"sliceof.String":     sliceof.String{"bcc"},
		"database date":      testDateTime(1),
		"database ID":        testObjectID{1},
		"database decimal":   testDecimal{h: 1},
		"database timestamp": testTimestamp{T: 1},
		"database regex":     testRegex{Pattern: "bcc"},
		"database null":      testNull{},
		"simple structs":     []testTimestamp{{T: 1}},
		"nil *sliceof.Any":   (*sliceof.Any)(nil),
		"nil map":            map[string]any(nil),
		"nil named map":      testDocument(nil),
		"nil slice":          []any(nil),
	}

	for name, value := range values {
		t.Run(name, func(t *testing.T) {

			expected, err := Copy(value)
			require.NoError(t, err)

			require.NoError(t, Remove(value, "bcc"))
			require.Equal(t, expected, value)
		})
	}
}

// TestRemove_Names proves how names are matched: whole, case-sensitive, and several at once
func TestRemove_Names(t *testing.T) {

	t.Run("several names", func(t *testing.T) {

		value := map[string]any{"bto": 1, "bcc": 2, "to": 3, "object": testDocument{"bto": 1, "cc": 2}}

		require.NoError(t, Remove(value, "bto", "bcc"))
		require.Equal(t, map[string]any{"to": 3, "object": testDocument{"cc": 2}}, value)
	})

	t.Run("no names", func(t *testing.T) {

		value := map[string]any{"bcc": 1, "object": testDocument{"bcc": 1}}

		require.NoError(t, Remove(value))
		require.Equal(t, map[string]any{"bcc": 1, "object": testDocument{"bcc": 1}}, value)
	})

	t.Run("whole names only", func(t *testing.T) {

		value := map[string]any{"bccs": 1, "BCC": 2, "": 3, "list": testDocument{"bccs": 1}}

		require.NoError(t, Remove(value, "bcc"))
		require.Len(t, value, 4)
		require.Equal(t, testDocument{"bccs": 1}, value["list"])
	})

	t.Run("a removed value is not searched", func(t *testing.T) {

		// The struct would be an error, but it is deleted along with its key
		value := map[string]any{"bcc": testStruct{}}

		require.NoError(t, Remove(value, "bcc"))
		require.Empty(t, value)
	})
}

// TestRemove_UnsupportedTypes proves that a value Remove cannot search is an error, at any
// depth, so that no property is left behind without a warning
func TestRemove_UnsupportedTypes(t *testing.T) {

	number := 1

	values := map[string]any{
		"struct":                    testStruct{},
		"pointer":                   &number,
		"pointer to slice":          &[]any{map[string]any{"bcc": 1}},
		"struct in sliceof pointer": &sliceof.Any{testStruct{}},
		"struct with an interface":  struct{ V any }{V: map[string]any{"bcc": 1}},
		"pointer to map":            &map[string]any{"bcc": 1},
		"nil pointer":               (*int)(nil),
		"function":                  func() {},
		"channel":                   make(chan int),
		"struct in map":             map[string]any{"a": testStruct{}},
		"struct in slice":           []any{testStruct{}},
		"struct in typed map":       map[string]mapof.Any{"a": {"b": testStruct{}}},
		"struct in slice of maps":   sliceof.MapOfAny{{"a": testStruct{}}},
		"struct in reflected map":   testDocument{"a": testStruct{}},
		"struct in reflected slice": testArray{testStruct{}},
		"struct in array":           [1]any{testStruct{}},
		"slice of structs":          []testStruct{{}},
		"map of structs":            map[string]testStruct{"a": {}},
	}

	for name, value := range values {
		t.Run(name, func(t *testing.T) {

			err := Remove(value, "bcc")

			require.Error(t, err)
			require.True(t, derp.IsInternalServerError(err), "an unsupported type is an Internal error")
		})
	}
}

// TestRemove_ErrorNamesTypeNotValue proves that an unsupported value is never written into
// the error
func TestRemove_ErrorNamesTypeNotValue(t *testing.T) {

	err := Remove(map[string]any{"object": testStruct{Secret: "hunter2-recipient", Hidden: map[string]any{"bcc": 1}}}, "bcc")

	require.Error(t, err)
	require.NotContains(t, errorText(t, err), "hunter2-recipient")
	require.Contains(t, errorText(t, err), "deep.testStruct")
}

// TestRemove_Remover proves that a value which removes its own properties is given every name
func TestRemove_Remover(t *testing.T) {

	t.Run("receives the names", func(t *testing.T) {

		removed := []string{}
		value := map[string]any{"a": testRemover{Removed: &removed}}

		require.NoError(t, Remove(value, "bto", "bcc"))
		require.Equal(t, []string{"bto", "bcc"}, removed)
	})

	t.Run("is told its depth, and continues from it", func(t *testing.T) {

		box := &testBox{Child: map[string]any{"bcc": 1, "to": 2}}
		value := map[string]any{"list": []any{box}}

		require.NoError(t, Remove(value, "bcc"))
		require.Equal(t, 2, box.Depth)
		require.Equal(t, map[string]any{"to": 2}, box.Child)
	})

	t.Run("a cycle through a Remover ends in an error", func(t *testing.T) {

		box := &testBox{}
		box.Child = map[string]any{"box": box}

		err := Remove(box, "bcc")

		require.Error(t, err)
		require.Less(t, len(errorText(t, err)), 2_000)
	})

	t.Run("its error is returned", func(t *testing.T) {
		require.Error(t, Remove([]any{testFailingCopier{}}, "bcc"))
	})
}

// TestRemove_MaxDepth proves that a walk stops at the maximum depth, which is what stops a
// value that contains itself
func TestRemove_MaxDepth(t *testing.T) {

	t.Run("map that contains itself", func(t *testing.T) {

		value := map[string]any{}
		value["self"] = value

		err := Remove(value, "bcc")

		require.Error(t, err)
		require.True(t, derp.IsInternalServerError(err))
		require.Less(t, len(errorText(t, err)), 2_000)
	})

	t.Run("reflected containers that contain themselves", func(t *testing.T) {

		document := testDocument{}
		document["list"] = testArray{document}

		require.Error(t, Remove(document, "bcc"))
	})

	t.Run("pointer that leads to itself", func(t *testing.T) {

		value := &sliceof.Any{nil}
		(*value)[0] = value

		require.Error(t, Remove(value, "bcc"))
	})

	t.Run("the limit allows what a JSON decoder allows", func(t *testing.T) {

		var value any = map[string]any{"bcc": 1}

		for range maxDepth {
			value = []any{value}
		}

		require.NoError(t, Remove(value, "bcc"))
		require.False(t, hasKey(value, "bcc"))

		require.Error(t, Remove([]any{value}, "bcc"))
	})
}
