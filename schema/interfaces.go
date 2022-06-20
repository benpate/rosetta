package schema

// Nullable interface wraps the IsNull method, that helps an object
// to identify if it contains a null value or not.  This mirrors
// the null.Nullable interface here, for convenience.
type Nullable interface {
	IsNull() bool
}

type Enumerator interface {
	Enumerate() []string
}

type PathGetter interface {
	GetPath(path string) (any, error)
}

type PathSetter interface {
	SetPath(path string, value any) error
}

type ElementMap map[string]Element
