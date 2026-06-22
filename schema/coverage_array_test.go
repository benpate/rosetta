package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArray_Validate_Required(t *testing.T) {
	// A required array must contain at least one item
	schema := New(Array{Items: String{}, Required: true})

	_, _, err := Validate(schema, &testArrayA{})
	require.Error(t, err)
}

func TestArray_Validate_MaxLength(t *testing.T) {
	// An array longer than MaxLength is rejected
	schema := New(Array{Items: String{}, MaxLength: 1})

	_, _, err := Validate(schema, &testArrayA{"one", "two"})
	require.Error(t, err)
}

func TestArray_DefaultValue(t *testing.T) {
	require.Equal(t, []any{}, Array{}.DefaultValue())
}

func TestArray_IsRequired(t *testing.T) {
	require.False(t, Array{}.IsRequired())
	require.True(t, Array{Required: true}.IsRequired())
}

func TestArray_GetElement_Self(t *testing.T) {
	element := Array{Items: String{}}
	found, ok := element.GetElement("")
	require.True(t, ok)
	require.Equal(t, element, found)
}

func TestArray_GetElement_Index(t *testing.T) {
	found, ok := Array{Items: String{}}.GetElement("0")
	require.True(t, ok)
	require.Equal(t, String{}, found)
}

func TestArray_GetElement_BoundedIndex(t *testing.T) {
	found, ok := Array{Items: String{}, MaxLength: 5}.GetElement("2")
	require.True(t, ok)
	require.Equal(t, String{}, found)
}

func TestArray_GetElement_NotAnIndex(t *testing.T) {
	_, ok := Array{Items: String{}}.GetElement("nope")
	require.False(t, ok)
}

func TestArray_Inherit(t *testing.T) {
	require.NotPanics(t, func() { Array{}.Inherit(String{}) })
}

func TestArray_AllProperties(t *testing.T) {
	element := Array{Items: String{}}
	require.Equal(t, ElementMap{"": element}, element.AllProperties())
}

func TestArray_GetLength(t *testing.T) {
	value := testArrayA{"one", "two"}

	length, ok := Array{}.GetLength(value)
	require.True(t, ok)
	require.Equal(t, 2, length)
}

func TestArray_GetLength_NotLengthGetter(t *testing.T) {
	_, ok := Array{}.GetLength("not-an-array")
	require.False(t, ok)
}

func TestArray_GetIndex(t *testing.T) {
	value := testArrayA{"one", "two"}

	item, ok := Array{}.GetIndex(value, 1)
	require.True(t, ok)
	require.Equal(t, "two", item)
}

func TestArray_GetIndex_NotArrayGetter(t *testing.T) {
	_, ok := Array{}.GetIndex("not-an-array", 0)
	require.False(t, ok)
}

func TestArray_SetIndex(t *testing.T) {
	value := testArrayA{"one", "two"}

	ok := Array{}.SetIndex(&value, 1, "TWO")
	require.True(t, ok)
	require.Equal(t, "TWO", value[1])
}

func TestArray_SetIndex_NotArraySetter(t *testing.T) {
	require.False(t, Array{}.SetIndex("not-an-array", 0, "x"))
}

// growableArray is a string slice whose SetIndex grows the slice as needed,
// which is required to exercise the Append behavior.
type growableArray []string

func (a growableArray) Length() int {
	return len(a)
}

func (a *growableArray) SetIndex(index int, value any) bool {
	stringValue, ok := value.(string)
	if !ok {
		return false
	}

	for index >= len(*a) {
		*a = append(*a, "")
	}
	(*a)[index] = stringValue
	return true
}

func TestArray_Append(t *testing.T) {
	value := growableArray{"one"}

	err := Array{}.Append(&value, "two")
	require.NoError(t, err)
	require.Equal(t, growableArray{"one", "two"}, value)
}

func TestArray_Append_Fails(t *testing.T) {
	// SetIndex rejects non-string values, so Append returns an error
	value := growableArray{"one"}
	require.Error(t, Array{}.Append(&value, 42))
}

func TestArray_MarshalMap(t *testing.T) {
	result := Array{Items: String{}, MinLength: 1, MaxLength: 5, Required: true}.MarshalMap()

	require.Equal(t, TypeArray, result["type"])
	require.Equal(t, 1, result["minLength"])
	require.Equal(t, 5, result["maxLength"])
	require.Equal(t, true, result["required"])
	require.Contains(t, result, "items")
}

func TestArray_UnmarshalMap(t *testing.T) {
	element := Array{}
	err := element.UnmarshalMap(map[string]any{
		"type":      "array",
		"items":     map[string]any{"type": "string"},
		"minLength": 1,
		"maxLength": 5,
		"required":  true,
	})

	require.NoError(t, err)
	require.IsType(t, String{}, element.Items)
	require.Equal(t, 1, element.MinLength)
	require.Equal(t, 5, element.MaxLength)
	require.True(t, element.Required)
}

func TestArray_UnmarshalMap_WrongType(t *testing.T) {
	require.Error(t, (&Array{}).UnmarshalMap(map[string]any{"type": "string"}))
}

func TestArray_UnmarshalMap_BadItems(t *testing.T) {
	// "items" with an unrecognized type cannot be unmarshalled
	err := (&Array{}).UnmarshalMap(map[string]any{
		"type":  "array",
		"items": map[string]any{"type": "not-a-real-type"},
	})
	require.Error(t, err)
}
