package schema

// getLength returns the length of an object, if it is an ArrayGetter
func getLength(object any) (int, bool) {

	if getter, ok := object.(LengthGetter); ok {
		return getter.Length(), true
	}

	return 0, false
}

// getIndex returns the value at a specific index, if the object is an ArrayGetter
func getIndex(object any, index int) (any, bool) {

	if getter, ok := object.(ArrayGetter); ok {
		return getter.GetIndex(index)
	}

	return nil, false
}
