package schema

import (
	"testing"
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

	result := s.AllProperties()

	t.Log(result)
}
