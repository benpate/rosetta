package mapof

type Int map[string]int

func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetIntOK(key string) (int, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int) SetIntOK(key string, value int) bool {
	if *x == nil {
		*x = make(Int)
	}
	(*x)[key] = value
	return true
}
