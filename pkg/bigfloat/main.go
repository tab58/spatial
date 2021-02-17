package bigfloat

import (
	"math/big"

	"github.com/tab58/v1/spatial/pkg/errors"
)

// ErrNaN expresses that a computation has resulted in an invalid results (NaN).
var ErrNaN = errors.ErrNaN

// ErrOverflow expresses that a computation has resulted in numeric overflow.
var ErrOverflow = errors.ErrOverflow

// ErrUnderflow expresses that a computation has resulted in numeric underflow.
var ErrUnderflow = errors.ErrUnderflow

// Comparator simplifies comparison between float64 and big.Float numbers.
type Comparator struct {
	tmp *big.Float
}

// LT (less than) returns true if a < b, false if not.
func (c *Comparator) LT(a, b *big.Float) bool {
	return a.Cmp(b) < 0
}

// LTFloat64 (less than) returns true if a < b, false if not.
func (c *Comparator) LTFloat64(a *big.Float, b float64) bool {
	c.tmp.SetFloat64(b)
	return c.LT(a, c.tmp)
}

// LTE (less than or equal) returns true if a <= b, false if not.
func (c *Comparator) LTE(a, b *big.Float) bool {
	return a.Cmp(b) <= 0
}

// LTEFloat64 (less than or equal) returns true if a <= b, false if not.
func (c *Comparator) LTEFloat64(a *big.Float, b float64) bool {
	c.tmp.SetFloat64(b)
	return c.LTE(a, c.tmp)
}

// GT (greater than) returns true if a > b, false if not.
func (c *Comparator) GT(a, b *big.Float) bool {
	return a.Cmp(b) > 0
}

// GTFloat64 (greater than) returns true if a > b, false if not.
func (c *Comparator) GTFloat64(a *big.Float, b float64) bool {
	c.tmp.SetFloat64(b)
	return c.GT(a, c.tmp)
}

// GTE (greater than or equal) returns true if a >= b, false if not.
func (c *Comparator) GTE(a, b *big.Float) bool {
	return a.Cmp(b) >= 0
}

// GTEFloat64 (greater than or equal) returns true if a >= b, false if not.
func (c *Comparator) GTEFloat64(a *big.Float, b float64) bool {
	c.tmp.SetFloat64(b)
	return c.GTE(a, c.tmp)
}

// NewComparator creates a new Comparator.
func NewComparator() *Comparator {
	return &Comparator{
		tmp: big.NewFloat(0),
	}
}

// SortMaxMin sorts the values, maximum value first and minimum value second. Returns the values in argument order if the same.
func SortMaxMin(a, b *big.Float) (max *big.Float, min *big.Float) {
	if a.Cmp(b) < 0 {
		return b, a
	}
	return a, b
}

// Nrm2 computes the 2-norm of the vector for which the components are given as arguments.
func Nrm2(a, b *big.Float) *big.Float {
	valA, _ := a.Float64()
	valB, _ := b.Float64()

	if valA == 0 && valB == 0 {
		return big.NewFloat(0)
	}

	x, y := big.NewFloat(0), big.NewFloat(0)
	x.Copy(a)
	y.Copy(b)

	x.Abs(x)
	y.Abs(y)
	u, t := SortMaxMin(x, y)
	t.Quo(t, u)

	res := big.NewFloat(0)
	res.Copy(t)
	res.Mul(res, t)
	res.Add(res, big.NewFloat(1))
	res.Sqrt(res)
	res.Mul(res, u)

	return res
}
