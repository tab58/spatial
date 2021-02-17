package geometry

import (
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas/blas64"
)

// Matrix4D is a row-major representation of a 4x4 matrix.
type Matrix4D struct {
	elements [16]float64
}

// Rows returns the number of rows in the matrix.
func (m *Matrix4D) Rows() uint { return 4 }

// Cols returns the number of columns in the matrix.
func (m *Matrix4D) Cols() uint { return 4 }

// ElementAt returns the value of the element at the given indices.
func (m *Matrix4D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, errors.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix4D) ToBlas64General() *blas64.General {
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
func (m *Matrix4D) SetElementAt(i, j uint, value float64) error {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return errors.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}
