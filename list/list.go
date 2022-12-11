package list

// IsEmpty returns TRUE if the list is empty.
func IsEmpty[T Stringlike](value T) bool {
	return len(value) == 0
}

// IsEmptyTail returns TRUE if this list only has one element
func IsEmptyTail[T Stringlike](value T, delimiter byte) bool {
	place := Index(value, delimiter)
	return (place == -1) || (place == len(value)-1)
}

// Index finds the first occurrance of the delimiter (-1 if not found)
func Index[T Stringlike](value T, delimiter byte) int {
	for i := 0; i < len(value); i++ {
		if value[i] == delimiter {
			return i
		}
	}
	return -1
}

// LastIndex finds the last occurrance of the delimiter (-1 if not found)
func LastIndex[T Stringlike](value T, delimiter byte) int {
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == delimiter {
			return i
		}
	}
	return -1

}

// Head returns the FIRST item in a list
func Head[T Stringlike](value T, delimiter byte) string {

	index := Index(value, delimiter)

	if index == -1 {
		return string(value)
	}

	return string(value[:index])
}

// Tail returns any values in the list AFTER the first item
func Tail[T Stringlike](value T, delimiter byte) T {
	index := Index(value, delimiter)

	if index == -1 {
		return T("")
	}

	return value[index+1:]
}

// RemoveLast returns the full list, with the last element removed.
func RemoveLast[T Stringlike](value T, delimiter byte) T {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return T("")
	}

	return value[:index]
}

// First returns the FIRST item in a list (alias for Head)
func First[T Stringlike](value T, delimiter byte) string {
	return Head(value, delimiter)
}

// Last returns the LAST item in a T-based-list
func Last[T Stringlike](value T, delimiter byte) string {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return string(value)
	}

	return string(value[index+1:])
}

// Split returns the FIRST element, and the REST element in one function call
func Split[T Stringlike](value T, delimiter byte) (string, T) {

	index := Index(value, delimiter)

	if index == -1 {
		return string(value), T("")
	}

	return string(value[:index]), value[index+1:]
}

// SplitTail behaves like split, but splits the beginning of the list from the last item in the list.  So, the list "a,b,c" => "a,b", "c"
func SplitTail[T Stringlike](value T, delimiter byte) (T, string) {

	index := LastIndex(value, delimiter)

	if index == -1 {
		return value, ""
	}

	return value[:index], string(value[index+1:])
}

// at returns the list vaue at a particular index
func At[T Stringlike](value T, delimiter byte, index int) string {

	if IsEmpty(value) {
		return ""
	}

	if index == 0 {
		return Head(value, delimiter)
	}

	tail := Tail(value, delimiter)

	return At(tail, delimiter, index-1)
}

// PushHead adds a new item to the beginning of the list
func PushHead[T Stringlike](value T, headValue string, delimiter byte) T {

	// If the new value is empty, NOOP
	if len(headValue) == 0 {
		return value
	}

	// If the value is empty, make a copy in a new variable
	if len(value) == 0 {
		return T(headValue)
	}

	return T(headValue + string(delimiter) + string(value))
}

// PushTail adds a new item to the end of the list
func PushTail[T Stringlike](value T, tailValue string, delimiter byte) T {

	// If the new value is empty, return the value
	if len(tailValue) == 0 {
		return value
	}

	// If the value is empty, make a copy in a new variable
	if len(value) == 0 {
		return T(tailValue)
	}

	return T(string(value) + string(delimiter) + tailValue)
}
