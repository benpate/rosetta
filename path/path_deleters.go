package path

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

func deleteSliceOfString(path string, object []string) error {
	return derp.NewInternalError("path.deleteSliceOfString", "Unimplemented")
}

func deleteSliceOfInt(path string, object []int) error {
	return derp.NewInternalError("path.deleteSliceOfString", "Unimplemented")
}

func deleteSliceOfDeleter(path string, object []Deleter) error {
	return derp.NewInternalError("path.deleteSliceOfString", "Unimplemented")
}

func deleteSliceOfInterface(path string, object []interface{}) error {
	return derp.NewInternalError("path.deleteSliceOfString", "Unimplemented")
}

func deleteMapOfString(path string, object map[string]string) error {

	head, tail := list.Split(path, ".")

	if tail != "" {
		return derp.NewInternalError("path.deleteMapOfString", "Cannot delete sub-elements of string", path)
	}

	delete(object, head)
	return nil
}

func deleteMapOfInterface(name string, object map[string]interface{}) error {

	head, tail := list.Split(name, ".")

	if tail != "" {
		return Delete(object[head], tail)
	}

	delete(object, head)
	return nil
}
