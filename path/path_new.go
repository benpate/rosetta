package path

import (
	"github.com/benpate/rosetta/list"
)

/*************************************************
 * New Style Getters
 *************************************************/

func GetBool(object any, path string) bool {

	if leaf, last, ok := getLeaf(object, list.Dot(path)); ok {
		if getter, ok := leaf.(BoolGetter); ok {
			return getter.GetBool(last)
		}
	}

	return false
}

func GetFloat(object any, path string) float64 {

	if leaf, last, ok := getLeaf(object, list.Dot(path)); ok {
		if getter, ok := leaf.(FloatGetter); ok {
			return getter.GetFloat(last)
		}
	}

	return 0
}

func GetInt(object any, path string) int {

	if leaf, last, ok := getLeaf(object, list.Dot(path)); ok {
		if getter, ok := leaf.(IntGetter); ok {
			return getter.GetInt(last)
		}
	}

	return 0
}

func GetInt64(object any, path string) int64 {

	if leaf, last, ok := getLeaf(object, list.Dot(path)); ok {
		if getter, ok := leaf.(Int64Getter); ok {
			return getter.GetInt64(last)
		}
	}

	return 0
}

func GetString(object any, path string) string {

	if leaf, last, ok := getLeaf(object, list.Dot(path)); ok {
		if getter, ok := leaf.(StringGetter); ok {
			return getter.GetString(last)
		}
	}

	return ""
}

/*************************************************
 * New Style Setters
 *************************************************/

func SetBool(object any, path string, value bool) bool {

	leaf, last, ok := getLeaf(object, list.Dot(path))

	if !ok {
		return false
	}

	if setter, ok := leaf.(BoolSetter); ok {
		return setter.SetBool(last, value)
	}

	return false
}

func SetFloat(object any, path string, value float64) bool {

	leaf, last, ok := getLeaf(object, list.Dot(path))

	if !ok {
		return false
	}

	if setter, ok := leaf.(FloatSetter); ok {
		return setter.SetFloat(last, value)
	}

	return false
}

func SetInt(object any, path string, value int) bool {

	leaf, last, ok := getLeaf(object, list.Dot(path))

	if !ok {
		return false
	}

	if setter, ok := leaf.(IntSetter); ok {
		return setter.SetInt(last, value)
	}

	return false
}

func SetInt64(object any, path string, value int64) bool {

	leaf, last, ok := getLeaf(object, list.Dot(path))

	if !ok {
		return false
	}

	if setter, ok := leaf.(Int64Setter); ok {
		return setter.SetInt64(last, value)
	}

	return false
}

func SetString(object any, path string, value string) bool {

	leaf, last, ok := getLeaf(object, list.Dot(path))

	if !ok {
		return false
	}

	if setter, ok := leaf.(StringSetter); ok {
		return setter.SetString(last, value)
	}

	return false
}

/*************************************************
 * Tree Traversal
 *************************************************/

func getLeaf(object any, path list.List) (any, string, bool) {

	head, tail := path.Split()

	if tail.IsEmpty() {
		return object, head, true
	}

	if childGetter, ok := object.(ChildGetter); ok {
		child, ok := childGetter.GetChild(head)
		if ok {
			return getLeaf(child, tail)
		}
	}

	return nil, "", false
}
