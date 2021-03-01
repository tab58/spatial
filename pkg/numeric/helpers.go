package numeric

import (
	"math"
)

// IsOverflow returns true if the number has overflowed, false if not.
func IsOverflow(num float64) bool {
	return math.IsInf(num, 0)
}

// AreAnyOverflow returns true if any number has overflowed, false if not.
func AreAnyOverflow(nums ...float64) bool {
	for _, num := range nums {
		if IsOverflow(num) {
			return true
		}
	}
	return false
}

// Signum returns the sign of the float64 provided.
func Signum(a float64) (int, error) {
	if math.IsNaN(a) {
		return 0, ErrNaN
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
