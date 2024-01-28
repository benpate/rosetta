package channel

import (
	"testing"
)

func TestMap(t *testing.T) {

	ch := Map(testChannel(), func(value string) string {
		return value + "!"
	})

	for value := range ch {
		t.Log(value)
	}

	t.Log("DONE")
}
