package schema

import (
	"testing"

	"github.com/benpate/rosetta/maps"
	"github.com/stretchr/testify/require"
)

func TestObjectComplex(t *testing.T) {

	value := struct {
		FirstName string   `path:"firstName"`
		LastName  string   `path:"lastName"`
		Data      maps.Map `path:"data"`
	}{
		FirstName: "John",
		LastName:  "Connor",
	}

	schema := New(Object{
		Properties: map[string]Element{
			"firstName": String{},
			"lastName":  String{},
			"data": Object{
				Properties: map[string]Element{
					"age":    Integer{},
					"height": Integer{},
					"weight": Integer{},
					"extras": Object{},
				},
			},
		},
	})

	err := schema.Set(&value, "data.extras", maps.Map{"foo": "bar"})
	require.Nil(t, err)
	t.Log(value)

}
