package bigfloat

import (
	"errors"
	"math"
	"math/big"
)

// ErrOverflow expresses that a computation has resulted in numeric overflow.
var ErrOverflow = errors.New("numeric overflow")

// ErrUnderflow expresses that a computation has resulted in numeric underflow.
var ErrUnderflow = errors.New("numeric underflow")

// ErrNaN expresses that a computation has resulted in an invalid results (NaN).
var ErrNaN = errors.New("result is NaN")

func hasUnderflowed(f float64, acc big.Accuracy) bool {
	return (f == 0 && acc == big.Below) || (f == -0 && acc == big.Above)
}

func hasOverflowed(f float64, acc big.Accuracy) bool {
	return math.IsInf(f, 0)
}

// HasNumericErr returns errors if computations have resulted in overflow, underflow, or NaN.
func HasNumericErr(result float64, accuracy big.Accuracy) error {
	if hasOverflowed(result, accuracy) {
		return ErrOverflow
	} else if hasUnderflowed(result, accuracy) {
		return ErrUnderflow
	} else if math.IsNaN(result) {
		return ErrNaN
	} else {
		return nil
	}
}

// Calculator is a data structure meant to be a running accumulator of arithmetic operations.
type Calculator struct {
	result *big.Float
	tmp    *big.Float
}

// NewCalculator creates a extended precision calculator.
func NewCalculator(z float64) *Calculator {
	return &Calculator{
		result: big.NewFloat(z),
		tmp:    big.NewFloat(0),
	}
}

func (f *Calculator) arithBinaryOp(a float64, op func(*big.Float, *big.Float)) *Calculator {
	f.tmp.SetFloat64(a)
	op(f.result, f.tmp)
	return f
}

func (f *Calculator) arithUnaryOp(op func(*big.Float)) *Calculator {
	op(f.result)
	return f
}

// Value gets a clone of the current value of the result.
func (f *Calculator) Value() *big.Float {
	res := big.NewFloat(0)
	res.Copy(f.result)
	return res
}

// SetFloat64 sets the value of the result to z.
func (f *Calculator) SetFloat64(z float64) *Calculator {
	f.result.SetFloat64(z)
	return f
}

// AccuracyFloat64 returns the result converted to a float64 with an indication of its accuracy.
func (f *Calculator) AccuracyFloat64() (float64, big.Accuracy) {
	return f.result.Float64()
}

// Float64 converts the result into a float64 and indicates if there is an error.
func (f *Calculator) Float64() (float64, error) {
	res, acc := f.result.Float64()
	return res, HasNumericErr(res, acc)
}

// Add adds a to the result.
func (f *Calculator) Add(a float64) *Calculator {
	return f.arithBinaryOp(a, addOp)
}

// Sub subtracts a to the result.
func (f *Calculator) Sub(a float64) *Calculator {
	return f.arithBinaryOp(a, subOp)
}

// Mul multiplies a to the result.
func (f *Calculator) Mul(a float64) *Calculator {
	return f.arithBinaryOp(a, mulOp)
}

// Quo divides a from the result.
func (f *Calculator) Quo(a float64) *Calculator {
	return f.arithBinaryOp(a, quoOp)
}

// Neg negates -a and returns the result.
func (f *Calculator) Neg() *Calculator {
	return f.arithUnaryOp(negOp)
}

// Sqrt computes the square root of the result.
func (f *Calculator) Sqrt(a float64) *Calculator {
	return f.arithUnaryOp(sqrtOp)
}

// Abs computes the absolute value of the result.
func (f *Calculator) Abs(a float64) *Calculator {
	return f.arithUnaryOp(absOp)
}

func addOp(res, tmp *big.Float) {
	res.Add(res, tmp)
}

func subOp(res, tmp *big.Float) {
	res.Sub(res, tmp)
}

func mulOp(res, tmp *big.Float) {
	res.Mul(res, tmp)
}

func quoOp(res, tmp *big.Float) {
	res.Quo(res, tmp)
}

func negOp(res *big.Float) {
	res.Neg(res)
}

func sqrtOp(res *big.Float) {
	res.Sqrt(res)
}

func absOp(res *big.Float) {
	res.Abs(res)
}
