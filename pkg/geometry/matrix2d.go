package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/numeric"
	"gonum.org/v1/gonum/blas/blas64"
)

// // MatrixReader is a read-only interface for a matrix.
// type MatrixReader interface {
// 	Rows() uint
// 	Cols() uint
// 	ElementAt(i, j uint) (float64, error)
// 	ToBlas64General() blas64.General
// }

// // MatrixWriter is a write-only interface for a matrix.
// type MatrixWriter interface {
// 	SetElementAt(i, j uint, value float64) error
// }

// Matrix2D is a row-major representation of a 2x2 matrix.
type Matrix2D struct {
	elements [4]float64
}

// Rows returns the number of rows in the matrix.
func (m *Matrix2D) Rows() uint { return 2 }

// Cols returns the number of columns in the matrix.
func (m *Matrix2D) Cols() uint { return 2 }

// Clone returns a deep copy of the matrix.
func (m *Matrix2D) Clone() *Matrix2D {
	a := m.elements
	tmp := [4]float64{a[0], a[1], a[2], a[3]}
	return &Matrix2D{
		elements: tmp,
	}
}

// Copy copies the elements of the matrix to this one.
func (m *Matrix2D) Copy(mat *Matrix2D) {
	a := mat.elements
	m.elements[0] = a[0]
	m.elements[1] = a[1]
	m.elements[2] = a[2]
	m.elements[3] = a[3]
}

// Identity sets the matrix to the identity matrix.
func (m *Matrix2D) Identity() {
	// ignoring error since all elements will not overflow
	m.SetElements(1, 0, 0, 1)
}

// Scale multiplies the elements of the matrix by the given scalar.
func (m *Matrix2D) Scale(z float64) error {
	out := [4]float64{}
	for i, v := range m.elements {
		val := v * z
		if numeric.IsOverflow(val) {
			return numeric.ErrOverflow
		}
		out[i] = val
	}
	m.elements = out
	return nil
}

// ElementAt returns the value of the element at the given indices.
func (m *Matrix2D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, numeric.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix2D) ToBlas64General() blas64.General {
	data := make([]float64, len(m.elements))
	copy(data, m.elements[:])
	return blas64.General{
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
		return numeric.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}

// SetElements sets the elements in the matrix.
func (m *Matrix2D) SetElements(m00, m01, m10, m11 float64) error {
	if numeric.AreAnyOverflow(m00, m01, m10, m11) {
		return numeric.ErrOverflow
	}

	m.elements[0] = m00
	m.elements[1] = m01
	m.elements[2] = m10
	m.elements[3] = m11
	return nil
}

// Elements clones the elements of the matrix and returns them.
func (m *Matrix2D) Elements() [4]float64 {
	tmp := [4]float64{}
	for i, v := range m.elements {
		tmp[i] = v
	}
	return tmp
}

// Add adds the elements of the given matrix to the elements of this matrix.
func (m *Matrix2D) Add(mat *Matrix2D) error {
	tmp := [4]float64{
		m.elements[0] + mat.elements[0],
		m.elements[1] + mat.elements[1],
		m.elements[2] + mat.elements[2],
		m.elements[3] + mat.elements[3],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return numeric.ErrOverflow
	}
	m.elements = tmp
	return nil
}

// Sub subtracts the elements of the given matrix to the elements of this matrix.
func (m *Matrix2D) Sub(mat *Matrix2D) error {
	tmp := [4]float64{
		m.elements[0] - mat.elements[0],
		m.elements[1] - mat.elements[1],
		m.elements[2] - mat.elements[2],
		m.elements[3] - mat.elements[3],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return numeric.ErrOverflow
	}
	m.elements = tmp
	return nil
}

func multiply2DMatrices(a, b [4]float64) ([4]float64, error) {
	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]

	b0 := b[0]
	b1 := b[1]
	b2 := b[2]
	b3 := b[3]

	out := [4]float64{}
	out[0] = a0*b0 + a2*b1
	out[1] = a1*b0 + a3*b1
	out[2] = a0*b2 + a2*b3
	out[3] = a1*b2 + a3*b3
	return out, nil
}

// Premultiply left-multiplies the given matrix with this one.
func (m *Matrix2D) Premultiply(mat *Matrix2D) error {
	res, err := multiply2DMatrices(mat.elements, m.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Postmultiply right-multiplies the given matrix with this one.
func (m *Matrix2D) Postmultiply(mat *Matrix2D) error {
	res, err := multiply2DMatrices(m.elements, mat.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Invert inverts this matrix in-place.
func (m *Matrix2D) Invert() error {
	a := m.elements
	a0, a1, a2, a3 := a[0], a[1], a[2], a[3]

	// Calculate the determinant
	det := a0*a3 - a2*a1
	if math.Abs(det) < 1e-13 {
		return numeric.ErrSingularMatrix
	}
	det = 1.0 / det

	out := [4]float64{}
	out[0] = a3 * det
	out[1] = -a1 * det
	out[2] = -a2 * det
	out[3] = a0 * det
	m.elements = out

	return nil
}

// Determinant calculates the determinant of the matrix.
func (m *Matrix2D) Determinant() float64 {
	a := m.elements
	return a[0]*a[3] - a[2]*a[1]
}

// Adjoint calculates the adjoint/adjugate matrix.
func (m *Matrix2D) Adjoint() *Matrix2D {
	a := m.elements
	// Caching this value is nessecary if out == a
	a0 := a[0]

	out := [4]float64{}
	out[0] = a[3]
	out[1] = -a[1]
	out[2] = -a[2]
	out[3] = a0
	return &Matrix2D{
		elements: out,
	}
}

// Transpose transposes the matrix in-place.
func (m *Matrix2D) Transpose() {
	a1 := m.elements[1]
	m.elements[1] = m.elements[2]
	m.elements[2] = a1
}

// IsSingular returns true if the matrix determinant is exactly zero, false if not.
func (m *Matrix2D) IsSingular() bool {
	return m.Determinant() == 0
}

// IsNearSingular returns true if the matrix determinant is equal or below the given tolerance, false if not.
func (m *Matrix2D) IsNearSingular(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	return math.Abs(m.Determinant()) <= tol, nil
}
