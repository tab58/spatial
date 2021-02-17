package blasmatrix

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas/blas64"
)

// Vector3DReader is a read-only interface for a Vector3D.
type Vector3DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64
	Length() (float64, error)
	ToBlasVector() blas64.Vector
}

// BuildMatrix3DSkewSymmetric builds the skew symmetric matrix based on vector component values.
func BuildMatrix3DSkewSymmetric(v Vector3DReader) *blas64.General {
	// TODO: do check on v's Length, etc.

	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	return &blas64.General{
		Rows:   3,
		Cols:   3,
		Data:   []float64{0, -z, y, z, 0, -x, -y, x, 0},
		Stride: 3,
	}
}

// Get3DRotMatrix returns a rotation matrix that rotates about with the specified angle.
func Get3DRotMatrix(axis Vector3DReader, angle float64) (*blas64.General, error) {
	UU := NewBlas64General(3, 3)
	l, err := axis.Length()
	if err != nil {
		return UU, err
	}

	// TODO: remove this magic number
	lengthTol := 1e-14
	if math.Abs(l) < lengthTol {
		return UU, errors.ErrVectorZeroLength
	}

	if math.IsNaN(angle) {
		return UU, ErrNaN
	}
	if math.IsInf(angle, 0) {
		return UU, ErrInfinity
	}

	c := math.Cos(angle)
	s := math.Sin(angle)
	c1 := 1 - c

	u := axis.ToBlasVector()
	blas64.Ger(c1, u, u, *UU)
	// UU = (1 - cos(t)) * outerProduct(u, u)

	// cI = cos(t) * I
	cI := &blas64.General{
		Rows:   3,
		Cols:   3,
		Data:   []float64{c, 0, 0, 0, c, 0, 0, 0, c},
		Stride: 3,
	}

	Ux := BuildMatrix3DSkewSymmetric(axis)
	calc := NewCalculator(3, 3)
	// Ux = sin(t) * skewSymm(axis)
	calc.Add(Ux)
	calc.Scale(s)

	calc.Add(cI)
	calc.Add(UU)

	return calc.Value(), nil
}
