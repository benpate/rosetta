package format

// UnsafeAny leaves the string format untouched.
func UnsafeAny(arg string) StringFormat {
	return func(value string) (string, error) {
		return value, nil
	}
}
