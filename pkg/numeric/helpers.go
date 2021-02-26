package numeric

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/errors"
)

// Signum returns the sign of the float64 provided.
func Signum(a float64) (int, error) {
	if math.IsNaN(a) {
		return 0, errors.ErrNaN
	}
	if a < 0 || math.IsInf(a, -1) {
		return -1, nil
	} else if a > 0 || math.IsInf(a, 1) {
		return 1, nil
	}
	return 0, nil
}

// Nrm2 computes the 2-norm of a vector in a numerically-stable way.
func Nrm2(a float64, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	x := math.Abs(a)
	y := math.Abs(b)
	u := math.Max(x, y)
	t := math.Min(x, y) / u
	return u * math.Sqrt(1+t*t)
}
