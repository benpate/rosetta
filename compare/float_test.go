package compare

import "testing"

func TestCompareFloat(t *testing.T) {
	assertCompare(t, Float32, float32(1.5), float32(2.5))
	assertCompare(t, Float64, float64(1.5), float64(2.5))
}

func TestCompareFloat_Negatives(t *testing.T) {
	assertCompare(t, Float64, -2.5, -1.5)
	assertCompare(t, Float32, float32(-2.5), float32(-1.5))
}
