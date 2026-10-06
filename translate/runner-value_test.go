package translate

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

// TestValue_NeverSharesItsValue requires each execution to write a copy of a value rule's
// object, because the rule is parsed once and reused by every caller
func TestValue_NeverSharesItsValue(t *testing.T) {

	pipeline, err := NewFromMap(
		map[string]any{"target": "url", "value": map[string]any{"type": "Link", "tag": map[string]any{"a": "b"}, "list": []any{"x"}}},
		map[string]any{"target": "url.href", "path": "href"},
		map[string]any{"target": "url.tag.c", "path": "href"},
		map[string]any{"target": "url.list.1", "path": "href"},
	)
	require.NoError(t, err)

	inSchema := schema.New(schema.Object{Properties: schema.ElementMap{"href": schema.String{}}})
	outSchema := schema.New(schema.Object{Wildcard: schema.Any{}})

	first := mapof.Any{}
	require.NoError(t, pipeline.Execute(inSchema, mapof.Any{"href": "https://one.example"}, outSchema, &first))

	second := mapof.Any{}
	require.NoError(t, pipeline.Execute(inSchema, mapof.Any{"href": "https://two.example"}, outSchema, &second))

	// Writing the second document changes nothing in the first, at any depth
	expected := mapof.Any{"type": "Link", "href": "https://one.example", "tag": mapof.Any{"a": "b", "c": "https://one.example"}, "list": []any{"x", "https://one.example"}}
	require.Equal(t, expected, normalize(first["url"]))

	// And the rule still holds its original value
	require.Equal(t, mapof.Any{"type": "Link", "tag": mapof.Any{"a": "b"}, "list": []any{"x"}}, normalize(pipeline[0].Runner.(valueRunner).Value))
}

// TestAppend_NeverSharesItsValue requires each execution to append a copy of an append rule's object
func TestAppend_NeverSharesItsValue(t *testing.T) {

	pipeline, err := NewFromMap(
		map[string]any{"target": "list", "value": []any{}},
		map[string]any{"target": "list", "append": map[string]any{"type": "Item"}},
		map[string]any{"target": "list.0.id", "path": "id"},
	)
	require.NoError(t, err)

	inSchema := schema.New(schema.Object{Properties: schema.ElementMap{"id": schema.String{}}})
	outSchema := schema.New(schema.Object{Wildcard: schema.Any{}})

	first := mapof.Any{}
	require.NoError(t, pipeline.Execute(inSchema, mapof.Any{"id": "one"}, outSchema, &first))

	second := mapof.Any{}
	require.NoError(t, pipeline.Execute(inSchema, mapof.Any{"id": "two"}, outSchema, &second))

	require.Equal(t, []any{mapof.Any{"type": "Item", "id": "one"}}, normalize(first["list"]))
	require.Equal(t, mapof.Any{"type": "Item"}, normalize(pipeline[1].Runner.(appendRunner).Append))
}

// normalize converts every rosetta map and list beneath a value into mapof.Any and []any, so
// that values can be compared without regard to which container type holds them
func normalize(value any) any {

	switch typed := value.(type) {

	case map[string]any:
		return normalize(mapof.Any(typed))

	case mapof.Any:
		result := mapof.Any{}
		for key, item := range typed {
			result[key] = normalize(item)
		}
		return result

	case *mapof.Any:
		if typed == nil {
			return nil
		}
		return normalize(*typed)

	case sliceof.Any:
		return normalize([]any(typed))

	case *sliceof.Any:
		if typed == nil {
			return nil
		}
		return normalize(*typed)

	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			result[index] = normalize(item)
		}
		return result
	}

	return value
}

// TestCopyValue requires a copy that shares no map or list with its original, at any depth
func TestCopyValue(t *testing.T) {

	original := mapof.Any{
		"plain":   map[string]any{"a": "b"},
		"list":    []any{map[string]any{"c": "d"}},
		"rosetta": sliceof.Any{mapof.Any{"e": "f"}},
		"number":  1,
	}

	result := copyValue(original).(mapof.Any)
	require.Equal(t, original, result)

	result["plain"].(map[string]any)["a"] = "changed"
	result["list"].([]any)[0].(map[string]any)["c"] = "changed"
	result["rosetta"].(sliceof.Any)[0].(mapof.Any)["e"] = "changed"
	result["number"] = 2

	require.Equal(t, mapof.Any{
		"plain":   map[string]any{"a": "b"},
		"list":    []any{map[string]any{"c": "d"}},
		"rosetta": sliceof.Any{mapof.Any{"e": "f"}},
		"number":  1,
	}, original)

	require.Equal(t, "text", copyValue("text"))
	require.Nil(t, copyValue(nil))
}
