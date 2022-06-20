package format

func Date(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

func DateTime(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

func Time(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}
