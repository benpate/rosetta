package convert

import "reflect"

func ReflectValue(value any) reflect.Value {

	if valueOf, ok := value.(reflect.Value); ok {
		return valueOf
	}

	return reflect.ValueOf(value)
}

func ReflectType(value any) reflect.Type {

	if typeOf, ok := value.(reflect.Type); ok {
		return typeOf
	}

	return reflect.TypeOf(value)
}
