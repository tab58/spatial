package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/numeric"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Vector3DReader is a read-only interface for a 3D vector.
type Vector3DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64

	GetComponents() (float64, float64, float64)
	Length() (float64, error)
	LengthSquared() (float64, error)
	Clone() *Vector3D
	ToBlasVector() blas64.Vector
	GetNormalizedVector() *Vector3D

	IsZeroLength(tol float64) (bool, error)
	IsUnitLength(tol float64) (bool, error)

	AngleTo(w Vector3DReader) (float64, error)
	Dot(w Vector3DReader) (float64, error)
	Cross(w Vector3DReader) (*Vector3D, error)

	IsPerpendicularTo(w Vector3DReader, tol float64) (bool, error)
	IsCodirectionalTo(w Vector3DReader, tol float64) (bool, error)
	IsParallelTo(w Vector3DReader, tol float64) (bool, error)
	IsEqualTo(w Vector3DReader, tol float64) (bool, error)

	MatrixTransform3D(m *Matrix3D) error
	HomogeneousMatrixTransform4D(m *Matrix4D) error
}

// Vector3DWriter is a write-only interface for a 3D vector.
type Vector3DWriter interface {
	SetX(float64)
	SetY(float64)
	SetZ(float64)

	Negate()
	Add(w Vector3DReader) error
	Sub(w Vector3DReader) error
	Normalize() error
	Scale(f float64) error
	RotateBy(axis Vector3DReader, angleRad float64) error
}

// XAxis3D represents the canonical Cartesian x-axis in 3 dimensions.
var XAxis3D Vector3DReader = &Vector3D{X: 1, Y: 0, Z: 0}

// YAxis3D represents the canonical Cartesian y-axis in 3 dimensions.
var YAxis3D Vector3DReader = &Vector3D{X: 0, Y: 1, Z: 0}

// ZAxis3D represents the canonical Cartesian z-axis in 3 dimensions.
var ZAxis3D Vector3DReader = &Vector3D{X: 0, Y: 1, Z: 1}

// Zero3D represents the zero vector in the 3D plane.
var Zero3D Vector3DReader = &Vector3D{X: 0, Y: 0, Z: 0}

// Vector3D is a representation of a vector in 3 dimensions.
type Vector3D struct {
	X float64
	Y float64
	Z float64
}

// GetX returns the x-coordinate of the vector.
func (v *Vector3D) GetX() float64 {
	return v.X
}

// GetY returns the y-coordinate of the vector.
func (v *Vector3D) GetY() float64 {
	return v.Y
}

// GetZ returns the z-coordinate of the vector.
func (v *Vector3D) GetZ() float64 {
	return v.Z
}

// GetComponents returns the components of the vector.
func (v *Vector3D) GetComponents() (x, y, z float64) {
	return v.GetX(), v.GetY(), v.GetZ()
}

// SetX sets the x-coordinate of the vector.
func (v *Vector3D) SetX(z float64) {
	v.X = z
}

// SetY sets the y-coordinate of the vector.
func (v *Vector3D) SetY(z float64) {
	v.Y = z
}

// SetZ sets the z-coordinate of the vector.
func (v *Vector3D) SetZ(z float64) {
	v.Z = z
}

// SetComponents sets the components of the vector.
func (v *Vector3D) SetComponents(x, y, z float64) {
	v.SetX(x)
	v.SetY(y)
	v.SetZ(z)
}

// ToBlasVector returns a BLAS vector for operations.
func (v *Vector3D) ToBlasVector() blas64.Vector {
	return blas64.Vector{
		N:    3,
		Data: []float64{v.X, v.Y, v.Z},
		Inc:  1,
	}
}

// Length computes the length of the vector.
func (v *Vector3D) Length() (float64, error) {
	x, y, z := v.GetComponents()

	r := numeric.Nrm2(numeric.Nrm2(x, y), z)
	if numeric.AreAnyOverflow(r) {
		return 0, numeric.ErrOverflow
	}
	return r, nil
}

// LengthSquared computes the squared length of the vector.
func (v *Vector3D) LengthSquared() (float64, error) {
	x, y, z := v.GetComponents()

	r := x*x + y*y + z*z
	if numeric.IsOverflow(r) {
		return 0, numeric.ErrOverflow
	}
	return r, nil
}

// Clone creates a new Vector3D with the same component values.
func (v *Vector3D) Clone() *Vector3D {
	return &Vector3D{
		X: v.GetX(),
		Y: v.GetY(),
		Z: v.GetZ(),
	}
}

// GetNormalizedVector gets the unit vector codirectional to this vector.
func (v *Vector3D) GetNormalizedVector() *Vector3D {
	w := v.Clone()
	w.Normalize()
	return w
}

// IsZeroLength returns true if the vector is of zero length (within a tolerance), false if not.
func (v *Vector3D) IsZeroLength(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}
	return v.IsEqualTo(Zero3D, tol)
}

// IsUnitLength returns true if the vector is equal to the normalized vector within the given tolerance, false if not.
func (v *Vector3D) IsUnitLength(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	return v.IsEqualTo(vv, tol)
}

// AngleTo gets the angle between this vector and another vector.
func (v *Vector3D) AngleTo(u Vector3DReader) (float64, error) {
	// code based on Kahan's formula for angles between 3D vectors
	// https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf, see Mangled Angles section
	lv, err := v.Length()
	if err != nil {
		return 0, err
	}
	lu, err := u.Length()
	if err != nil {
		return 0, err
	}

	nVu := u.Clone()
	err = nVu.Scale(lv)
	if err != nil {
		return 0, err
	}

	nUv := v.Clone()
	err = nUv.Scale(lu)
	if err != nil {
		return 0, err
	}

	// Y = norm(v) * u - norm(u) * v
	Y := nVu.Clone()
	Y.Sub(nUv)

	// X = norm(v) * u + norm(u) * v
	X := nVu.Clone()
	X.Add(nUv)

	ay, err := Y.Length()
	if err != nil {
		return 0, err
	}
	ax, err := X.Length()
	if err != nil {
		return 0, err
	}
	return 2 * math.Atan2(ay, ax), nil
}

// Dot computes the dot product between this vector and another Vector3DReader.
func (v *Vector3D) Dot(w Vector3DReader) (float64, error) {
	ax, ay, az := v.GetComponents()
	bx, by, bz := w.GetComponents()

	r := ax*bx + ay*by + az*bz
	if numeric.AreAnyOverflow(r) {
		return 0, numeric.ErrOverflow
	}

	return r, nil
}

// Cross computes the cross product between this vector and another Vector3DReader.
func (v *Vector3D) Cross(w Vector3DReader) (*Vector3D, error) {
	ax, ay, az := v.GetComponents()
	bx, by, bz := w.GetComponents()

	ux := ay*bz - az*by
	uy := az*bx - ax*bz
	uz := ax*by - ay*bx
	if numeric.AreAnyOverflow(ux, uy, uz) {
		return nil, numeric.ErrOverflow
	}

	cross := &Vector3D{
		X: ux,
		Y: uy,
		Z: uz,
	}
	return cross, nil
}

// IsEqualTo returns true if the vector components are equal within a tolerance of each other, false if not.
func (v *Vector3D) IsEqualTo(w Vector3DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vx, vy, vz := v.GetComponents()
	wx, wy, wz := w.GetComponents()

	x := math.Abs(wx - vx)
	y := math.Abs(wy - vy)
	z := math.Abs(wz - vz)
	isEqual := x <= tol && y <= tol && z <= tol
	return isEqual, nil
}

// IsParallelTo returns true if the vector is in the direction (either same or opposite) of the given vector within the given tolerance, false if not.
func (v *Vector3D) IsParallelTo(w Vector3DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	D, err := vv.Dot(ww)
	if err != nil {
		return false, err
	}

	d, err := numeric.Signum(D)
	if err != nil {
		return false, err
	}
	if d == 0 {
		return false, nil
	}

	err = vv.Scale(float64(d)) // flips vv in the direction into the ww
	if err != nil {
		return false, err
	}
	return vv.IsEqualTo(ww, tol)
}

// IsPerpendicularTo returns true if the vector is pointed in the same direction as the given vector within the given tolerance, false if not.
func (v *Vector3D) IsPerpendicularTo(w Vector3DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	d, err := vv.Dot(ww)
	if err != nil {
		return false, err
	}

	return math.Abs(d) <= tol, nil
}

// IsCodirectionalTo returns true if the vector is pointed in the same direction as the given vector within the given tolerance, false if not.
func (v *Vector3D) IsCodirectionalTo(w Vector3DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	return vv.IsEqualTo(ww, tol)
}

// Negate negates the vector components.
func (v *Vector3D) Negate() {
	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	v.SetX(-x)
	v.SetY(-y)
	v.SetZ(-z)
}

// Add adds the given displacement vector to this point.
func (v *Vector3D) Add(w Vector3DReader) error {
	vx, vy, vz := v.GetComponents()
	wx, wy, wz := w.GetComponents()

	newX := vx + wx
	newY := vy + wy
	newZ := vz + wz
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}

// Sub subtracts the given displacement vector to this point.
func (v *Vector3D) Sub(w Vector3DReader) error {
	vx, vy, vz := v.GetComponents()
	wx, wy, wz := w.GetComponents()

	newX := vx - wx
	newY := vy - wy
	newZ := vz - wz
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}

// Normalize scales the vector to unit length.
func (v *Vector3D) Normalize() error {
	x, y, z := v.GetComponents()
	l, err := v.Length()
	if err != nil {
		return err
	}
	if math.Abs(l) == 0 {
		return numeric.ErrDivideByZero
	}

	newX := x / l
	newY := y / l
	newZ := z / l
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}

// Scale scales the vector by the given factor.
func (v *Vector3D) Scale(f float64) error {
	if math.IsNaN(f) {
		return numeric.ErrInvalidArgument
	}

	x, y, z := v.GetComponents()

	newX := x * f
	newY := y * f
	newZ := z * f
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}

// MatrixTransform3D transforms this vector by left-multiplying the given matrix.
func (v *Vector3D) MatrixTransform3D(m *Matrix3D) error {
	isSingular, err := m.IsNearSingular(1e-12)
	if err != nil {
		return err
	}
	if isSingular {
		return numeric.ErrSingularMatrix
	}

	vv := v.ToBlasVector()
	mm := m.ToBlas64General()
	uu := blas64.Vector{
		N:    3,
		Data: []float64{0, 0, 0},
		Inc:  1,
	}
	blas64.Gemv(blas.NoTrans, 1, mm, vv, 0, uu)

	newX := uu.Data[0]
	newY := uu.Data[1]
	newZ := uu.Data[2]
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}

// HomogeneousMatrixTransform4D transforms this vector by left-multiplying the given matrix
// by the homogeneous vector and then projected back into this space.
func (v *Vector3D) HomogeneousMatrixTransform4D(m *Matrix4D) error {
	u := &Vector4D{X: v.X, Y: v.Y, Z: v.Z, W: 1.0}
	err := u.MatrixTransform4D(m)
	if err != nil {
		return err
	}

	ux, uy, uz, uw := u.GetComponents()
	if uw != 0 {
		return numeric.ErrDivideByZero
	}

	newX := ux / uw
	newY := uy / uw
	newZ := uz / uw
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ)
	return nil
}
