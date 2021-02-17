package geometry

import (
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas/blas64"
)

// MatrixReader is a read-only interface for a matrix.
type MatrixReader interface {
	Rows() uint
	Cols() uint
	ElementAt(i, j uint) (float64, error)
	ToBlas64General() *blas64.General
}

// MatrixWriter is a write-only interface for a matrix.
type MatrixWriter interface {
	SetElementAt(i, j uint, value float64) error
}

// Matrix2D is a row-major representation of a 2x2 matrix.
type Matrix2D struct {
	elements [4]float64
}

// Rows returns the number of rows in the matrix.
func (m *Matrix2D) Rows() uint { return 2 }

// Cols returns the number of columns in the matrix.
func (m *Matrix2D) Cols() uint { return 2 }

// ElementAt returns the value of the element at the given indices.
func (m *Matrix2D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, errors.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix2D) ToBlas64General() *blas64.General {
	data := make([]float64, len(m.elements))
	copy(data, m.elements[:])
	return &blas64.General{
		Rows:   int(m.Rows()),
		Cols:   int(m.Cols()),
		Stride: int(m.Cols()),
		Data:   data,
	}
}

// SetElementAt sets the value of the element at the given indices.
func (m *Matrix2D) SetElementAt(i, j uint, value float64) error {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return errors.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}
