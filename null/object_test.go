package null

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// testObject is a sample struct used to exercise Object[T] with a non-primitive type
type testObject struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestObject_ZeroValue(t *testing.T) {

	var value Object[testObject]

	require.True(t, value.IsNull())
	require.False(t, value.IsPresent())
	require.Equal(t, testObject{}, value.Object())
	require.Nil(t, value.Interface())
}

func TestObject_NewObject(t *testing.T) {

	value := NewObject(testObject{Name: "Han", Count: 12})

	require.False(t, value.IsNull())
	require.True(t, value.IsPresent())
	require.Equal(t, testObject{Name: "Han", Count: 12}, value.Object())
	require.Equal(t, testObject{Name: "Han", Count: 12}, value.Interface())
}

func TestObject_NewObject_ZeroValueIsPresent(t *testing.T) {

	// A present-but-zero value must NOT read as null
	value := NewObject(testObject{})

	require.False(t, value.IsNull())
	require.True(t, value.IsPresent())
}

func TestObject_SetAndUnset(t *testing.T) {

	var value Object[testObject]

	value.Set(testObject{Name: "Leia", Count: 1})
	require.True(t, value.IsPresent())
	require.Equal(t, "Leia", value.Object().Name)

	// Unset must reset BOTH the value and the present flag
	value.Unset()
	require.True(t, value.IsNull())
	require.Equal(t, testObject{}, value.Object())
	require.Nil(t, value.Interface())
}

func TestObject_Nullable(t *testing.T) {

	// Object[T] must satisfy the Nullable interface
	var value Nullable = Object[string]{}
	require.True(t, value.IsNull())

	value = NewObject("present")
	require.False(t, value.IsNull())
}

func TestObject_Primitives(t *testing.T) {

	// Object[T] works with primitives, too
	stringValue := NewObject("hello")
	require.Equal(t, "hello", stringValue.Object())

	intValue := NewObject(42)
	require.Equal(t, 42, intValue.Object())

	boolValue := NewObject(false)
	require.False(t, boolValue.Object())
	require.True(t, boolValue.IsPresent())
}

func TestObject_Slices(t *testing.T) {

	value := NewObject([]string{"a", "b"})
	require.Equal(t, []string{"a", "b"}, value.Object())

	// An empty (but present) slice is not null
	empty := NewObject([]string{})
	require.True(t, empty.IsPresent())
	require.Empty(t, empty.Object())

	// Unset returns the slice to its nil zero value
	empty.Unset()
	require.Nil(t, empty.Object())
}

func TestObject_Pointers(t *testing.T) {

	// A present-but-nil pointer is still "present"
	var nilPointer *testObject
	value := NewObject(nilPointer)

	require.True(t, value.IsPresent())
	require.Nil(t, value.Object())

	// ...but Interface() returns a TYPED nil, which is not equal to an untyped nil
	require.False(t, value.Interface() == nil)
}

func TestObject_IsNil(t *testing.T) {

	// IsNil must track IsNull exactly
	var value Object[testObject]
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Set(testObject{Name: "Lando"})
	require.False(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())

	value.Unset()
	require.True(t, value.IsNil())
	require.Equal(t, value.IsNull(), value.IsNil())
}

func TestObject_IsZero(t *testing.T) {

	// A null value is always zero
	var value Object[testObject]
	require.True(t, value.IsZero())

	// A present-but-zero value is ALSO zero
	value.Set(testObject{})
	require.True(t, value.IsZero())
	require.True(t, value.IsPresent())

	// A present, non-zero value is not
	value.Set(testObject{Name: "Lando"})
	require.False(t, value.IsZero())

	value.Unset()
	require.True(t, value.IsZero())
}

func TestObject_IsZero_Primitives(t *testing.T) {

	require.True(t, NewObject("").IsZero())
	require.False(t, NewObject("x").IsZero())

	require.True(t, NewObject(0).IsZero())
	require.False(t, NewObject(1).IsZero())

	require.True(t, NewObject(false).IsZero())
	require.False(t, NewObject(true).IsZero())
}

func TestObject_IsZero_NilTypes(t *testing.T) {

	// Nil pointers, slices, and maps are all zero
	var nilPointer *testObject
	require.True(t, NewObject(nilPointer).IsZero())
	require.True(t, NewObject([]int(nil)).IsZero())
	require.True(t, NewObject(map[string]int(nil)).IsZero())

	// GOTCHA: an EMPTY but non-nil slice/map is not the zero value.  This follows
	// Go's own definition of zero, and differs from compare.IsZero, which
	// treats any zero-length collection as zero.
	require.False(t, NewObject([]int{}).IsZero())
	require.False(t, NewObject(map[string]int{}).IsZero())
}

func TestObject_IsZero_Interface(t *testing.T) {

	// An untyped nil inside an interface T has no reflected value, and reads as zero
	var value Object[any]
	require.True(t, value.IsZero())

	value.Set(nil)
	require.True(t, value.IsPresent())
	require.True(t, value.IsZero())

	// The boxed value is what gets tested, not the interface
	value.Set(0)
	require.True(t, value.IsZero())

	value.Set("Nien Nunb")
	require.False(t, value.IsZero())
}

/******************************************
 * JSON Marshalling
 ******************************************/

func TestObject_MarshalJSON_Null(t *testing.T) {

	var value Object[testObject]

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `null`, string(result))
}

func TestObject_MarshalJSON_Present(t *testing.T) {

	value := NewObject(testObject{Name: "Chewie", Count: 7})

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `{"name":"Chewie","count":7}`, string(result))
}

func TestObject_MarshalJSON_PresentZero(t *testing.T) {

	// A present zero value marshals as itself, NOT as null
	value := NewObject(testObject{})

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `{"name":"","count":0}`, string(result))
}

func TestObject_MarshalJSON_NilPointer(t *testing.T) {

	// A present-but-nil pointer marshals to "null" (indistinguishable from unset)
	var nilPointer *testObject
	value := NewObject(nilPointer)

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `null`, string(result))
}

func TestObject_MarshalJSON_Error(t *testing.T) {

	// Channels cannot be marshalled into JSON
	value := NewObject(make(chan int))

	result, err := value.MarshalJSON()

	require.NotNil(t, err)
	require.Nil(t, result)
}

func TestObject_UnmarshalJSON_Null(t *testing.T) {

	value := NewObject(testObject{Name: "overwrite me"})

	require.Nil(t, value.UnmarshalJSON([]byte(`null`)))
	require.True(t, value.IsNull())
	require.Equal(t, testObject{}, value.Object())
}

func TestObject_UnmarshalJSON_Empty(t *testing.T) {

	value := NewObject(testObject{Name: "overwrite me"})

	require.Nil(t, value.UnmarshalJSON([]byte(``)))
	require.True(t, value.IsNull())
}

func TestObject_UnmarshalJSON_Present(t *testing.T) {

	var value Object[testObject]

	require.Nil(t, value.UnmarshalJSON([]byte(`{"name":"Rey","count":3}`)))
	require.True(t, value.IsPresent())
	require.Equal(t, testObject{Name: "Rey", Count: 3}, value.Object())
}

func TestObject_UnmarshalJSON_Partial(t *testing.T) {

	// Missing properties fall back to their zero values
	var value Object[testObject]

	require.Nil(t, value.UnmarshalJSON([]byte(`{"name":"Rey"}`)))
	require.True(t, value.IsPresent())
	require.Equal(t, testObject{Name: "Rey", Count: 0}, value.Object())
}

func TestObject_UnmarshalJSON_Error(t *testing.T) {

	var value Object[testObject]

	require.NotNil(t, value.UnmarshalJSON([]byte(`"not-an-object"`)))
	require.True(t, value.IsNull())
}

func TestObject_UnmarshalJSON_ErrorKeepsPriorValue(t *testing.T) {

	// A failed parse leaves the previous value untouched (matching null.Int/Float/Bool)
	value := NewObject(testObject{Name: "keep me", Count: 9})

	require.NotNil(t, value.UnmarshalJSON([]byte(`[1,2,3]`)))
	require.True(t, value.IsPresent())
	require.Equal(t, testObject{Name: "keep me", Count: 9}, value.Object())
}

func TestObject_UnmarshalJSON_Malformed(t *testing.T) {

	var value Object[testObject]

	require.NotNil(t, value.UnmarshalJSON([]byte(`{"name":`)))
	require.True(t, value.IsNull())
}

func TestObject_UnmarshalJSON_PaddedNull(t *testing.T) {

	// GOTCHA: only the exact literal "null" is recognized as null.  A padded
	// null is handed to encoding/json, which decodes it as a present zero value.
	// encoding/json never passes padded literals, so this only affects direct calls.
	var value Object[testObject]

	require.Nil(t, value.UnmarshalJSON([]byte(` null`)))
	require.True(t, value.IsPresent())
	require.Equal(t, testObject{}, value.Object())
}

/******************************************
 * JSON Round Trips
 ******************************************/

// testObjectContainer exercises Object[T] as a struct field
type testObjectContainer struct {
	Struct Object[testObject] `json:"struct"`
	String Object[string]     `json:"string"`
	Slice  Object[[]int]      `json:"slice"`
}

func TestObject_Container_Unmarshal_Empty(t *testing.T) {

	var value testObjectContainer

	require.Nil(t, json.Unmarshal([]byte(`{}`), &value))
	require.True(t, value.Struct.IsNull())
	require.True(t, value.String.IsNull())
	require.True(t, value.Slice.IsNull())
}

func TestObject_Container_Unmarshal_Nulls(t *testing.T) {

	var value testObjectContainer

	require.Nil(t, json.Unmarshal([]byte(`{"struct":null,"string":null,"slice":null}`), &value))
	require.True(t, value.Struct.IsNull())
	require.True(t, value.String.IsNull())
	require.True(t, value.Slice.IsNull())
}

func TestObject_Container_Unmarshal_Full(t *testing.T) {

	var value testObjectContainer

	j := []byte(`{"struct":{"name":"Finn","count":8},"string":"","slice":[1,2,3]}`)

	require.Nil(t, json.Unmarshal(j, &value))

	require.True(t, value.Struct.IsPresent())
	require.Equal(t, testObject{Name: "Finn", Count: 8}, value.Struct.Object())

	// A present empty string is present, not null
	require.True(t, value.String.IsPresent())
	require.Equal(t, "", value.String.Object())

	require.True(t, value.Slice.IsPresent())
	require.Equal(t, []int{1, 2, 3}, value.Slice.Object())
}

func TestObject_Container_Unmarshal_Error(t *testing.T) {

	var value testObjectContainer

	require.NotNil(t, json.Unmarshal([]byte(`{"string":42}`), &value))
	require.True(t, value.String.IsNull())
}

func TestObject_Container_Marshal_Empty(t *testing.T) {

	var value testObjectContainer

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `{"struct":null,"string":null,"slice":null}`, string(result))
}

func TestObject_Container_Marshal_Full(t *testing.T) {

	var value testObjectContainer

	value.Struct.Set(testObject{Name: "Poe", Count: 4})
	value.String.Set("")
	value.Slice.Set([]int{})

	result, err := json.Marshal(value)

	require.Nil(t, err)
	require.Equal(t, `{"struct":{"name":"Poe","count":4},"string":"","slice":[]}`, string(result))
}

func TestObject_Container_RoundTrip(t *testing.T) {

	original := testObjectContainer{
		Struct: NewObject(testObject{Name: "BB-8", Count: 88}),
		Slice:  NewObject([]int{9}),
	}

	marshaled, err := json.Marshal(original)
	require.Nil(t, err)

	var result testObjectContainer
	require.Nil(t, json.Unmarshal(marshaled, &result))
	require.Equal(t, original, result)
}
