package format

func IPv4(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

func IPv6(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

func Hostname(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

func URI(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}
