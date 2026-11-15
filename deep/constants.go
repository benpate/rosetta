package deep

// maxDepth is the deepest level of nesting that Copy and Remove walk before they return an
// error, which is also the deepest document that encoding/json will read
const maxDepth = 10_000
