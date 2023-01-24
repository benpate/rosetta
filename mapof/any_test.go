package mapof

import (
	"testing"

	"github.com/benpate/rosetta/schema"
	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {

	s := schema.New(schema.Object{
		Wildcard: schema.Object{
			Wildcard: schema.String{},
		},
	})

	v := NewAny()

	require.NotNil(t, s.Set(&v, "foo.bar", "baz"))
}
