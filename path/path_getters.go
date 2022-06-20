package path

import (
	"github.com/benpate/rosetta/list"
)

/*******************************
 * Getters based on Data Types
 *******************************/

func getFromSliceOfString(value []string, path string) (interface{}, bool) {

	if path == "" {
		return value, true
	}

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(value))

	if err != nil {
		return nil, false
	}

	// If this is the last element in the path, then return it
	if tail == "" {
		return value[index], true
	}

	// Cannot dig deeper on a string value
	return nil, false
}

func getSliceOfInt(value []int, path string) (interface{}, bool) {

	if path == "" {
		return value, true
	}

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(value))

	if err != nil {
		return nil, false
	}

	// If this is the last element in the path, then return it
	if tail == "" {
		return value[index], true
	}

	// Cannot dig deeper on an int value
	return nil, false
}

func getSliceOfInterface(value []interface{}, path string) (interface{}, bool) {

	if path == "" {
		return value, true
	}

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(value))

	if err != nil {
		return nil, false
	}

	// If this is the last element in the path, then return it
	if tail == "" {
		return value[index], true
	}

	// Otherwise, continue digging.
	return GetOK(value[index], tail)
}

func getSliceOfGetter(value []Getter, path string) (interface{}, bool) {

	if path == "" {
		return value, true
	}

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(value))

	if err != nil {
		return nil, false
	}

	// If this is the last element in the path, then return it
	if tail == "" {
		return value[index], true
	}

	// Otherwise, continue digging.
	return GetOK(value[index], tail)
}

func getMapOfString(value map[string]string, path string) (interface{}, bool) {

	head, tail := list.Split(path, ".")

	if tail != "" {
		return nil, false
	}

	result, ok := value[head]
	return result, ok
}

func getMapOfInterface(value map[string]interface{}, path string) (interface{}, bool) {

	head, tail := list.Split(path, ".")

	if tail == "" {
		result, ok := value[head]
		return result, ok
	}

	return GetOK(value[head], tail)
}
