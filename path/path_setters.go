package path

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
)

func setSliceOfString(path string, object []string, value interface{}) error {

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(object))

	if err != nil {
		return err
	}

	if tail == "" {
		object[index] = convert.String(value)
		return nil
	}

	return derp.NewInternalError("path.Set", "Invalid Path", path)
}

func setSliceOfInt(path string, object []int, value interface{}) error {

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(object))

	if err != nil {
		return err
	}

	if tail == "" {
		object[index] = convert.Int(value)
		return nil
	}

	return derp.NewInternalError("path.Set", "Invalid Path", path)
}

func setSliceOfInterface(path string, object []interface{}, value interface{}) error {

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(object))

	if err != nil {
		return err
	}

	if tail == "" {
		object[index] = value
		return nil
	}

	return Set(&object[index], tail, value)
}

func setSliceOfSetter(path string, object []Setter, value interface{}) error {

	head, tail := list.Split(path, ".")
	index, err := Index(head, len(object))

	if err != nil {
		return err
	}

	if tail == "" {
		if setter, ok := value.(Setter); ok {
			object[index] = setter
			return nil
		}

		return derp.NewInternalError("path.Set", "Value is not a setter", value)
	}

	return Set(object[index], tail, value)
}

func setMapOfInterface(path string, object map[string]interface{}, value interface{}) error {

	head, tail := list.Split(path, ".")

	if tail == "" {
		object[head] = value
		return nil
	}

	return Set(object[head], tail, value)
}

func setMapOfString(path string, object map[string]string, value interface{}) error {

	head, tail := list.Split(path, ".")

	if tail != "" {
		return derp.NewInternalError("path.Set", "Cannot set sub-properties of a string", path)
	}

	object[head] = convert.String(value)
	return nil
}
