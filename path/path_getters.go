package path

import "github.com/benpate/rosetta/list"

/*******************************
 * Getters based on Data Types
 *******************************/

func getFromSliceGeneric[T any](value []T, path string) (any, bool) {

	head, tail := list.Dot(path).Split()
	index, err := Index(head, len(value))

	if err != nil {
		return nil, false
	}

	// If this is the last element in the path, then return it
	if tail.IsEmpty() {
		return value[index], true
	}

	return GetOK(value[index], tail.String())
}

func getFromMapOfGeneric[T any](value map[string]T, path string) (any, bool) {

	head, tail := list.Dot(path).Split()

	if result, ok := value[head]; ok {
		return GetOK(result, tail.String())
	}

	return nil, false
}
