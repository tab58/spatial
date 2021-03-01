package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/numeric"
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

// Clone returns a deep copy of the matrix.
func (m *Matrix4D) Clone() *Matrix4D {
	a := m.elements
	tmp := [16]float64{a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14], a[15]}
	return &Matrix4D{
		elements: tmp,
	}
}

// Copy copies the elements of the matrix to this one.
func (m *Matrix4D) Copy(mat *Matrix4D) {
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
	m.elements[9] = a[9]
	m.elements[10] = a[10]
	m.elements[11] = a[11]
	m.elements[12] = a[12]
	m.elements[13] = a[13]
	m.elements[14] = a[14]
	m.elements[15] = a[15]
}

// Identity sets the matrix to the identity matrix.
func (m *Matrix4D) Identity() {
	// ignoring error since all elements will not overflow
	m.SetElements(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)
}

// Scale multiplies the elements of the matrix by the given scalar.
func (m *Matrix4D) Scale(z float64) error {
	out := [16]float64{}
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
func (m *Matrix4D) ElementAt(i, j uint) (float64, error) {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return 0, numeric.ErrMatrixOutOfRange
	}
	return m.elements[i*cols+j], nil
}

// ToBlas64General returns a blas64.General with the same values as the matrix.
func (m *Matrix4D) ToBlas64General() blas64.General {
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
func (m *Matrix4D) SetElementAt(i, j uint, value float64) error {
	cols := m.Cols()
	if i <= m.Rows() || j <= cols {
		return numeric.ErrMatrixOutOfRange
	}
	m.elements[i*cols+j] = value
	return nil
}

// SetElements sets the elements in the matrix.
func (m *Matrix4D) SetElements(m00, m01, m02, m03, m10, m11, m12, m13, m20, m21, m22, m23, m30, m31, m32, m33 float64) error {
	if numeric.AreAnyOverflow(m00, m01, m02, m03, m10, m11, m12, m13, m20, m21, m22, m23, m30, m31, m32, m33) {
		return numeric.ErrOverflow
	}

	m.elements[0] = m00
	m.elements[1] = m01
	m.elements[2] = m02
	m.elements[3] = m03

	m.elements[4] = m10
	m.elements[5] = m11
	m.elements[6] = m12
	m.elements[7] = m13

	m.elements[8] = m20
	m.elements[9] = m21
	m.elements[10] = m22
	m.elements[11] = m23

	m.elements[12] = m30
	m.elements[13] = m31
	m.elements[14] = m32
	m.elements[15] = m33
	return nil
}

// Elements clones the elements of the matrix and returns them.
func (m *Matrix4D) Elements() [16]float64 {
	tmp := [16]float64{}
	for i, v := range m.elements {
		tmp[i] = v
	}
	return tmp
}

// Add adds the elements of the given matrix to the elements of this matrix.
func (m *Matrix4D) Add(mat *Matrix4D) error {
	tmp := [16]float64{
		m.elements[0] + mat.elements[0],
		m.elements[1] + mat.elements[1],
		m.elements[2] + mat.elements[2],
		m.elements[3] + mat.elements[3],
		m.elements[4] + mat.elements[4],
		m.elements[5] + mat.elements[5],
		m.elements[6] + mat.elements[6],
		m.elements[7] + mat.elements[7],
		m.elements[8] + mat.elements[8],
		m.elements[9] + mat.elements[9],
		m.elements[10] + mat.elements[10],
		m.elements[11] + mat.elements[11],
		m.elements[12] + mat.elements[12],
		m.elements[13] + mat.elements[13],
		m.elements[14] + mat.elements[14],
		m.elements[15] + mat.elements[15],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return numeric.ErrOverflow
	}
	m.elements = tmp
	return nil
}

// Sub subtracts the elements of the given matrix to the elements of this matrix.
func (m *Matrix4D) Sub(mat *Matrix4D) error {
	tmp := [16]float64{
		m.elements[0] - mat.elements[0],
		m.elements[1] - mat.elements[1],
		m.elements[2] - mat.elements[2],
		m.elements[3] - mat.elements[3],
		m.elements[4] - mat.elements[4],
		m.elements[5] - mat.elements[5],
		m.elements[6] - mat.elements[6],
		m.elements[7] - mat.elements[7],
		m.elements[8] - mat.elements[8],
		m.elements[9] - mat.elements[9],
		m.elements[10] - mat.elements[10],
		m.elements[11] - mat.elements[11],
		m.elements[12] - mat.elements[12],
		m.elements[13] - mat.elements[13],
		m.elements[14] - mat.elements[14],
		m.elements[15] - mat.elements[15],
	}

	if numeric.AreAnyOverflow(tmp[:]...) {
		return numeric.ErrOverflow
	}
	m.elements = tmp
	return nil
}

func multiply4DMatrices(a, b [16]float64) ([16]float64, error) {
	a00, a01, a02, a03 := a[0], a[1], a[2], a[3]
	a10, a11, a12, a13 := a[4], a[5], a[6], a[7]
	a20, a21, a22, a23 := a[8], a[9], a[10], a[11]
	a30, a31, a32, a33 := a[12], a[13], a[14], a[15]

	// Cache only the current line of the second matrix
	b0, b1, b2, b3 := b[0], b[1], b[2], b[3]

	out := [16]float64{}
	out[0] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	out[1] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	out[2] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	out[3] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0, b1, b2, b3 = b[4], b[5], b[6], b[7]
	out[4] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	out[5] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	out[6] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	out[7] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0, b1, b2, b3 = b[8], b[9], b[10], b[11]
	out[8] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	out[9] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	out[10] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	out[11] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0, b1, b2, b3 = b[12], b[13], b[14], b[15]
	out[12] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	out[13] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	out[14] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	out[15] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	return out, nil
}

// Premultiply left-multiplies the given matrix with this one.
func (m *Matrix4D) Premultiply(mat *Matrix4D) error {
	res, err := multiply4DMatrices(mat.elements, m.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Postmultiply right-multiplies the given matrix with this one.
func (m *Matrix4D) Postmultiply(mat *Matrix4D) error {
	res, err := multiply4DMatrices(m.elements, mat.elements)
	if err != nil {
		return err
	}
	m.elements = res
	return nil
}

// Invert inverts this matrix in-place.
func (m *Matrix4D) Invert() error {
	a := m.elements
	a00, a01, a02, a03 := a[0], a[1], a[2], a[3]
	a10, a11, a12, a13 := a[4], a[5], a[6], a[7]
	a20, a21, a22, a23 := a[8], a[9], a[10], a[11]
	a30, a31, a32, a33 := a[12], a[13], a[14], a[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	// Calculate the determinant
	det := b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
	if math.Abs(det) < 1e-13 {
		return numeric.ErrSingularMatrix
	}
	det = 1.0 / det

	out := [16]float64{}
	out[0] = (a11*b11 - a12*b10 + a13*b09) * det
	out[1] = (a02*b10 - a01*b11 - a03*b09) * det
	out[2] = (a31*b05 - a32*b04 + a33*b03) * det
	out[3] = (a22*b04 - a21*b05 - a23*b03) * det
	out[4] = (a12*b08 - a10*b11 - a13*b07) * det
	out[5] = (a00*b11 - a02*b08 + a03*b07) * det
	out[6] = (a32*b02 - a30*b05 - a33*b01) * det
	out[7] = (a20*b05 - a22*b02 + a23*b01) * det
	out[8] = (a10*b10 - a11*b08 + a13*b06) * det
	out[9] = (a01*b08 - a00*b10 - a03*b06) * det
	out[10] = (a30*b04 - a31*b02 + a33*b00) * det
	out[11] = (a21*b02 - a20*b04 - a23*b00) * det
	out[12] = (a11*b07 - a10*b09 - a12*b06) * det
	out[13] = (a00*b09 - a01*b07 + a02*b06) * det
	out[14] = (a31*b01 - a30*b03 - a32*b00) * det
	out[15] = (a20*b03 - a21*b01 + a22*b00) * det
	m.elements = a
	return nil
}

// Determinant calculates the determinant of the matrix.
func (m *Matrix4D) Determinant() float64 {
	a := m.elements
	a00, a01, a02, a03 := a[0], a[1], a[2], a[3]
	a10, a11, a12, a13 := a[4], a[5], a[6], a[7]
	a20, a21, a22, a23 := a[8], a[9], a[10], a[11]
	a30, a31, a32, a33 := a[12], a[13], a[14], a[15]

	b0 := a00*a11 - a01*a10
	b1 := a00*a12 - a02*a10
	b2 := a01*a12 - a02*a11
	b3 := a20*a31 - a21*a30
	b4 := a20*a32 - a22*a30
	b5 := a21*a32 - a22*a31
	b6 := a00*b5 - a01*b4 + a02*b3
	b7 := a10*b5 - a11*b4 + a12*b3
	b8 := a20*b2 - a21*b1 + a22*b0
	b9 := a30*b2 - a31*b1 + a32*b0

	return a13*b6 - a03*b7 + a33*b8 - a23*b9
}

// Adjoint calculates the adjoint/adjugate matrix.
func (m *Matrix4D) Adjoint() *Matrix4D {
	a := m.elements
	a00, a01, a02, a03 := a[0], a[1], a[2], a[3]
	a10, a11, a12, a13 := a[4], a[5], a[6], a[7]
	a20, a21, a22, a23 := a[8], a[9], a[10], a[11]
	a30, a31, a32, a33 := a[12], a[13], a[14], a[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	out := [16]float64{}
	out[0] = a11*b11 - a12*b10 + a13*b09
	out[1] = a02*b10 - a01*b11 - a03*b09
	out[2] = a31*b05 - a32*b04 + a33*b03
	out[3] = a22*b04 - a21*b05 - a23*b03
	out[4] = a12*b08 - a10*b11 - a13*b07
	out[5] = a00*b11 - a02*b08 + a03*b07
	out[6] = a32*b02 - a30*b05 - a33*b01
	out[7] = a20*b05 - a22*b02 + a23*b01
	out[8] = a10*b10 - a11*b08 + a13*b06
	out[9] = a01*b08 - a00*b10 - a03*b06
	out[10] = a30*b04 - a31*b02 + a33*b00
	out[11] = a21*b02 - a20*b04 - a23*b00
	out[12] = a11*b07 - a10*b09 - a12*b06
	out[13] = a00*b09 - a01*b07 + a02*b06
	out[14] = a31*b01 - a30*b03 - a32*b00
	out[15] = a20*b03 - a21*b01 + a22*b00
	return &Matrix4D{
		elements: out,
	}
}

// Transpose transposes the matrix in-place.
func (m *Matrix4D) Transpose() {
	a01 := m.elements[1]
	a02 := m.elements[2]
	a03 := m.elements[3]
	a12 := m.elements[6]
	a13 := m.elements[7]
	a23 := m.elements[11]

	m.elements[1] = m.elements[4]
	m.elements[2] = m.elements[8]
	m.elements[3] = m.elements[12]
	m.elements[4] = a01
	m.elements[6] = m.elements[9]
	m.elements[7] = m.elements[13]
	m.elements[8] = a02
	m.elements[9] = a12
	m.elements[11] = m.elements[14]
	m.elements[12] = a03
	m.elements[13] = a13
	m.elements[14] = a23
}

// IsSingular returns true if the matrix determinant is exactly zero, false if not.
func (m *Matrix4D) IsSingular() bool {
	return m.Determinant() == 0
}

// IsNearSingular returns true if the matrix determinant is equal or below the given tolerance, false if not.
func (m *Matrix4D) IsNearSingular(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	return math.Abs(m.Determinant()) <= tol, nil
}
