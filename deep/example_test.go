package deep

import (
	"fmt"
	"strings"
)

// ExampleCopy copies a document, then changes the copy without changing the original
func ExampleCopy() {

	original := map[string]any{
		"type": "Note",
		"to":   []any{"alice"},
	}

	duplicate, err := Copy(original)

	if err != nil {
		fmt.Println(err)
		return
	}

	duplicate["to"].([]any)[0] = "bob"

	fmt.Println(original["to"])
	fmt.Println(duplicate["to"])

	// Output:
	// [alice]
	// [bob]
}

// ExampleRemove deletes two property names from every level of a document
func ExampleRemove() {

	document := map[string]any{
		"type": "Create",
		"bcc":  []any{"carol"},
		"object": map[string]any{
			"type": "Note",
			"bto":  "dave",
		},
	}

	if err := Remove(document, "bto", "bcc"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(document)

	// Output:
	// map[object:map[type:Note] type:Create]
}

// ExampleRemoveFunc deletes every property in a reserved namespace
func ExampleRemoveFunc() {

	document := map[string]any{
		"type":           "Note",
		"emissary:token": "secret",
		"tag":            []any{map[string]any{"name": "go", "x-internal": true}},
	}

	if err := RemoveFunc(document, HasPrefix("emissary:", "x-")); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(document)

	// Output:
	// map[tag:[map[name:go]] type:Note]
}

// ExampleAnyOf combines the Matchers that come with the package with one of the caller's own
func ExampleAnyOf() {

	match := AnyOf(
		Named("bto", "bcc"),
		HasPrefix("x-"),
		func(name string) bool { return strings.HasSuffix(name, "_private") },
	)

	fmt.Println(match("bcc"), match("x-debug"), match("notes_private"), match("to"))

	// Output:
	// true true true false
}
