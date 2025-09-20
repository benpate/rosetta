package compare

import (
	"io"

	"github.com/benpate/rosetta/convert"
)

// IsZero returns TRUE if the value is the ZERO VALUE for its datatype or NIL
func IsZero(value any) bool {

	if value == nil {
		return true
	}

	if IsNil(value) {
		return true
	}

	switch v := value.(type) {

	case bool:
		return !v
	case string:
		return v == ""
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	case []any:
		return len(v) == 0

	case LengthGetter:
		return v.Length() == 0

	case Nuller:
		return v.IsNull()

	case Inter:
		return IsZero(v.Int())

	case Floater:
		return IsZero(v.Float())

	case Hexer:
		return IsZero(convert.Int64(v))

	case Stringer:
		return IsZero(v.String())

	case io.Reader:
		return IsZero(convert.String(v))
	}

	return false
}

// NotZero returns TRUE if the value is NOT the zero value for its datatype
func NotZero(value any) bool {
	return !IsZero(value)
}
