package deep

// Copier is implemented by a value that returns a deep copy of itself, which Copy uses for a
// type it does not otherwise know. It copies what it holds with CopyWithDepth and depth+1.
type Copier interface {
	Copy(depth int) (any, error)
}

// Remover is implemented by a value that deletes the named properties from itself, which Remove
// uses for a type it does not otherwise know. It searches what it holds with RemoveWithDepth
// and depth+1.
type Remover interface {
	Remove(depth int, propertyNames ...string) error
}
