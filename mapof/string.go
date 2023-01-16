package mapof

type String map[string]string

func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetStringOK(key string) (string, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *String) SetStringOK(key string, value string) bool {
	(*x)[key] = value
	return true
}
