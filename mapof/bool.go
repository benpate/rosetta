package mapof

type Bool map[string]bool

func (x Bool) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Bool) GetBoolOK(key string) (bool, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Bool) SetBoolOK(key string, value bool) bool {
	if *x == nil {
		*x = make(Bool)
	}
	(*x)[key] = value
	return true
}
