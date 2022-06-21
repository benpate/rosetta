package list

// IsEmpty returns TRUE if the list is empty.
func IsEmpty(value []byte) bool {
	return len(value) == 0
}

// IsEmptyTail returns TRUE if this list only has one element
func IsEmptyTail(value []byte, delimiter byte) bool {
	place := Index(value, delimiter)
	return (place == -1) || (place == len(value)-1)
}

// Index finds the first occurrance of the delimiter (-1 if not found)
func Index(value []byte, delimiter byte) int {
	for i := 0; i < len(value); i++ {
		if value[i] == delimiter {
			return i
		}
	}
	return -1
}

// LastIndex finds the last occurrance of the delimiter (-1 if not found)
func LastIndex(value []byte, delimiter byte) int {
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == delimiter {
			return i
		}
	}
	return -1

}

// Head returns the FIRST item in a list
func Head(value []byte, delimiter byte) []byte {

	index := Index(value, delimiter)

	if index == -1 {
		return []byte(value)
	}

	return []byte(value[:index])
}

// Tail returns any values in the list AFTER the first item
func Tail(value []byte, delimiter byte) []byte {
	index := Index(value, delimiter)

	if index == -1 {
		return []byte{}
	}

	return value[index+1:]
}

// RemoveLast returns the full list, with the last element removed.
func RemoveLast(value []byte, delimiter byte) []byte {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return []byte{}
	}

	return value[:index]
}

// Last returns the LAST item in a string-based-list
func Last(value []byte, delimiter byte) []byte {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return value
	}

	return value[index+1:]
}

// Split returns the FIRST element, and the REST element in one function call
func Split(value []byte, delimiter byte) ([]byte, []byte) {

	index := Index(value, delimiter)

	if index == -1 {
		return value, []byte{}
	}

	return value[:index], value[index+1:]
}

// SplitTail behaves like split, but splits the beginning of the list from the last item in the list.  So, the list "a,b,c" => "a,b", "c"
func SplitTail(value []byte, delimiter byte) ([]byte, []byte) {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return value, []byte{}
	}

	return value[:index], value[index+1:]
}

// at returns the list vaue at a particular index
func At(value []byte, delimiter byte, index int) []byte {

	if IsEmpty(value) {
		return value
	}

	if index == 0 {
		return Head(value, delimiter)
	}

	tail := Tail(value, delimiter)

	return At(tail, delimiter, index-1)
}

// PushHead adds a new item to the beginning of the list
func PushHead(value []byte, newValue []byte, delimiter byte) []byte {

	// If the new value is empty, NOOP
	if len(newValue) == 0 {
		if len(value) == 0 {
			return []byte{}
		}
		return value
	}

	// If the value is empty, make a copy in a new variable
	if len(value) == 0 {
		result := make([]byte, len(newValue))
		copy(result, newValue)
		return result
	}

	// Otherwise, make a new variable with the new value at the HEAD of the value
	newLength := len(newValue) + 1 + len(value)

	result := make([]byte, newLength)

	copy(result, newValue)
	result[len(newValue)] = delimiter
	copy(result[len(newValue)+1:], value)

	return result
}

// PushTail adds a new item to the end of the list
func PushTail(value []byte, newValue []byte, delimiter byte) []byte {

	// If the new value is empty, return the value
	if len(newValue) == 0 {
		if len(value) == 0 {
			return []byte{}
		}
		return value
	}

	// If the value is empty, make a copy in a new variable
	if len(value) == 0 {
		result := make([]byte, len(newValue))
		copy(result, newValue)
		return result
	}

	// Otherwise, make a new variable with the newValue at the end
	newLength := len(value) + 1 + len(newValue)

	result := make([]byte, newLength)

	copy(result, value)
	result[len(value)] = delimiter
	copy(result[len(value)+1:], newValue)

	return result
}
