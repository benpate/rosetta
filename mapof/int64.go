package mapof

type Int64 map[string]int64

func (x Int64) GetInt64(key string) int64 {
	return x[key]
}

func (x *Int64) SetInt64(key string, value int64) bool {
	(*x)[key] = value
	return true
}
