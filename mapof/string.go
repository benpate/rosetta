package mapof

type String map[string]string

func NewString() String {
	return make(String)
}

func (x String) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x String) GetStringOK(key string) (string, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *String) SetString(key string, value string) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *String) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *String) makeNotNil() {
	if *x == nil {
		*x = make(String)
	}
}
