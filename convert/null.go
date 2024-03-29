package convert

import "github.com/benpate/rosetta/null"

// NullBool converts a value into a nullable value.
// The value is only set if the input value is a natural match for this data type.
func NullBool(value any) null.Bool {

	var result null.Bool

	if v, ok := BoolOk(value, false); ok {
		result.Set(v)
	}

	return result
}

// NullInt converts a value into a nullable value.
// The value is only set if the input value is a natural match for this data type.
func NullInt(value any) null.Int {

	var result null.Int

	if v, ok := IntOk(value, 0); ok {
		result.Set(v)
	}

	return result

}

// NullInt64 converts a value into a nullable value.
// The value is only set if the input value is a natural match for this data type.
func NullInt64(value any) null.Int64 {

	var result null.Int64

	if v, ok := Int64Ok(value, 0); ok {
		result.Set(v)
	}

	return result

}

// NullFloat converts a value into a nullable value.
// The value is only set if the input value is a natural match for this data type.
func NullFloat(value any) null.Float {

	var result null.Float

	if v, ok := FloatOk(value, 0); ok {
		result.Set(v)
	}

	return result

}
