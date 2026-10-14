package bcp47

// maxTagLength is the longest language tag that Parse will read.  Real tags, even with
// a script, region, and variant, fit in well under half of this.
const maxTagLength = 64
