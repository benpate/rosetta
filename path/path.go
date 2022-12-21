package path

import (
	"reflect"
	"strconv"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// SetAll adds every path from the dataset into the object.  It returns an aggregate error containing all errors generated.
func SetAll(object any, dataset map[string]any) error {

	var errorReport error

	// Put approved form data into the stream
	for key, value := range dataset {
		if err := Set(object, key, value); err != nil {
			errorReport = derp.Append(errorReport, err)
		}
	}

	return errorReport
}

// Get tries to return the value of the object at this path.
func Get(object any, path string) any {

	result, _ := GetOK(object, path)
	return result
}

// GetOK returns the value of the object at the provided path.
// If a value does not already exist, then the OK boolean is false.
func GetOK(object any, path string) (any, bool) {

	// If the object is empty, then there's nothing left to traverse.
	if object == nil {
		return nil, false
	}

	// If the path is empty, then we have arrived at the correct value.
	if path == "" {
		return object, true
	}

	// Next steps depend on the type of object we're working with.
	switch obj := object.(type) {

	case Getter:
		return obj.GetPath(path)

	case []Getter:
		return getFromSliceGeneric(obj, path)

	case []string:
		return getFromSliceGeneric(obj, path)

	case []int:
		return getFromSliceGeneric(obj, path)

	case []int64:
		return getFromSliceGeneric(obj, path)

	case []float64:
		return getFromSliceGeneric(obj, path)

	case []any:
		return getFromSliceGeneric(obj, path)

	case map[string]Getter:
		return getFromMapOfGeneric(obj, path)

	case map[string]string:
		return getFromMapOfGeneric(obj, path)

	case map[string]int:
		return getFromMapOfGeneric(obj, path)

	case map[string]int64:
		return getFromMapOfGeneric(obj, path)

	case map[string]float64:
		return getFromMapOfGeneric(obj, path)

	case map[string]any:
		return getFromMapOfGeneric(obj, path)

	default:
		return GetWithReflection(reflect.ValueOf(obj), path)
	}
}

// Set tries to return the value of the object at this path.
func Set(object any, name string, value any) error {

	switch obj := object.(type) {

	case Setter:
		return obj.SetPath(name, value)

	case []Setter:
		return setSliceOfSetter(name, obj, value)

	case []string:
		return setSliceOfString(name, obj, value)

	case []int:
		return setSliceOfInt(name, obj, value)

	case []any:
		return setSliceOfInterface(name, obj, value)

	case map[string]string:
		return setMapOfString(name, obj, value)

	case map[string]any:
		return setMapOfInterface(name, obj, value)

	default:
		return SetWithReflection(reflect.ValueOf(object), name, value)
	}
}

// Delete tries to remove a value from ths object at this path
func Delete(object any, name string) error {

	switch obj := object.(type) {

	case Deleter:
		return obj.DeletePath(name)

	case []string:
		return deleteSliceOfString(name, obj)

	case []int:
		return deleteSliceOfInt(name, obj)

	case []any:
		return deleteSliceOfInterface(name, obj)

	case []Deleter:
		return deleteSliceOfDeleter(name, obj)

	case map[string]string:
		return deleteMapOfString(name, obj)

	case map[string]any:
		return deleteMapOfInterface(name, obj)
	}

	return derp.NewInternalError("path.Delete", "Unable to delete from this type of record.")
}

// Index is useful for vetting array indices.  It attempts to convert the Head() token int
// an integer, and then check that the integer is within the designated array bounds (is greater than zero,
// and less than the maximum value provided to the function).
//
// It returns the array index and an error
func Index(value string, maximum int) (int, error) {

	result, err := strconv.Atoi(value)

	if err != nil {
		return 0, derp.Wrap(err, "path.Index", "Index must be an integer", value, maximum)
	}

	if result < 0 {
		return 0, derp.NewInternalError("path.Index", "Index out of bounds", "cannot be less than zero", value)
	}

	if (maximum >= 0) && (result >= maximum) {
		return 0, derp.NewInternalError("path.Index", "Index out of bounds", "cannot be greater than (or equal to) maximum", value, maximum)
	}

	// Fall through means that this is a valid array index
	return result, nil
}

// Split splits the path into head and tail strings (separated by ".")
func Split(path string) (string, string) {
	head, tail := list.Dot(path).Split()
	return head, tail.String()
}
