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
