package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/errors"
	"github.com/tab58/v1/spatial/pkg/numeric"
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

// Clone returns a deep copy of the matrix.
func (m *Matrix3D) Clone() *Matrix3D {
	a := m.elements
	tmp := [9]float64{a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8]}
	return &Matrix3D{
		elements: tmp,
	}
}

// Copy copies the elements of the matrix to this one.
func (m *Matrix3D) Copy(mat *Matrix3D) {
	a := mat.elements
	m.elements[0] = a[0]
	m.elements[1] = a[1]
	m.elements[2] = a[2]
	m.elements[3] = a[3]
	m.elements[4] = a[4]
	m.elements[5] = a[5]
	m.elements[6] = a[6]
	m.elements[7] = a[7]
	m.elements[8] = a[8]
}

// Identity sets the matrix to the identity matrix.
func (m *Matrix3D) Identity() {
	// ignoring error since all elements will not overflow
	m.SetElements(1, 0, 0, 0, 1, 0, 0, 0, 1)
}

// Scale multiplies the elements of the matrix by the given scalar.
func (m *Matrix3D) Scale(z float64) error {
	out := [9]float64{}
	for i, v := range m.elements {
		val := v * z
		if numeric.IsOverflow(val) {
			return errors.ErrOverflow
		}
		out[i] = val
	}
	m.elements = out
	return nil
}

// ElementAt returns the value of the element at the given indices.
func (m *Matrix3D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, errors.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// SetElements sets the elements in the matrix.
func (m *Matrix3D) SetElements(m00, m01, m02, m10, m11, m12, m20, m21, m22 float64) error {
	if numeric.AreAnyOverflow(m00, m01, m02, m10, m11, m12, m20, m21, m22) {
		return errors.ErrOverflow
	}

	m.elements[0] = m00
	m.elements[1] = m01
	m.elements[2] = m02
	m.elements[3] = m10
	m.elements[4] = m11
	m.elements[5] = m12
	m.elements[6] = m20
	m.elements[7] = m21
	m.elements[8] = m22
	return nil
}

// Elements clones the elements of the matrix and returns them.
func (m *Matrix3D) Elements() [9]float64 {
	tmp := [9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i, v := range m.elements {
		tmp[i] = v
	}
	return tmp
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix3D) ToBlas64General() blas64.General {
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
func (m *Matrix3D) SetElementAt(i, j uint, value float64) error {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return errors.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}

// Add adds the elements of the given matrix to the elements of this matrix.
func (m *Matrix3D) Add(mat *Matrix3D) error {
	tmp := [9]float64{
		m.elements[0] + mat.elements[0],
		m.elements[1] + mat.elements[1],
		m.elements[2] + mat.elements[2],
		m.elements[3] + mat.elements[3],
		m.elements[4] + mat.elements[4],
		m.elements[5] + mat.elements[5],
		m.elements[6] + mat.elements[6],
		m.elements[7] + mat.elements[7],
		m.elements[8] + mat.elements[8],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return errors.ErrOverflow
	}
	m.elements = tmp
	return nil
}

// Sub subtracts the elements of the given matrix from the elements of this matrix.
func (m *Matrix3D) Sub(mat *Matrix3D) error {
	tmp := [9]float64{
		m.elements[0] - mat.elements[0],
		m.elements[1] - mat.elements[1],
		m.elements[2] - mat.elements[2],
		m.elements[3] - mat.elements[3],
		m.elements[4] - mat.elements[4],
		m.elements[5] - mat.elements[5],
		m.elements[6] - mat.elements[6],
		m.elements[7] - mat.elements[7],
		m.elements[8] - mat.elements[8],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return errors.ErrOverflow
	}
	m.elements = tmp
	return nil
}

func multiply3DMatrices(a, b [9]float64) ([9]float64, error) {
	a00, a01, a02 := a[0], a[1], a[2]
	a10, a11, a12 := a[3], a[4], a[5]
	a20, a21, a22 := a[6], a[7], a[8]

	b00, b01, b02 := b[0], b[1], b[2]
	b10, b11, b12 := b[3], b[4], b[5]
	b20, b21, b22 := b[6], b[7], b[8]

	out := [9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}
	out[0] = b00*a00 + b01*a10 + b02*a20
	out[1] = b00*a01 + b01*a11 + b02*a21
	out[2] = b00*a02 + b01*a12 + b02*a22

	out[3] = b10*a00 + b11*a10 + b12*a20
	out[4] = b10*a01 + b11*a11 + b12*a21
	out[5] = b10*a02 + b11*a12 + b12*a22

	out[6] = b20*a00 + b21*a10 + b22*a20
	out[7] = b20*a01 + b21*a11 + b22*a21
	out[8] = b20*a02 + b21*a12 + b22*a22

	if numeric.AreAnyOverflow(out[:]...) {
		return [9]float64{}, errors.ErrOverflow
	}
	return out, nil
}

// Premultiply left-multiplies the given matrix with this one.
func (m *Matrix3D) Premultiply(mat *Matrix3D) error {
	res, err := multiply3DMatrices(mat.elements, m.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Postmultiply right-multiplies the given matrix with this one.
func (m *Matrix3D) Postmultiply(mat *Matrix3D) error {
	res, err := multiply3DMatrices(m.elements, mat.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Invert inverts this matrix in-place.
func (m *Matrix3D) Invert() error {
	a := m.elements
	a00, a01, a02 := a[0], a[1], a[2]
	a10, a11, a12 := a[3], a[4], a[5]
	a20, a21, a22 := a[6], a[7], a[8]

	b01 := a22*a11 - a12*a21
	b11 := -a22*a10 + a12*a20
	b21 := a21*a10 - a11*a20

	// Calculate the determinant
	det := a00*b01 + a01*b11 + a02*b21
	if math.Abs(det) < 1e-13 {
		return errors.ErrSingularMatrix
	}
	det = 1.0 / det

	out := [9]float64{}
	out[0] = b01 * det
	out[1] = (-a22*a01 + a02*a21) * det
	out[2] = (a12*a01 - a02*a11) * det
	out[3] = b11 * det
	out[4] = (a22*a00 - a02*a20) * det
	out[5] = (-a12*a00 + a02*a10) * det
	out[6] = b21 * det
	out[7] = (-a21*a00 + a01*a20) * det
	out[8] = (a11*a00 - a01*a10) * det
	m.elements = out

	return nil
}

// Determinant calculates the determinant of the matrix.
func (m *Matrix3D) Determinant() float64 {
	a := m.elements
	a00, a01, a02 := a[0], a[1], a[2]
	a10, a11, a12 := a[3], a[4], a[5]
	a20, a21, a22 := a[6], a[7], a[8]

	res := a00*(a22*a11-a12*a21) +
		a01*(-a22*a10+a12*a20) +
		a02*(a21*a10-a11*a20)
	return res
}

// Adjoint calculates the adjoint/adjugate matrix.
func (m *Matrix3D) Adjoint() *Matrix3D {
	a := m.elements
	a00, a01, a02 := a[0], a[1], a[2]
	a10, a11, a12 := a[3], a[4], a[5]
	a20, a21, a22 := a[6], a[7], a[8]

	out := [9]float64{}
	out[0] = a11*a22 - a12*a21
	out[1] = a02*a21 - a01*a22
	out[2] = a01*a12 - a02*a11
	out[3] = a12*a20 - a10*a22
	out[4] = a00*a22 - a02*a20
	out[5] = a02*a10 - a00*a12
	out[6] = a10*a21 - a11*a20
	out[7] = a01*a20 - a00*a21
	out[8] = a00*a11 - a01*a10
	return &Matrix3D{
		elements: out,
	}
}

// Transpose transposes the matrix in-place.
func (m *Matrix3D) Transpose() {
	a01 := m.elements[1]
	a02 := m.elements[2]
	a12 := m.elements[5]
	m.elements[1] = m.elements[3]
	m.elements[2] = m.elements[6]
	m.elements[3] = a01
	m.elements[5] = m.elements[7]
	m.elements[6] = a02
	m.elements[7] = a12
}

// IsSingular returns true if the matrix determinant is exactly zero, false if not.
func (m *Matrix3D) IsSingular() bool {
	return m.Determinant() == 0
}

// IsNearSingular returns true if the matrix determinant is equal or below the given tolerance, false if not.
func (m *Matrix3D) IsNearSingular(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	return math.Abs(m.Determinant()) <= tol, nil
}
