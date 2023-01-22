package schema

import "github.com/benpate/rosetta/list"

func (schema Schema) Remove(object any, path string) bool {

	return schema.remove(object, list.ByDot(path))
}

func (schema Schema) remove(object any, path list.List) bool {

	head, tail := path.Split()

	// If the list has more than one item, then we need to
	// keep digging into it before we can remove the value.
	if !tail.IsEmpty() {

		if getter, ok := object.(ObjectGetter); ok {

			if child, ok := getter.GetObjectOK(head); ok {
				return schema.remove(child, tail)
			}
		}

		return false
	}

	// Fall through means now it's time to actually
	// delete the thing

	if remover, ok := object.(Remover); ok {
		return remover.Remove(head)
	}

	// Can't delete from this object.  Dernit.
	return false
}
