package geometry

import (
	"math"
	"math/big"

	"github.com/tab58/v1/spatial/pkg/errors"
)

// BigFloatPrecision is the number of bits every default *big.Float should have in the mantissa.
const BigFloatPrecision = 80

// NewBigFloat creates a new *big.Float with specific parameters.
func NewBigFloat(z float64) *big.Float {
	return new(big.Float).SetPrec(BigFloatPrecision).SetFloat64(z)
}

// IsInvalidTolerance returns true if the tolerance value is invalid, false if valid.
func IsInvalidTolerance(tol float64) bool {
	return tol < 0
}

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
