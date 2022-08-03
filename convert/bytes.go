package convert

// Bytes forces a conversion from an arbitrary value into a slice of bytes.
func Bytes(value interface{}) []byte {

	if value == nil {
		return []byte{}
	}

	switch v := value.(type) {
	case byte:
		return []byte{v}

	case []byte:
		return v

	case string:
		return []byte(v)

	case Stringer:
		return []byte(v.String())
	}

	return []byte(String(value))
}
