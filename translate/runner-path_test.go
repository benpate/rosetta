package translate

import (
	"testing"

	deepcopy "github.com/tiendc/go-deepcopy"
)

func TestDeepCopy_Map(t *testing.T) {

	source := map[string]any{
		"name": "John",
		"age":  30,
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "Anytown",
		},
	}

	var target any

	if err := deepcopy.Copy(&target, source); err != nil {
		t.Errorf("Error during deepcopy: %v", err)
	}

	t.Log(source, target)
}

func TestDeepCopy_Struct(t *testing.T) {

	source := struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	}

	var target any

	if err := deepcopy.Copy(&target, source); err != nil {
		t.Errorf("Error during deepcopy: %v", err)
	}

	t.Log(source, target)
}
