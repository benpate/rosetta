package mapof

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/maps"
)

type String map[string]string

func NewString() String {
	return make(String)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x String) Keys() []string {
	return maps.KeysSorted(x)
}

func (x String) Equal(value map[string]string) bool {
	// Lengths must be identical
	if len(x) != len(value) {
		return false
	}

	// Items at each index must be identical
	for key := range x {
		if x[key] != value[key] {
			return false
		}
	}

	return true
}

func (x String) NotEqual(value map[string]string) bool {
	return !x.Equal(value)
}

func (x String) IsEmpty() bool {
	return len(x) == 0
}

func (x String) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter/Setter Interfaces
 ******************************************/

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

	if value == "" {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
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

func (x String) MapOfAny() map[string]any {
	return convert.MapOfAny(x.MapOfString())
}

func (x String) MapOfString() map[string]string {
	return x
}
