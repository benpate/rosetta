package mapof

import "github.com/benpate/rosetta/maps"

type Int64 map[string]int64

func NewInt64() Int64 {
	return make(Int64)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Int64) Keys() []string {
	return maps.KeysSorted(x)
}

func (x Int64) Equal(value map[string]int64) bool {
	return maps.Equal(x, value)
}

func (x Int64) NotEqual(value map[string]int64) bool {
	return maps.NotEqual(x, value)
}

func (x Int64) IsEmpty() bool {
	return len(x) == 0
}

func (x Int64) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x Int64) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Int64) GetInt64OK(key string) (int64, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int64) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Int64) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Int64) makeNotNil() {
	if *x == nil {
		*x = make(Int64)
	}
}
