package schema

import (
	"testing"

	"github.com/benpate/rosetta/sliceof"
	"github.com/stretchr/testify/require"
)

func TestAppend_String(t *testing.T) {

	testArray := sliceof.String{"one", "two", "three"}
	testElement := testArrayA_Schema()
	testSchema := New(testElement)

	err := testSchema.Append(&testArray, "", "four")
	require.Nil(t, err)

	require.Equal(t, 4, len(testArray))
	require.Equal(t, "one", testArray[0])
	require.Equal(t, "two", testArray[1])
	require.Equal(t, "three", testArray[2])
	require.Equal(t, "four", testArray[3])
}

func TestAppend_Object(t *testing.T) {

	testArray := sliceof.Object[string]{"one", "two", "three"}
	testElement := testArrayA_Schema()
	testSchema := New(testElement)

	err := testSchema.Append(&testArray, "", "four")

	require.Nil(t, err)
	require.Equal(t, 4, len(testArray))
	require.Equal(t, "one", testArray[0])
	require.Equal(t, "two", testArray[1])
	require.Equal(t, "three", testArray[2])
	require.Equal(t, "four", testArray[3])
}
