package mapof

type String map[string]string

func (x String) GetString(key string) string {
	return x[key]
}

func (x *String) SetString(key string, value string) bool {
	(*x)[key] = value
	return true
}
