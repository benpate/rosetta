package schema

import (
	"strings"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

// TestValidateObject_Wildcard_Normalize confirms that Normalize sanitizes a wildcard value
func TestValidateObject_Wildcard_Normalize(t *testing.T) {

	s := New(Object{Wildcard: String{Format: "html"}})
	value := mapof.Any{"bio": "<script>alert(1)</script><b>hello</b>"}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"bio"}, rewrites)
	require.NotContains(t, value["bio"], "<script>")
	require.Contains(t, value["bio"], "hello")
}

// TestValidateObject_Wildcard_Validate confirms that Validate rejects a bad wildcard value
func TestValidateObject_Wildcard_Validate(t *testing.T) {

	s := New(Object{Wildcard: String{MaxLength: 5}})

	_, err := s.Validate(mapof.Any{"short": "abc"})
	require.Nil(t, err)

	_, err = s.Validate(mapof.Any{"long": "abcdefgh"})
	require.NotNil(t, err)
}

// TestValidateObject_Wildcard_DeclaredPropertyWins confirms that named properties skip the wildcard
func TestValidateObject_Wildcard_DeclaredPropertyWins(t *testing.T) {

	s := New(Object{
		Properties: ElementMap{"name": String{MaxLength: 100}},
		Wildcard:   String{MaxLength: 3},
	})

	value := mapof.Any{"name": "a longer name", "tag": "abc"}

	_, err := s.Validate(value)
	require.Nil(t, err)
}

// TestValidateObject_Wildcard_Nested confirms that nested wildcards reach the innermost value
func TestValidateObject_Wildcard_Nested(t *testing.T) {

	// The shape of a translation map: language, then field
	s := New(Object{Wildcard: Object{Wildcard: String{Format: "html"}}})
	value := mapof.Any{"pt-BR": mapof.Any{"summary": "<script>x</script>ok"}}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"pt-BR.summary"}, rewrites)

	inner, ok := value["pt-BR"].(mapof.Any)
	require.True(t, ok)
	require.False(t, strings.Contains(inner["summary"].(string), "<script>"))
}

// TestValidateObject_Wildcard_IgnoresNonMaps confirms that the wildcard skips non-listable maps
func TestValidateObject_Wildcard_IgnoresNonMaps(t *testing.T) {

	s := New(Object{Wildcard: String{MaxLength: 3}})

	// declaredMap declares IsMap but has no Keys, so its keys cannot be listed
	_, err := s.Validate(declaredMap{"anything": "longer than three"})
	require.Nil(t, err)
}

// translationEntry is the shape of one language's translation: values a form may write, beside
// status that no schema lists.  It is held by value in a typed map, as Emissary will store it.
type translationEntry struct {
	Values mapof.Any
	Origin string
}

// GetPointer exposes only the values, so the status can never be reached through a schema
func (entry *translationEntry) GetPointer(name string) (any, bool) {
	if name == "values" {
		return &entry.Values, true
	}
	return nil, false
}

// translationSchema describes a translation map: language, then the translatable values,
// with content carrying its raw source and its rendered HTML
func translationSchema() Schema {
	return New(Object{Wildcard: Object{Properties: ElementMap{
		"values": Object{Properties: ElementMap{
			"summary": String{Format: "html"},
			"content": Object{Properties: ElementMap{
				"raw":  String{Format: "unsafe-any"},
				"html": String{Format: "html"},
			}},
		}},
	}}})
}

// TestValidateObject_Wildcard_StructEntry confirms that Normalize sanitizes values nested inside
// a struct held by value in a typed map, and leaves the struct's other fields alone
func TestValidateObject_Wildcard_StructEntry(t *testing.T) {

	value := mapof.Object[translationEntry]{
		"pt-BR": {Values: mapof.Any{"summary": "<script>x</script>resumo"}, Origin: "HUMAN"},
	}

	rewrites, err := translationSchema().Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"pt-BR.values.summary"}, rewrites)

	entry := value["pt-BR"]
	require.Equal(t, "HUMAN", entry.Origin, "fields outside the schema are kept")
	require.NotContains(t, entry.Values["summary"], "<script>")
	require.Contains(t, entry.Values["summary"], "resumo")
}

// TestValidateObject_Wildcard_StructEntry_Content confirms that Normalize sanitizes rendered
// content HTML inside a translation, and keeps its raw source verbatim
func TestValidateObject_Wildcard_StructEntry_Content(t *testing.T) {

	value := mapof.Object[translationEntry]{
		"pt-BR": {Values: mapof.Any{
			"content": mapof.Any{"raw": "# <script>raw</script>", "html": "<script>y</script><p>corpo</p>"},
		}},
	}

	rewrites, err := translationSchema().Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"pt-BR.values.content.html"}, rewrites)

	content, ok := value["pt-BR"].Values["content"].(mapof.Any)
	require.True(t, ok)
	require.Equal(t, "# <script>raw</script>", content["raw"], "unsafe-any keeps the raw source verbatim")
	require.NotContains(t, content["html"], "<script>")
	require.Contains(t, content["html"], "corpo")
}

// TestValidateObject_NamedNestedMap_Normalize confirms that Normalize can write a rewrite back
// into a map held by value beneath a named property
func TestValidateObject_NamedNestedMap_Normalize(t *testing.T) {

	s := New(Object{Properties: ElementMap{"content": Object{Properties: ElementMap{"html": String{Format: "html"}}}}})
	value := mapof.Any{"content": mapof.Any{"html": "<script>y</script><p>ok</p>"}}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Equal(t, []string{"content.html"}, rewrites)

	inner, ok := value["content"].(mapof.Any)
	require.True(t, ok)
	require.Equal(t, "<p>ok</p>", inner["html"])
}

// TestValidateObject_Wildcard_StructEntry_Validate confirms that Validate rejects unsafe values
// nested inside a struct held by value in a typed map
func TestValidateObject_Wildcard_StructEntry_Validate(t *testing.T) {

	value := mapof.Object[translationEntry]{
		"pt-BR": {Values: mapof.Any{"summary": "<script>x</script>resumo"}},
	}

	_, err := translationSchema().Validate(&value)
	require.NotNil(t, err)
}

// TestValidateObject_Wildcard_UnlistedKeysAreSkipped confirms that a key the schema does not
// list is neither checked nor removed, so a caller that stores one must sanitize it itself
func TestValidateObject_Wildcard_UnlistedKeysAreSkipped(t *testing.T) {

	s := New(Object{Wildcard: Object{Properties: ElementMap{"raw": String{Format: "unsafe-any"}}}})
	value := mapof.Any{"pt-BR": mapof.Any{"raw": "texto", "html": "<script>z</script>"}}

	rewrites, err := s.Normalize(&value)
	require.Nil(t, err)
	require.Empty(t, rewrites)

	inner, ok := value["pt-BR"].(mapof.Any)
	require.True(t, ok)
	require.Equal(t, "<script>z</script>", inner["html"])
}
