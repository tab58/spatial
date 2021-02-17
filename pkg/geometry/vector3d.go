package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/bigfloat"
	"github.com/tab58/v1/spatial/pkg/blasmatrix"
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Vector3DReader is a read-only interface for a 3D vector.
type Vector3DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64

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
	x, y, z := NewBigFloat(v.GetX()), NewBigFloat(v.GetY()), NewBigFloat(v.GetZ())
	r := bigfloat.Nrm2(bigfloat.Nrm2(x, y), z)

	res, acc := r.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
}

// LengthSquared computes the squared length of the vector.
func (v *Vector3D) LengthSquared() (float64, error) {
	x, y, z := NewBigFloat(v.GetX()), NewBigFloat(v.GetY()), NewBigFloat(v.GetZ())

	x.Mul(x, x)
	y.Mul(y, y)
	z.Mul(z, z)
	r := NewBigFloat(0)
	r.Add(x, y)
	r.Add(r, z)

	res, acc := r.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
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
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}
	return v.IsEqualTo(Zero3D, tol)
}

// IsUnitLength returns true if the vector is equal to the normalized vector within the given tolerance, false if not.
func (v *Vector3D) IsUnitLength(tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
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
	ax, ay, az := NewBigFloat(v.GetX()), NewBigFloat(v.GetY()), NewBigFloat(v.GetZ())
	bx, by, bz := NewBigFloat(w.GetX()), NewBigFloat(w.GetY()), NewBigFloat(w.GetZ())

	rAcc := NewBigFloat(0)
	t := NewBigFloat(0)
	t.Mul(ax, bx)
	rAcc.Add(rAcc, t)
	t.Mul(ay, by)
	rAcc.Add(rAcc, t)
	t.Mul(az, bz)
	rAcc.Add(rAcc, t)

	res, acc := rAcc.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
}

// Cross computes the cross product between this vector and another Vector3DReader.
func (v *Vector3D) Cross(w Vector3DReader) (*Vector3D, error) {
	ax, ay, az := NewBigFloat(v.GetX()), NewBigFloat(v.GetY()), NewBigFloat(v.GetZ())
	bx, by, bz := NewBigFloat(w.GetX()), NewBigFloat(w.GetY()), NewBigFloat(w.GetZ())

	t1, t2, r := NewBigFloat(0), NewBigFloat(0), NewBigFloat(0)

	t1.Mul(ay, bz)
	t2.Mul(az, by)
	r.Sub(t1, t2)
	ux, acc := r.Float64()
	err := bigfloat.HasNumericErr(ux, acc)
	if err != nil {
		return nil, err
	}

	t1.Mul(ax, bz)
	t2.Mul(az, bx)
	r.Sub(t1, t2)
	r.Neg(r)
	uy, acc := r.Float64()
	err = bigfloat.HasNumericErr(uy, acc)
	if err != nil {
		return nil, err
	}

	t1.Mul(ax, by)
	t2.Mul(ay, bx)
	r.Sub(t1, t2)
	uz, acc := r.Float64()
	err = bigfloat.HasNumericErr(uz, acc)
	if err != nil {
		return nil, err
	}

	return &Vector3D{
		X: ux,
		Y: uy,
		Z: uz,
	}, nil
}

// IsEqualTo returns true if the vector components are equal within a tolerance of each other, false if not.
func (v *Vector3D) IsEqualTo(w Vector3DReader, tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	vx, vy, vz := v.GetX(), v.GetY(), v.GetZ()
	wx, wy, wz := w.GetX(), w.GetY(), v.GetZ()

	x, y, z := NewBigFloat(0), NewBigFloat(0), NewBigFloat(0)
	tmp := NewBigFloat(0)

	x.SetFloat64(wx)
	tmp.SetFloat64(vx)
	x.Sub(x, tmp)

	y.SetFloat64(wy)
	tmp.SetFloat64(vy)
	y.Sub(y, tmp)

	z.SetFloat64(wz)
	tmp.SetFloat64(vz)
	z.Sub(z, tmp)

	cmp := bigfloat.NewComparator()
	isEqual := cmp.LTEFloat64(x, tol) && cmp.LTEFloat64(y, tol) && cmp.LTEFloat64(z, tol)
	return isEqual, nil
}

// IsParallelTo returns true if the vector is in the direction (either same or opposite) of the given vector within the given tolerance, false if not.
func (v *Vector3D) IsParallelTo(w Vector3DReader, tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	D, err := vv.Dot(ww)
	if err != nil {
		return false, err
	}

	d, err := Signum(D)
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
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
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
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
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
	vx, vy, vz := v.GetX(), v.GetY(), v.GetZ()
	wx, wy, wz := w.GetX(), w.GetY(), v.GetZ()

	c := bigfloat.NewCalculator(v.GetX())
	c.SetFloat64(vx)
	c.Add(wx)
	newX, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(vy)
	c.Add(wy)
	newY, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(vz)
	c.Add(wz)
	newZ, err := c.Float64()
	if err != nil {
		return err
	}

	v.SetX(newX)
	v.SetY(newY)
	v.SetZ(newZ)
	return nil
}

// Sub subtracts the given displacement vector to this point.
func (v *Vector3D) Sub(w Vector3DReader) error {
	vx, vy, vz := v.GetX(), v.GetY(), v.GetZ()
	wx, wy, wz := w.GetX(), w.GetY(), v.GetZ()

	c := bigfloat.NewCalculator(v.GetX())
	c.SetFloat64(vx)
	c.Sub(wx)
	newX, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(vy)
	c.Sub(wy)
	newY, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(vz)
	c.Sub(wz)
	newZ, err := c.Float64()
	if err != nil {
		return err
	}

	v.SetX(newX)
	v.SetY(newY)
	v.SetZ(newZ)
	return nil
}

// Normalize scales the vector to unit length.
func (v *Vector3D) Normalize() error {
	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	l, err := v.Length()
	if err != nil {
		return err
	}
	if math.Abs(l) == 0 {
		return errors.ErrDivideByZero
	}

	c := bigfloat.NewCalculator(x).Quo(l)
	newX, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(y).Quo(l)
	newY, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(z).Quo(l)
	newZ, err := c.Float64()
	if err != nil {
		return err
	}

	v.SetX(newX)
	v.SetY(newY)
	v.SetZ(newZ)
	return nil
}

// Scale scales the vector by the given factor.
func (v *Vector3D) Scale(f float64) error {
	if math.IsNaN(f) {
		return errors.ErrInvalidArgument
	}
	x, y, z := v.GetX(), v.GetY(), v.GetZ()

	c := bigfloat.NewCalculator(x).Mul(f)
	newX, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(y).Mul(f)
	newY, err := c.Float64()
	if err != nil {
		return err
	}

	c.SetFloat64(z).Mul(f)
	newZ, err := c.Float64()
	if err != nil {
		return err
	}

	v.SetX(newX)
	v.SetY(newY)
	v.SetZ(newZ)
	return nil
}

// RotateBy rotates the vector about the specified axis by the angle specified.
// The positive direction is counterclockwise about the specified axis of rotation.
func (v *Vector3D) RotateBy(axis Vector3DReader, angleRad float64) error {
	rotMat, err := blasmatrix.Get3DRotMatrix(axis, angleRad)
	if err != nil {
		return err
	}

	vv := v.ToBlasVector()
	tmp := Zero3D.ToBlasVector()
	blas64.Gemv(blas.NoTrans, 1, *rotMat, vv, 0, tmp)

	x, y, z := tmp.Data[0], tmp.Data[1], tmp.Data[2]
	v.SetX(x)
	v.SetY(y)
	v.SetZ(z)
	return nil
}
