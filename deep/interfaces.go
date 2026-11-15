package deep

// Copier is implemented by a value that returns a deep copy of itself, which Copy uses
// for a type it does not otherwise know
type Copier interface {
	Copy() (any, error)
}

// Remover is implemented by a value that deletes the named properties from itself, at
// every level, which Remove uses for a type it does not otherwise know
type Remover interface {
	Remove(propertyNames ...string) error
}
