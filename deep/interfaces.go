package deep

// Copier is implemented by a value that returns a deep copy of itself. It copies what it
// holds with CopyWithDepth and depth+1.
type Copier interface {
	Copy(depth int) (any, error)
}

// Remover is implemented by a value that deletes the matching properties from itself. It
// searches what it holds with RemoveWithDepth and depth+1.
type Remover interface {
	Remove(depth int, match Matcher) error
}
