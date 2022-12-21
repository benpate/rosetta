package schema

type ElementMap map[string]Element

// Nullable interface wraps the IsNull method, that helps an object
// to identify if it contains a null value or not.  This mirrors
// the null.Nullable interface here, for convenience.
type Nullable interface {
	IsNull() bool
}

type Enumerator interface {
	Enumerate() []string
}

type Getter interface {
	Get(string) (any, error)
}

type Setter interface {
	Set(string, any) error
}
