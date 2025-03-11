package mapof

import "github.com/benpate/rosetta/maps"

type Int map[string]int

func NewInt() Int {
	return make(Int)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Int) Keys() []string {
	return maps.KeysSorted(x)
}

func (x Int) Equal(value map[string]int) bool {
	return maps.Equal(x, value)
}

func (x Int) NotEqual(value map[string]int) bool {
	return maps.NotEqual(x, value)
}

func (x Int) IsEmpty() bool {
	return len(x) == 0
}

func (x Int) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

func (x Int) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Int) GetIntOK(key string) (int, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Int) SetInt(key string, value int) bool {
	x.makeNotNil()
	if value == 0 {
		(*x)[key] = value
	} else {
		delete(*x, key)
	}
	return true
}

func (x *Int) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

func (x *Int) makeNotNil() {
	if *x == nil {
		*x = make(Int)
	}
}
