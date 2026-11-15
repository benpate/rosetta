package deep

import (
	"fmt"
	"testing"
)

// benchmarkSink keeps each result alive, so that the compiler cannot delete the work
var benchmarkSink any

// benchmarkListed returns a small ActivityPub document made of the types that Copy and
// Remove list by name
func benchmarkListed() map[string]any {
	return map[string]any{
		"type":      "Create",
		"to":        []any{"a", "b"},
		"published": 1.5,
		"object": map[string]any{
			"type":    "Note",
			"content": "hello",
			"tag": []any{
				map[string]any{"type": "Mention", "href": "h"},
				map[string]any{"type": "Hashtag", "name": "n"},
			},
		},
	}
}

// benchmarkReflected returns the same document made of types that are reached by reflection
func benchmarkReflected() testDocument {
	return testDocument{
		"type":      "Create",
		"to":        testArray{"a", "b"},
		"published": 1.5,
		"object": testDocument{
			"type":    "Note",
			"content": "hello",
			"tag": testArray{
				testDocument{"type": "Mention", "href": "h"},
				testDocument{"type": "Hashtag", "name": "n"},
			},
		},
	}
}

// benchmarkStrings returns a map of 100 strings, which Remove asks its Matcher about one by one
func benchmarkStrings() map[string]string {

	result := make(map[string]string, 100)

	for index := range 100 {
		result[fmt.Sprint("key", index)] = "value"
	}

	return result
}

// BenchmarkCopy_Listed measures a copy of the types that Copy lists by name
func BenchmarkCopy_Listed(b *testing.B) {

	for document := benchmarkListed(); b.Loop(); {
		benchmarkSink, _ = Copy(document)
	}
}

// BenchmarkCopy_Reflected measures a copy of the types that Copy reaches by reflection
func BenchmarkCopy_Reflected(b *testing.B) {

	for document := benchmarkReflected(); b.Loop(); {
		benchmarkSink, _ = Copy(document)
	}
}

// BenchmarkRemove_Listed measures a search of the types that Remove lists by name
func BenchmarkRemove_Listed(b *testing.B) {

	for document := benchmarkListed(); b.Loop(); {
		benchmarkSink = Remove(document, "bto", "bcc")
	}
}

// BenchmarkRemove_Reflected measures a search of the types that Remove reaches by reflection
func BenchmarkRemove_Reflected(b *testing.B) {

	for document := benchmarkReflected(); b.Loop(); {
		benchmarkSink = Remove(document, "bto", "bcc")
	}
}

// BenchmarkRemove_LargeStringMap measures a search of one large map of simple values
func BenchmarkRemove_LargeStringMap(b *testing.B) {

	for document := benchmarkStrings(); b.Loop(); {
		benchmarkSink = Remove(document, "bto", "bcc")
	}
}
