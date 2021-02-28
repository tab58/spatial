package geometry

import "gonum.org/v1/gonum/blas/blas64"

// Vector4DReader is a read-only interface for a 4D vector.
type Vector4DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64
	GetW() float64

	Length() (float64, error)
	LengthSquared() (float64, error)
	Clone() *Vector4D
	ToBlasVector() blas64.Vector
	GetNormalizedVector() *Vector4D

	IsZeroLength(tol float64) (bool, error)
	IsUnitLength(tol float64) (bool, error)

	AngleTo(w Vector4DReader) (float64, error)
	Dot(w Vector4DReader) (float64, error)

	IsPerpendicularTo(w Vector4DReader, tol float64) (bool, error)
	IsCodirectionalTo(w Vector4DReader, tol float64) (bool, error)
	IsParallelTo(w Vector4DReader, tol float64) (bool, error)
	IsEqualTo(w Vector4DReader, tol float64) (bool, error)

	MatrixTransform4D(m *Matrix4D) error
}

// Vector4DWriter is a write-only interface for a 4D vector.
type Vector4DWriter interface {
	SetX(float64)
	SetY(float64)
	SetZ(float64)
	SetW(float64)

	Negate()
	Add(w Vector4DReader) error
	Sub(w Vector4DReader) error
	Normalize() error
	Scale(f float64) error
}

// XAxis4D represents the canonical Cartesian x-axis in 3 dimensions.
var XAxis4D Vector4DReader = &Vector4D{X: 1, Y: 0, Z: 0, W: 0}

// YAxis4D represents the canonical Cartesian y-axis in 3 dimensions.
var YAxis4D Vector4DReader = &Vector4D{X: 0, Y: 1, Z: 0, W: 0}

// ZAxis4D represents the canonical Cartesian z-axis in 3 dimensions.
var ZAxis4D Vector4DReader = &Vector4D{X: 0, Y: 1, Z: 1, W: 0}

// WAxis4D represents the canonical Cartesian z-axis in 4 dimensions.
var WAxis4D Vector4DReader = &Vector4D{X: 0, Y: 1, Z: 1, W: 0}

// Zero4D represents the zero vector in the 3D plane.
var Zero4D Vector4DReader = &Vector4D{X: 0, Y: 0, Z: 0, W: 0}

// Vector4D is a representation of a vector in 4 dimensions.
type Vector4D struct {
	X float64
	Y float64
	Z float64
	W float64
}

func (v *Vector4D) GetX() float64 { return v.X }

func (v *Vector4D) GetY() float64 { return v.Y }

func (v *Vector4D) GetZ() float64 { return v.Z }

func (v *Vector4D) GetW() float64 { return v.W }

func (v *Vector4D) Length() (float64, error) {

}

func (v *Vector4D) LengthSquared() (float64, error) {

}

func (v *Vector4D) Clone() *Vector4D {

}

func (v *Vector4D) ToBlasVector() blas64.Vector {

}

func (v *Vector4D) GetNormalizedVector() *Vector4D {

}

func (v *Vector4D) IsZeroLength(tol float64) (bool, error) {

}

func (v *Vector4D) IsUnitLength(tol float64) (bool, error) {

}

func (v *Vector4D) AngleTo(w Vector4DReader) (float64, error) {

}

func (v *Vector4D) Dot(w Vector4DReader) (float64, error) {

}

func (v *Vector4D) IsPerpendicularTo(w Vector4DReader, tol float64) (bool, error) {

}

func (v *Vector4D) IsCodirectionalTo(w Vector4DReader, tol float64) (bool, error) {

}

func (v *Vector4D) IsParallelTo(w Vector4DReader, tol float64) (bool, error) {

}

func (v *Vector4D) IsEqualTo(w Vector4DReader, tol float64) (bool, error) {

}

func (v *Vector4D) MatrixTransform4D(m *Matrix4D) error {

}

func (v *Vector4D) SetX(z float64) { v.X = z }

func (v *Vector4D) SetY(z float64) { v.Y = z }

func (v *Vector4D) SetZ(z float64) { v.Z = z }

func (v *Vector4D) SetW(z float64) { v.W = z }

func (v *Vector4D) Negate() {
	x, y, z, w := v.GetX(), v.GetY(), v.GetZ(), v.GetW()
	v.SetX(-x)
	v.SetY(-y)
	v.SetZ(-z)
	v.SetW(-w)
}

func (v *Vector4D) Add(w Vector4DReader) error {

}

func (v *Vector4D) Sub(w Vector4DReader) error {

}

func (v *Vector4D) Normalize() error {

}

func (v *Vector4D) Scale(f float64) error {

}
