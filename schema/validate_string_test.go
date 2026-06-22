package schema

import (
	"testing"
	"unicode/utf8"

	"github.com/benpate/rosetta/schema/format"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

// testStringGetter tests getting/setting string values
type testStringGetter struct {
	value string
}

// GetStringOK gets a string value
func (t testStringGetter) GetStringOK(_ string) (string, bool) {
	return t.value, true
}

// SetString sets a string value
func (t *testStringGetter) SetString(_ string, value string) bool {
	t.value = value
	return true
}

// TestStringGetter tests getting/setting string values
func TestStringGetter_Empty(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Required: true},
		},
	})

	var getter testStringGetter

	value, err := schema.Get(&getter, "value")
	require.NoError(t, err)
	require.Equal(t, "", value)

	_, changed, err := Validate(schema, &getter) // Invalid because field is required
	require.False(t, changed)
	require.Error(t, err)
}

// TestStringSetter tests getting/setting string values
func TestStringSetter(t *testing.T) {

	spew.Config.DisableMethods = true

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Required: true},
		},
	})

	var getter testStringGetter

	err := schema.Set(&getter, "value", "hello")
	require.NoError(t, err)
	require.Equal(t, "hello", getter.value)

	value, err := schema.Get(&getter, "value")
	require.Equal(t, "hello", value)
	require.NoError(t, err)
}

// TestStringSetter tests getting/setting string values
func TestStringSetter_Required(t *testing.T) {

	spew.Config.DisableMethods = true

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Required: true},
		},
	})

	var getter testStringGetter

	err := schema.Set(&getter, "value", "hello")
	require.NoError(t, err)
	require.Equal(t, "hello", getter.value)

	value, err := schema.Get(&getter, "value")
	require.Equal(t, "hello", value)
	require.NoError(t, err)

	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)
}

// TestStringGetter tests getting/setting string values
func TestStringGetter_ValidateOnly(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Required: true},
		},
	})

	var getter testStringGetter

	err := schema.Set(&getter, "value", "hello")
	require.NoError(t, err)
	require.Equal(t, "hello", getter.value)

	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)
	// require.NoError(t, err)
}

func TestStringValidator(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{MinLength: 10, MaxLength: 20},
		},
	})

	var getter testStringGetter

	// Test Min-Length
	require.Error(t, schema.Set(&getter, "value", "too-short"))

	_, _, err := Validate(schema, &getter)
	require.Error(t, err)

	// Test Max-Length
	require.NoError(t, schema.Set(&getter, "value", "this-one-is-way-too-long-so-it-shouldn't-get-set-like-this"))
	require.Equal(t, "this-one-is-way-too-", getter.value)
	value, err := schema.Get(getter, "value")
	require.NoError(t, err)
	require.Equal(t, "this-one-is-way-too-", value)

	// Test OK
	require.NoError(t, schema.Set(&getter, "value", "this-is-just-right"))
	require.Equal(t, "this-is-just-right", getter.value)

	_, _, err = Validate(schema, &getter)
	require.NoError(t, err)

}

// TestStringValidator_Multibyte confirms that length limits are measured in
// runes (not bytes) and that truncation lands on a rune boundary, never
// splitting a multi-byte character into invalid UTF-8.
func TestStringValidator_Multibyte(t *testing.T) {

	// Register a passthrough format so this test exercises the rune-aware
	// length/truncation logic in isolation, rather than the default `no-html`
	// sanitizer. Reset the freeze latch first, since earlier tests may have
	// frozen the registry by reading from it.
	formatsFrozen.Store(false)
	UseFormat("test-passthrough", func(_ string) format.StringFormat {
		return func(value string) (string, error) { return value, nil }
	})

	schema := New(Object{
		Properties: ElementMap{
			"value": String{MinLength: 3, MaxLength: 5, Format: "test-passthrough"},
		},
	})

	var getter testStringGetter

	// Five 3-byte runes (15 bytes) is exactly MaxLength=5 runes, so it must pass unchanged.
	require.NoError(t, schema.Set(&getter, "value", "日本語です"))
	require.Equal(t, "日本語です", getter.value)

	// Seven runes must truncate to the first five runes (not the first five bytes),
	// and the result must remain valid UTF-8.
	require.NoError(t, schema.Set(&getter, "value", "日本語ですよね"))
	require.Equal(t, "日本語です", getter.value)
	require.True(t, utf8.ValidString(getter.value))

	// Two runes is below MinLength=3 (even though it is 6 bytes), so it must fail.
	require.Error(t, schema.Set(&getter, "value", "日本"))
}

// testStringPointer tests getting/setting bool values via a pointer
type testStringPointer struct {
	value string
}

func (test *testStringPointer) GetPointer(_ string) (any, bool) {
	return &test.value, true
}

func TestStringPointer(t *testing.T) {

	schema := New(Object{
		Properties: ElementMap{
			"value": String{Enum: []string{"hello", "there", "general", "kenobi"}},
		},
	})

	var getter testStringPointer

	// Validate correct value
	require.NoError(t, schema.Set(&getter, "value", "kenobi"))
	require.Equal(t, "kenobi", getter.value)

	// Validate incorrect value
	require.Error(t, schema.Set(&getter, "value", "invalid"))

	// Re-get the original value that was set
	value, err := schema.Get(&getter, "value")
	require.NoError(t, err)
	require.Equal(t, "kenobi", value)
}
