package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/numeric"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Vector4DReader is a read-only interface for a 4D vector.
type Vector4DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64
	GetW() float64

	GetComponents() (float64, float64, float64, float64)
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
	ProjectToVector3D() (*Vector3D, error)
}

// Vector4DWriter is a write-only interface for a 4D vector.
type Vector4DWriter interface {
	SetX(float64)
	SetY(float64)
	SetZ(float64)
	SetW(float64)

	SetComponents(x, y, z, w float64)

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

// GetX returns the x-coordinate of the vector.
func (v *Vector4D) GetX() float64 { return v.X }

// GetY returns the y-coordinate of the vector.
func (v *Vector4D) GetY() float64 { return v.Y }

// GetZ returns the z-coordinate of the vector.
func (v *Vector4D) GetZ() float64 { return v.Z }

// GetW returns the w-coordinate of the vector.
func (v *Vector4D) GetW() float64 { return v.W }

// GetComponents returns the components of the vector.
func (v *Vector4D) GetComponents() (x, y, z, w float64) {
	return v.GetX(), v.GetY(), v.GetZ(), v.GetW()
}

// Length computes the length of the vector.
func (v *Vector4D) Length() (float64, error) {
	x, y, z, w := v.GetComponents()

	r := numeric.Nrm2(numeric.Nrm2(numeric.Nrm2(x, y), z), w)
	if numeric.IsOverflow(r) {
		return 0, numeric.ErrOverflow
	}
	return r, nil
}

// LengthSquared computes the squared length of the vector.
func (v *Vector4D) LengthSquared() (float64, error) {
	x, y, z, w := v.GetComponents()

	r := x*x + y*y + z*z + w*w
	if numeric.IsOverflow(r) {
		return 0, numeric.ErrOverflow
	}
	return r, nil
}

// Clone creates a new Vector3D with the same component values.
func (v *Vector4D) Clone() *Vector4D {
	return &Vector4D{
		X: v.GetX(),
		Y: v.GetY(),
		Z: v.GetZ(),
		W: v.GetW(),
	}
}

// ToBlasVector returns a BLAS vector for operations.
func (v *Vector4D) ToBlasVector() blas64.Vector {
	return blas64.Vector{
		N:    4,
		Data: []float64{v.X, v.Y, v.Z, v.W},
		Inc:  1,
	}
}

// GetNormalizedVector gets the unit vector codirectional to this vector.
func (v *Vector4D) GetNormalizedVector() *Vector4D {
	w := v.Clone()
	w.Normalize()
	return w
}

// IsZeroLength returns true if the vector is of zero length (within a tolerance), false if not.
func (v *Vector4D) IsZeroLength(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}
	return v.IsEqualTo(Zero4D, tol)
}

// IsUnitLength returns true if the vector is equal to the normalized vector within the given tolerance, false if not.
func (v *Vector4D) IsUnitLength(tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	return v.IsEqualTo(vv, tol)
}

// AngleTo gets the angle between this vector and another vector.
func (v *Vector4D) AngleTo(u Vector4DReader) (float64, error) {
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
func (v *Vector4D) Dot(w Vector4DReader) (float64, error) {
	ax, ay, az, aw := v.GetComponents()
	bx, by, bz, bw := v.GetComponents()

	r := ax*bx + ay*by + az*bz + aw*bw
	if numeric.AreAnyOverflow(r) {
		return 0, numeric.ErrOverflow
	}

	return r, nil
}

// IsPerpendicularTo returns true if the vector is pointed in the same direction as the given vector within the given tolerance, false if not.
func (v *Vector4D) IsPerpendicularTo(w Vector4DReader, tol float64) (bool, error) {
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
func (v *Vector4D) IsCodirectionalTo(w Vector4DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	return vv.IsEqualTo(ww, tol)
}

// IsParallelTo returns true if the vector is in the direction (either same or opposite) of the given vector within the given tolerance, false if not.
func (v *Vector4D) IsParallelTo(w Vector4DReader, tol float64) (bool, error) {
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

// IsEqualTo returns true if the vector components are equal within a tolerance of each other, false if not.
func (v *Vector4D) IsEqualTo(w Vector4DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	vx, vy, vz, vw := v.GetComponents()
	wx, wy, wz, ww := w.GetComponents()

	x := math.Abs(wx - vx)
	y := math.Abs(wy - vy)
	z := math.Abs(wz - vz)
	dw := math.Abs(ww - vw)
	isEqual := x <= tol && y <= tol && z <= tol && dw <= tol
	return isEqual, nil
}

// MatrixTransform4D transforms this vector by left-multiplying the given matrix.
func (v *Vector4D) MatrixTransform4D(m *Matrix4D) error {
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
		N:    4,
		Data: []float64{0, 0, 0, 0},
		Inc:  1,
	}
	blas64.Gemv(blas.NoTrans, 1, mm, vv, 0, uu)

	newX := uu.Data[0]
	newY := uu.Data[1]
	newZ := uu.Data[2]
	newW := uu.Data[3]
	if numeric.AreAnyOverflow(newX, newY, newZ, newW) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ, newW)
	return nil
}

// SetX sets the x-coordinate of the vector.
func (v *Vector4D) SetX(z float64) { v.X = z }

// SetY sets the y-coordinate of the vector.
func (v *Vector4D) SetY(z float64) { v.Y = z }

// SetZ sets the z-coordinate of the vector.
func (v *Vector4D) SetZ(z float64) { v.Z = z }

// SetW sets the w-coordinate of the vector.
func (v *Vector4D) SetW(z float64) { v.W = z }

// SetComponents sets the components of the vector.
func (v *Vector4D) SetComponents(x, y, z, w float64) {
	v.SetX(x)
	v.SetY(y)
	v.SetZ(z)
	v.SetW(w)
}

// Negate negates the vector components.
func (v *Vector4D) Negate() {
	x, y, z, w := v.GetComponents()
	v.SetComponents(-x, -y, -z, -w)
}

// Add adds the given displacement vector to this point.
func (v *Vector4D) Add(w Vector4DReader) error {
	vx, vy, vz, vw := v.GetComponents()
	wx, wy, wz, ww := w.GetComponents()

	newX := vx + wx
	newY := vy + wy
	newZ := vz + wz
	newW := vw + ww
	if numeric.AreAnyOverflow(newX, newY, newZ, newW) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ, newW)
	return nil
}

// Sub subtracts the given displacement vector to this point.
func (v *Vector4D) Sub(w Vector4DReader) error {
	vx, vy, vz, vw := v.GetComponents()
	wx, wy, wz, ww := w.GetComponents()

	newX := vx - wx
	newY := vy - wy
	newZ := vz - wz
	newW := vw - ww
	if numeric.AreAnyOverflow(newX, newY, newZ, newW) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ, newW)
	return nil
}

// Normalize scales the vector to unit length.
func (v *Vector4D) Normalize() error {
	x, y, z, w := v.GetComponents()
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
	newW := w / l
	if numeric.AreAnyOverflow(newX, newY, newZ, newW) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ, newW)
	return nil
}

// Scale scales the vector by the given factor.
func (v *Vector4D) Scale(f float64) error {
	if math.IsNaN(f) {
		return numeric.ErrInvalidArgument
	}

	x, y, z, w := v.GetComponents()

	newX := x * f
	newY := y * f
	newZ := z * f
	newW := w * f
	if numeric.AreAnyOverflow(newX, newY, newZ, newW) {
		return numeric.ErrOverflow
	}

	v.SetComponents(newX, newY, newZ, newW)
	return nil
}

// ProjectToVector3D treats this vector as a homogeneous projection of a 3D vector and projects this vector back to 3D space.
func (v *Vector4D) ProjectToVector3D() (*Vector3D, error) {
	x, y, z, w := v.GetComponents()

	if w == 0 {
		return nil, numeric.ErrDivideByZero
	}

	newX := x / w
	newY := y / w
	newZ := z / w
	if numeric.AreAnyOverflow(newX, newY, newZ) {
		return nil, numeric.ErrOverflow
	}

	u := &Vector3D{
		X: newX,
		Y: newY,
		Z: newZ,
	}
	return u, nil
}
