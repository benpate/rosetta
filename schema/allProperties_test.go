package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAllProperties tests the AllProperties function
func TestAllProperties(t *testing.T) {

	s := Schema{
		Element: Object{
			Properties: ElementMap{
				"name":      String{},
				"age":       Integer{BitSize: 32},
				"published": Integer{BitSize: 64},
				"other":     Array{Items: String{}},
				"more": Object{
					Properties: ElementMap{
						"first":  String{},
						"second": Integer{BitSize: 32},
						"third":  Integer{BitSize: 64},
					},
				},
			},
		},
	}

	expected := ElementMap{
		"name":        String{},
		"age":         Integer{BitSize: 32},
		"published":   Integer{BitSize: 64},
		"other":       Array{Items: String{}},
		"more.first":  String{},
		"more.second": Integer{BitSize: 32},
		"more.third":  Integer{BitSize: 64},
	}

	actual := s.AllProperties()

	require.Equal(t, expected, actual)
}
