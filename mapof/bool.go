package mapof

type Bool map[string]bool

func (x Bool) GetBool(key string) bool {
	return x[key]
}

func (x *Bool) SetBool(key string, value bool) bool {
	(*x)[key] = value
	return true
}
