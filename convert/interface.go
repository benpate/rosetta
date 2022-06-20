package convert

import "reflect"

func Interface(value any) any {

	// Safe handling of reflection values
	if v, ok := value.(reflect.Value); ok {

		if v.Kind() == reflect.Invalid {
			return nil
		}

		return v.Interface()
	}

	// Otherwise, just return the value
	return value
}
