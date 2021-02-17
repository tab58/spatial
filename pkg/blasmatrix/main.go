package blasmatrix

import (
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas/blas64"
)

// ErrNaN expresses that an argument is NaN.
var ErrNaN = errors.ErrNaN

// ErrInfinity expresses that a number is infinite.
var ErrInfinity = errors.ErrInfinity

// ErrMatrixDims expresses that the matrix dimensions for a specific operation don't match.
var ErrMatrixDims = errors.ErrMatrixDims

// CopyMatrix copies the contents of the src matrix into the dst matrix.
func CopyMatrix(src *blas64.General, dst *blas64.General) error {
	srows, scols, sstride := src.Rows, src.Cols, src.Stride
	drows, dcols, dstride := dst.Rows, dst.Cols, dst.Stride

	if srows != drows || scols != dcols || sstride != dstride {
		return ErrMatrixDims
	}

	k := srows * scols
	for i := 0; i < k; i++ {
		dst.Data[i] = src.Data[i]
	}
	return nil
}

// BlankFromMatrix creates a blank matrix from a given matrix.
func BlankFromMatrix(mat *blas64.General) *blas64.General {
	mrows, mcols := mat.Rows, mat.Cols
	return NewBlas64General(uint(mrows), uint(mcols))
}

// NewBlas64General creates a new blas64.General with the given dimensions.
func NewBlas64General(rows, cols uint) *blas64.General {
	data := make([]float64, 0, rows*cols)
	return &blas64.General{
		Rows:   int(rows),
		Cols:   int(cols),
		Stride: int(rows),
		Data:   data,
	}
}
