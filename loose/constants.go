package loose

// templateCacheLimit is the most templates that CachedTemplate keeps compiled. A template beyond
// it is compiled on every call instead.
const templateCacheLimit = 1000
