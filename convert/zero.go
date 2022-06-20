package convert

import "io"

// IsZeroValue returns TRUE if the value is the zero value for its datatype
func IsZeroValue(value interface{}) bool {

	if value == nil {
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
	case []interface{}:
		return len(v) == 0

	case Nuller:
		return v.IsNull()
	case Inter:
		return IsZeroValue(v.Int())
	case Floater:
		return IsZeroValue(v.Float())
	case Hexer:
		return IsZeroValue(Int64(v))
	case Stringer:
		return IsZeroValue(v.String())
	case io.Reader:
		return IsZeroValue(String(v))
	}

	return false
}
