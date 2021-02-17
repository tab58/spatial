package blasmatrix

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Calculator aggregates operations on a blas64.General matrix.
type Calculator struct {
	result *blas64.General
}

// Value returns a copy of the result matrix.
func (c *Calculator) Value() *blas64.General {
	res := BlankFromMatrix(c.result)
	CopyMatrix(c.result, res)
	return res
}

// Add adds the matrix to the current result.
func (c *Calculator) Add(mat *blas64.General) error {
	mrows, mcols := mat.Rows, mat.Cols
	rows, cols := c.result.Rows, c.result.Cols
	if mrows != rows || mcols != cols {
		return ErrMatrixDims
	}

	k := mrows * mcols
	for i := 0; i < k; i++ {
		x := c.result.Data[i]
		y := mat.Data[i]
		c.result.Data[i] = x + y
	}
	return nil
}

// Sub subtracts the matrix from the current result.
func (c *Calculator) Sub(mat *blas64.General) error {
	mrows, mcols := mat.Rows, mat.Cols
	rows, cols := c.result.Rows, c.result.Cols
	if mrows != rows || mcols != cols {
		return ErrMatrixDims
	}

	k := mrows * mcols
	for i := 0; i < k; i++ {
		x := c.result.Data[i]
		y := mat.Data[i]
		c.result.Data[i] = x - y
	}

	return nil
}

// Scale multiplies the matrix by the given factor.
func (c *Calculator) Scale(z float64) {
	rows, cols := c.result.Rows, c.result.Cols
	k := rows * cols
	for i := 0; i < k; i++ {
		x := c.result.Data[i]
		c.result.Data[i] = x * z
	}
}

// Premultiply computes the expression mat * result.
func (c *Calculator) Premultiply(mat *blas64.General) error {
	mcols := mat.Cols
	rows := c.result.Rows
	if mcols != rows {
		return ErrMatrixDims
	}

	tmp := *BlankFromMatrix(c.result)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, *mat, *c.result, 0, tmp)
	err := CopyMatrix(&tmp, c.result)
	if err != nil {
		return err
	}
	return nil
}

// Postmultiply computes the expression result * mat.
func (c *Calculator) Postmultiply(mat *blas64.General) error {
	mrows := mat.Rows
	cols := c.result.Cols
	if mrows != cols {
		return ErrMatrixDims
	}

	tmp := *BlankFromMatrix(c.result)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, *c.result, *mat, 0, tmp)
	err := CopyMatrix(&tmp, c.result)
	if err != nil {
		return err
	}
	return nil
}

// NewCalculator returns a new Calculator.
func NewCalculator(rows, cols uint) *Calculator {
	return &Calculator{
		result: NewBlas64General(rows, cols),
	}
}
