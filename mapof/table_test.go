package mapof

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/schema"
)

type testTableItem struct {
	key   string
	value any
}

func testTable(schema schema.Schema, object any, items []testTableItem) error {

	for _, item := range items {

		if err := schema.Set(object, item.key, item.value); err != nil {
			return err
		}

		if value, err := schema.Get(object, item.key); err != nil {
			return err
		} else if value != item.value {
			return derp.InternalError("mapof.testTable", "Unexpected value", item.key, item.value, value)
		}
	}

	return nil
}
