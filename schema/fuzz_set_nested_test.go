package schema_test

import (
	"strings"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
)

// fuzzNestedSchema returns a small nested schema, with String leaves, Any wildcards, and a list,
// that drives Set down real paths into a mapof.Any
func fuzzNestedSchema() schema.Schema {

	return schema.New(schema.Object{
		Wildcard: schema.Any{},
		Properties: schema.ElementMap{
			"name": schema.String{},
			"age":  schema.String{},
			"tags": schema.Array{Items: schema.Object{Wildcard: schema.Any{}}, MaxLength: 16},
			"child": schema.Object{
				Properties: schema.ElementMap{
					"name": schema.String{},
					"grandchild": schema.Object{
						Properties: schema.ElementMap{
							"name": schema.String{},
						},
					},
				},
			},
		},
	})
}

// FuzzSet_NestedPaths drives arbitrary dotted paths through Set into a mapof.Any.  Path segments
// come from whatever a client posted, so a hostile path must never panic, hang, or allocate
// without bound, and a write that Set accepts must read back through the same path.
func FuzzSet_NestedPaths(f *testing.F) {

	f.Add("name")
	f.Add("age")
	f.Add("0.0")
	f.Add("child.name")
	f.Add("child.grandchild.name")
	f.Add("tags.0.label")
	f.Add("tags.15.label")
	f.Add("tags.99999999.label")
	f.Add("other.0.id")
	f.Add("other.99999999.id")
	f.Add("other.007.id")
	f.Add("")
	f.Add(".")
	f.Add("..")
	f.Add(".name")
	f.Add("name.")
	f.Add("child..name")
	f.Add("nonexistent.path")
	f.Add("child.nonexistent")
	f.Add(strings.Repeat("child.", 64) + "name")
	f.Add(strings.Repeat("0.", 64) + "x")
	f.Add("\x00")
	f.Add("child.\x00.name")

	s := fuzzNestedSchema()

	f.Fuzz(func(t *testing.T, path string) {

		object := mapof.NewAny()

		// A hostile path may legitimately be REFUSED, but must never panic or hang
		if err := s.Set(&object, path, "VALUE"); err != nil {
			return
		}

		// RULE: A write that Set accepts must be readable back through that same path.
		// Silently accepting a write that goes nowhere is how a form field vanishes.
		result, err := s.Get(&object, path)

		if err != nil {
			t.Fatalf("Set accepted path %q, but Get could not read it back: %v", path, err)
		}

		if result != "VALUE" {
			t.Fatalf("Set accepted path %q, but Get read back %#v", path, result)
		}
	})
}
