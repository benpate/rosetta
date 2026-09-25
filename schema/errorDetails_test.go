package schema

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

// testSecret stands in for a password, key, or token held by an object being validated or set
const testSecret = "s3cr3t-value"

// TestErrorDetails_NameTypesNotValues requires that no error from validating or setting an object
// reports what that object holds
func TestErrorDetails_NameTypesNotValues(t *testing.T) {

	// These errors once attached the whole object, or the whole value being set.  Emissary
	// stores every reported error, so a failed save of an OAuth client stored its secret.
	schema := New(Object{Properties: ElementMap{
		"name":   String{Required: true},
		"secret": String{},
		"child":  Object{Properties: ElementMap{"inner": String{}}},
	}})

	// A plain map implements no getter, so every property read fails with the object in hand
	newObject := func() map[string]any {
		return map[string]any{"name": "", "secret": testSecret}
	}

	t.Run("Validate", func(t *testing.T) {
		object := newObject()
		_, err := schema.Validate(&object)
		require.Equal(t, "schema.Schema.Validate", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})

	t.Run("Normalize", func(t *testing.T) {
		object := newObject()
		_, err := schema.Normalize(&object)
		require.Equal(t, "schema.Schema.Normalize", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		object := newObject()
		err := schema.Set(&object, "child", testSecret)
		require.Equal(t, "schema.Schema.Set", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})

	t.Run("SetAll", func(t *testing.T) {
		object := newObject()
		err := schema.SetAll(&object, map[string]any{"child": testSecret})
		require.Equal(t, "schema.Schema.SetAll", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})

	t.Run("SetURLValues", func(t *testing.T) {
		object := newObject()
		err := schema.SetURLValues(&object, url.Values{"child": {testSecret}})
		require.Equal(t, "schema.Schema.SetURLValues", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})

	t.Run("SetByReflection", func(t *testing.T) {
		var target string
		err := setByReflection(target, testSecret)
		require.Equal(t, "schema.setByReflection", derp.Location(err))
		requireNoSecretInDetails(t, err)
	})
}

// requireNoSecretInDetails fails if any layer of the error chain carries the secret in its message
// or details.  "%#v" prints every field, including those hidden from JSON.
func requireNoSecretInDetails(t *testing.T, err error) {

	t.Helper()
	require.Error(t, err)

	for ; err != nil; err = errors.Unwrap(err) {
		require.NotContains(t, derp.Message(err), testSecret)

		for _, detail := range derp.Details(err) {
			rendered := fmt.Sprintf("%#v", detail)
			carries := strings.Contains(rendered, testSecret)
			require.False(t, carries, "%s carries the secret in a detail", derp.Location(err))
		}
	}
}
