package index

import "math"

// Normalize returns a new vector scaled to unit length (L2 norm = 1.0).
// It prevents slice aliasing by always returning a newly allocated slice.
func Normalize(v []float32) []float32 {
	var sqSum float64

	// Promote to float64 before squaring to maintain precision during accumulation
	for _, val := range v {
		f := float64(val)
		sqSum += f * f
	}

	mag := float32(math.Sqrt(sqSum))
	res := make([]float32, len(v))

	// Handle zero-magnitude vectors by returning a copy of the original
	if mag == 0 {
		copy(res, v)
		return res
	}

	for i, val := range v {
		res[i] = val / mag
	}

	return res
}

// Dot computes the dot product between two vectors.
// For the purpose of this index, it assumes both vectors are already normalized.
// It panics on length mismatch to catch programming or model-configuration errors early.
func Dot(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("vector length mismatch: check if embedding models are consistent")
	}

	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}

	return sum
}
