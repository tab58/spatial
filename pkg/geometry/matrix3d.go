package geometry

import (
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas/blas64"
)

// Matrix3D is a row-major representation of a 3x3 matrix.
type Matrix3D struct {
	elements [9]float64
}

// Rows returns the number of rows in the matrix.
func (m *Matrix3D) Rows() uint { return 3 }

// Cols returns the number of columns in the matrix.
func (m *Matrix3D) Cols() uint { return 3 }

// ElementAt returns the value of the element at the given indices.
func (m *Matrix3D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, errors.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix3D) ToBlas64General() *blas64.General {
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
func (m *Matrix3D) SetElementAt(i, j uint, value float64) error {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return errors.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}
