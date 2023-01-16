package path

/**************************
 * Getter Interfaces
 **************************/

type BoolGetter interface {
	GetBool(string) bool
}

type BytesGetter interface {
	GetBytes(string) []byte
}

type FloatGetter interface {
	GetFloat(string) float64
}

type IntGetter interface {
	GetInt(string) int
}

type Int64Getter interface {
	GetInt64(string) int64
}

type StringGetter interface {
	GetString(string) string
}

type ObjectGetter interface {
	GetObject(string) (any, bool)
}

/**************************
 * Setter Interfaces
 **************************/

type BoolSetter interface {
	SetBool(string, bool) bool
}

type BytesSetter interface {
	SetBytes(string, []byte) bool
}

type FloatSetter interface {
	SetFloat(string, float64) bool
}

type IntSetter interface {
	SetInt(string, int) bool
}

type Int64Setter interface {
	SetInt64(string, int64) bool
}

type StringSetter interface {
	SetString(string, string) bool
}

/**************************
 * Getter/Setter Interfaces
 **************************/

type BoolGetterSetter interface {
	BoolGetter
	BoolSetter
}

type BytesGetterSetter interface {
	BytesGetter
	BytesSetter
}

type FloatGetterSetter interface {
	FloatGetter
	FloatSetter
}

type IntGetterSetter interface {
	IntGetter
	IntSetter
}

type Int64GetterSetter interface {
	Int64Getter
	Int64Setter
}

type StringGetterSetter interface {
	StringGetter
	StringSetter
}

/**************************
 * Deleter Interface
 **************************/
type Deleter interface {
	Delete(string) bool
}
