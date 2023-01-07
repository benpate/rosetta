package mapof

type Int map[string]int

func (x Int) GetInt(key string) int {
	return x[key]
}

func (x *Int) SetInt(key string, value int) bool {
	(*x)[key] = value
	return true
}
