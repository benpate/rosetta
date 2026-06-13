package format

// IPv4 returns a StringFormat intended to validate IPv4 addresses; validation is not yet implemented.
func IPv4(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

// IPv6 returns a StringFormat intended to validate IPv6 addresses; validation is not yet implemented.
func IPv6(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

// Hostname returns a StringFormat intended to validate hostnames; validation is not yet implemented.
func Hostname(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}

// URI returns a StringFormat intended to validate URIs; validation is not yet implemented.
func URI(arg string) StringFormat {

	return func(value string) (string, error) {
		return value, nil
	}
}
