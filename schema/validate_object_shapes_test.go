package schema

import (
	"encoding/json"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

// shapeSchemas returns a translation map described two ways: by wildcard, and by named property
func shapeSchemas() map[string]Schema {
	html := String{Format: "html"}

	return map[string]Schema{
		"wildcard": New(Object{Wildcard: Object{Wildcard: html}}),
		"named":    New(Object{Wildcard: Object{Properties: ElementMap{"summary": html}}}),
	}
}

// decodeJSON decodes a document the way a request body is decoded, with every nested object
// held as a plain map[string]any
func decodeJSON(t *testing.T, document string) mapof.Any {
	result := mapof.Any{}
	require.NoError(t, json.Unmarshal([]byte(document), &result))
	return result
}

// TestSet_MapHeldByValue confirms that Set sanitizes and stores a map passed by value
func TestSet_MapHeldByValue(t *testing.T) {

	for name, s := range shapeSchemas() {
		t.Run(name, func(t *testing.T) {

			object := mapof.Any{}
			err := s.Set(&object, "pt-BR", mapof.Any{"summary": "<script>x</script>ok"})
			require.NoError(t, err)

			stored, ok := object["pt-BR"].(mapof.Any)
			require.True(t, ok, "the map is stored by value, as it was passed")
			require.Equal(t, "ok", stored["summary"])
		})
	}
}

// TestSet_PlainMapHeldByValue confirms that Set sanitizes a plain Go map, and stores it as one
func TestSet_PlainMapHeldByValue(t *testing.T) {

	for name, s := range shapeSchemas() {
		t.Run(name, func(t *testing.T) {

			object := mapof.Any{}
			err := s.Set(&object, "pt-BR", map[string]any{"summary": "<script>x</script>ok"})
			require.NoError(t, err)

			stored, ok := object["pt-BR"].(map[string]any)
			require.True(t, ok, "a plain map is stored as its own type")
			require.Equal(t, "ok", stored["summary"])
		})
	}
}

// TestNormalize_MapHeldByValue confirms that Normalize rewrites a map passed by value exactly as
// it rewrites one passed by pointer
func TestNormalize_MapHeldByValue(t *testing.T) {

	for name, s := range shapeSchemas() {
		t.Run(name, func(t *testing.T) {

			value := mapof.Any{"pt-BR": mapof.Any{"summary": "<script>x</script>ok"}}
			rewrites, err := s.Normalize(value)
			require.NoError(t, err)
			require.Equal(t, []string{"pt-BR.summary"}, rewrites)
			require.Equal(t, mapof.Any{"pt-BR": mapof.Any{"summary": "ok"}}, value)
		})
	}
}

// TestValidate_MapHeldByValue confirms that Validate reports a rewrite in a map passed by value
// as a validation error, as it does for one passed by pointer, and leaves the rewrite applied
func TestValidate_MapHeldByValue(t *testing.T) {

	for name, s := range shapeSchemas() {
		t.Run(name, func(t *testing.T) {

			byValue := mapof.Any{"pt-BR": mapof.Any{"summary": "<script>x</script>ok"}}
			_, valueErr := s.Validate(byValue)

			byPointer := mapof.Any{"pt-BR": mapof.Any{"summary": "<script>x</script>ok"}}
			_, pointerErr := s.Validate(&byPointer)

			require.True(t, derp.IsValidationError(valueErr))
			require.True(t, derp.IsValidationError(pointerErr))
			require.Equal(t, byPointer, byValue)

			_, err := s.Validate(mapof.Any{"pt-BR": mapof.Any{"summary": "clean"}})
			require.NoError(t, err)
		})
	}
}

// TestNormalize_PlainNestedMap confirms that a document decoded from JSON, whose nested objects
// are plain maps, is read and sanitized, and keeps its plain maps
func TestNormalize_PlainNestedMap(t *testing.T) {

	for name, s := range shapeSchemas() {
		t.Run(name, func(t *testing.T) {

			clean := decodeJSON(t, `{"pt-BR":{"summary":"clean"}}`)
			rewrites, err := s.Normalize(&clean)
			require.NoError(t, err)
			require.Empty(t, rewrites)

			_, err = s.Validate(&clean)
			require.NoError(t, err)

			dirty := decodeJSON(t, `{"pt-BR":{"summary":"<script>x</script>ok"}}`)
			rewrites, err = s.Normalize(&dirty)
			require.NoError(t, err)
			require.Equal(t, []string{"pt-BR.summary"}, rewrites)
			require.Equal(t, map[string]any{"summary": "ok"}, dirty["pt-BR"])
		})
	}
}

// TestValidateObject_WrongKind confirms that a single value or a list where an Object is expected
// is refused as a validation error that does not carry the value, and that nil is allowed
func TestValidateObject_WrongKind(t *testing.T) {

	wrongKinds := map[string]any{
		"string": "<script>" + testSecret + "</script>",
		"number": 42,
		"bool":   true,
		"list":   []any{"<script>" + testSecret + "</script>"},
	}

	for name, s := range shapeSchemas() {
		for kind, wrong := range wrongKinds {
			t.Run(name+"/"+kind, func(t *testing.T) {

				value := mapof.Any{"pt-BR": wrong}
				_, err := s.Normalize(&value)
				require.True(t, derp.IsValidationError(err))
				requireNoSecretInDetails(t, err)

				err = s.Set(&mapof.Any{}, "pt-BR", wrong)
				require.True(t, derp.IsValidationError(err))
				requireNoSecretInDetails(t, err)
			})
		}
	}

	// A wildcard Object allows a nil value, which holds nothing to sanitize
	value := mapof.Any{"pt-BR": nil}
	_, err := shapeSchemas()["wildcard"].Normalize(&value)
	require.NoError(t, err)
}
