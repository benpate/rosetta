package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRegisteredFormat(t *testing.T) {

	for _, registered := range []string{"url", "uri", "no-html", "token", "in", "webfinger"} {
		require.True(t, IsRegisteredFormat(registered), registered)
	}

	for _, unregistered := range []string{"", "bogus", "URL", "no html"} {
		require.False(t, IsRegisteredFormat(unregistered), unregistered)
	}
}

func TestValidateFormats(t *testing.T) {

	// requireValid asserts that every format name in the schema resolves.
	requireValid := func(t *testing.T, element Element) {
		t.Helper()
		require.Nil(t, New(element).ValidateFormats())
	}

	// requireInvalid asserts that at least one format name fails to resolve.
	requireInvalid := func(t *testing.T, element Element) {
		t.Helper()
		require.NotNil(t, New(element).ValidateFormats())
	}

	t.Run("EmptySchema", func(t *testing.T) {
		require.Nil(t, Schema{}.ValidateFormats())
	})

	t.Run("StringFormats", func(t *testing.T) {
		requireValid(t, String{})                             // empty format selects the no-html default
		requireValid(t, String{Format: " "})                  // whitespace-only behaves like empty
		requireValid(t, String{Format: "url"})                // single format
		requireValid(t, String{Format: "token in=red,green"}) // multiple formats with a parameter
		requireValid(t, String{Format: "url "})               // trailing space is skipped by the runtime
		requireInvalid(t, String{Format: "bogus"})            // unknown name
		requireInvalid(t, String{Format: "url bogus"})        // unknown name after a valid one
		requireInvalid(t, String{Format: "URL"})              // format names are case-sensitive
		requireInvalid(t, String{Format: "bogus=with-value"}) // unknown name with a parameter
	})

	t.Run("NestedElements", func(t *testing.T) {

		valid := Object{
			Properties: ElementMap{
				"link":  String{Format: "url"},
				"count": Integer{},
				"tags":  Array{Items: String{Format: "token"}},
				"child": Object{
					Properties: ElementMap{
						"email": String{Format: "email"},
					},
				},
			},
		}

		requireValid(t, valid)

		requireInvalid(t, Object{Properties: ElementMap{"link": String{Format: "bogus"}}})
		requireInvalid(t, Array{Items: String{Format: "bogus"}})
		requireInvalid(t, Array{Items: Object{Properties: ElementMap{"deep": String{Format: "bogus"}}}})
		requireInvalid(t, Object{Wildcard: String{Format: "bogus"}})
	})

	t.Run("LeafElementsWithoutFormats", func(t *testing.T) {
		requireValid(t, Any{})
		requireValid(t, Boolean{})
		requireValid(t, Integer{})
		requireValid(t, Number{})
		requireValid(t, Array{}) // nil Items is tolerated
	})

	t.Run("PointerElements", func(t *testing.T) {

		requireValid(t, &String{Format: "url"})
		requireInvalid(t, &String{Format: "bogus"})
		requireInvalid(t, &Array{Items: String{Format: "bogus"}})
		requireInvalid(t, &Object{Properties: ElementMap{"link": String{Format: "bogus"}}})

		var nilString *String
		var nilArray *Array
		var nilObject *Object
		requireValid(t, nilString)
		requireValid(t, nilArray)
		requireValid(t, nilObject)
	})
}

func TestFormatNames(t *testing.T) {

	require.Empty(t, formatNames(""))
	require.Empty(t, formatNames("   "))
	require.Equal(t, []string{"url"}, formatNames("url"))
	require.Equal(t, []string{"url", "token"}, formatNames("url token"))
	require.Equal(t, []string{"in", "token"}, formatNames("in=red,green token"))
	require.Equal(t, []string{"url"}, formatNames(" url "))
}
