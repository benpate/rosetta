package mapof

type Int64 map[string]int64

func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Int64) GetInt64OK(key string) (int64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int64) SetInt64OK(key string, value int64) bool {
	if *x == nil {
		*x = make(Int64)
	}
	(*x)[key] = value
	return true
}
