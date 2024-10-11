package convert

import (
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestURLValues(t *testing.T) {
	original := url.Values{}
	original.Add("name1", "value1")
	original.Add("name2", "value2")
	original.Add("name3", "value3")
	original.Add("name3", "another value3")

	mapOfAny := MapOfAny(original)
	result := URLValues(mapOfAny)

	// require.Equal(t, original, result)
	spew.Dump(original, mapOfAny, result)
}

func TestURLValues2(t *testing.T) {
	original := url.Values{}
	original.Add("name1", "value1")
	original.Add("name2", "value2")
	original.Add("name3", "value3")
	original.Add("name3", "another value3")

	converted := URLValues(original)

	require.Equal(t, original, converted)
}
