package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/bigfloat"
	"github.com/tab58/v1/spatial/pkg/errors"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Vector2DReader is a read-only interface for a 2D vector.
type Vector2DReader interface {
	GetX() float64
	GetY() float64
	Length() (float64, error)
	LengthSquared() (float64, error)
	Angle() (float64, error)
	Clone() *Vector2D
	ToBlasVector() blas64.Vector
	GetPerpendicularVector() *Vector2D
	GetNormalizedVector() *Vector2D

	IsZeroLength(tol float64) (bool, error)
	IsUnitLength(tol float64) (bool, error)

	AngleTo(w Vector2DReader) (float64, error)
	Dot(w Vector2DReader) (float64, error)

	IsPerpendicularTo(w Vector2DReader, tol float64) (bool, error)
	IsCodirectionalTo(w Vector2DReader, tol float64) (bool, error)
	IsParallelTo(w Vector2DReader, tol float64) (bool, error)
	IsEqualTo(w Vector2DReader, tol float64) (bool, error)
}

// Vector2DWriter is a write-only interface for a 2D vector.
type Vector2DWriter interface {
	SetX(float64)
	SetY(float64)
	Negate()
	Add(w Vector2DReader) error
	Sub(w Vector2DReader) error
	Normalize() error
	Scale(f float64) error
	RotateBy(angleRad float64) error
}

// XAxis2D represents the canonical Cartesian x-axis in the 2D plane.
var XAxis2D Vector2DReader = &Vector2D{X: 1, Y: 0}

// YAxis2D represents the canonical Cartesian y-axis in the 2D plane.
var YAxis2D Vector2DReader = &Vector2D{X: 0, Y: 1}

// Zero2D represents the zero vector in the 2D plane.
var Zero2D Vector2DReader = &Vector2D{X: 0, Y: 0}

// Vector2D is a representation of a mathematical vector.
type Vector2D struct {
	X float64
	Y float64
}

// GetX returns the x-coordinate of the vector.
func (v *Vector2D) GetX() float64 {
	return v.X
}

// GetY returns the y-coordinate of the vector.
func (v *Vector2D) GetY() float64 {
	return v.Y
}

// SetX sets the x-coordinate of the vector.
func (v *Vector2D) SetX(z float64) {
	v.X = z
}

// SetY sets the y-coordinate of the vector.
func (v *Vector2D) SetY(z float64) {
	v.Y = z
}

// Clone creates a new Vector2D with the same component values.
func (v *Vector2D) Clone() *Vector2D {
	return &Vector2D{
		X: v.X,
		Y: v.Y,
	}
}

// Dot computes the dot produce between this vector and another Vector2DReader.
func (v *Vector2D) Dot(w Vector2DReader) (float64, error) {
	ax, ay := NewBigFloat(v.GetX()), NewBigFloat(v.GetY())
	bx, by := NewBigFloat(w.GetX()), NewBigFloat(w.GetY())

	r := NewBigFloat(0)
	t := NewBigFloat(0)
	t.Mul(ax, bx)
	r.Add(r, t)
	t.Mul(ay, by)
	r.Add(r, t)

	res, acc := r.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
}

// Length computes the length of the vector.
func (v *Vector2D) Length() (float64, error) {
	x, y := NewBigFloat(v.GetX()), NewBigFloat(v.GetY())
	r := bigfloat.Nrm2(x, y)

	res, acc := r.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
}

// LengthSquared computes the squared length of the vector.
func (v *Vector2D) LengthSquared() (float64, error) {
	x, y := NewBigFloat(v.GetX()), NewBigFloat(v.GetY())

	x.Mul(x, x)
	y.Mul(y, y)
	r := NewBigFloat(0)
	r.Add(x, y)

	res, acc := r.Float64()
	return res, bigfloat.HasNumericErr(res, acc)
}

// Negate negates the vector components.
func (v *Vector2D) Negate() {
	x, y := v.GetX(), v.GetY()
	v.SetX(-x)
	v.SetY(-y)
}

// Add adds the given displacement vector to this point.
func (v *Vector2D) Add(w Vector2DReader) error {
	vx, vy := v.GetX(), v.GetY()
	wx, wy := w.GetX(), w.GetY()

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

	v.SetX(newX)
	v.SetY(newY)
	return nil
}

// Sub subtracts the given displacement vector to this point.
func (v *Vector2D) Sub(w Vector2DReader) error {
	vx, vy := v.GetX(), v.GetY()
	wx, wy := w.GetX(), w.GetY()

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

	v.SetX(newX)
	v.SetY(newY)
	return nil
}

// AngleTo gets the angle between this vector and another vector.
func (v *Vector2D) AngleTo(w Vector2DReader) (float64, error) {
	// code based on Kahan's formulas for needle-like triangles
	// https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf
	lv, err := v.Length()
	if err != nil {
		return 0, nil
	}
	lw, err := w.Length()
	if err != nil {
		return 0, nil
	}

	u := v.Clone()
	u.Sub(w)
	c, err := u.Length()
	if err != nil {
		return 0, nil
	}

	a := math.Max(lv, lw)
	b := math.Min(lv, lw)

	var mu float64
	if c > b {
		mu = b - (a - c)
	} else {
		mu = c - (a - b)
	}

	t1 := (a - b) - c
	t2 := a + (b + c)
	t3 := (a - c) + b

	t := (t1 * mu) / (t2 * t3)
	return 2 * math.Atan(math.Sqrt(t)), nil
}

// Angle computes the canonical angle from the x-axis on the plane.
func (v *Vector2D) Angle() (float64, error) {
	return XAxis2D.AngleTo(v)
}

// Normalize scales the vector to unit length.
func (v *Vector2D) Normalize() error {
	x, y := v.GetX(), v.GetY()
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

	v.SetX(newX)
	v.SetY(newY)
	return nil
}

// Scale scales the vector by the given factor.
func (v *Vector2D) Scale(f float64) error {
	if math.IsNaN(f) {
		return errors.ErrInvalidArgument
	}
	x, y := v.GetX(), v.GetY()

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

	v.SetX(newX)
	v.SetY(newY)
	return nil
}

// IsEqualTo returns true if the vector components are equal within a tolerance of each other, false if not.
func (v *Vector2D) IsEqualTo(w Vector2DReader, tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	vx, vy := v.GetX(), v.GetY()
	wx, wy := w.GetX(), w.GetY()

	x, y := NewBigFloat(0), NewBigFloat(0)
	tmp := NewBigFloat(0)

	x.SetFloat64(wx)
	tmp.SetFloat64(vx)
	x.Sub(x, tmp)

	y.SetFloat64(wy)
	tmp.SetFloat64(vy)
	y.Sub(y, tmp)

	cmp := bigfloat.NewComparator()
	isEqual := cmp.LTEFloat64(x, tol) && cmp.LTEFloat64(y, tol)
	return isEqual, nil
}

// IsParallelTo returns true if the vector is in the direction (either same or opposite) of the given vector within the given tolerance, false if not.
func (v *Vector2D) IsParallelTo(w Vector2DReader, tol float64) (bool, error) {
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

// IsCodirectionalTo returns true if the vector is pointed in the same direction as the given vector within the given tolerance, false if not.
func (v *Vector2D) IsCodirectionalTo(w Vector2DReader, tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	ww := w.Clone()
	ww.Normalize()

	return vv.IsEqualTo(ww, tol)
}

// IsPerpendicularTo returns true if the vector is pointed in the same direction as the given vector within the given tolerance, false if not.
func (v *Vector2D) IsPerpendicularTo(w Vector2DReader, tol float64) (bool, error) {
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

// IsUnitLength returns true if the vector is equal to the normalized vector within the given tolerance, false if not.
func (v *Vector2D) IsUnitLength(tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}

	vv := v.Clone()
	vv.Normalize()
	return v.IsEqualTo(vv, tol)
}

// IsZeroLength returns true if the vector is of zero length (within a tolerance), false if not.
func (v *Vector2D) IsZeroLength(tol float64) (bool, error) {
	if IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}
	return v.IsEqualTo(Zero2D, tol)
}

// GetPerpendicularVector gets a vector perpendicular to this one.
func (v *Vector2D) GetPerpendicularVector() *Vector2D {
	x, y := v.GetX(), v.GetY()
	return &Vector2D{X: -y, Y: x}
}

// GetNormalizedVector gets the unit vector codirectional to this vector.
func (v *Vector2D) GetNormalizedVector() *Vector2D {
	w := v.Clone()
	w.Normalize()
	return w
}

// RotateBy rotates the vector by the angle specified. The positive direction is counterclockwise about the cross product of the axes.
func (v *Vector2D) RotateBy(angleRad float64) error {
	if math.IsNaN(angleRad) {
		return bigfloat.ErrNaN
	}
	c := math.Cos(angleRad)
	s := math.Sin(angleRad)
	R := blas64.General{
		Rows:   2,
		Cols:   2,
		Stride: 2,
		Data:   []float64{c, -s, s, c},
	}
	w := v.ToBlasVector()

	Rw := Zero2D.Clone().ToBlasVector()
	blas64.Gemv(blas.NoTrans, 1, R, w, 1, Rw)

	// TODO: do we need some check here?
	newX, newY := Rw.Data[0], Rw.Data[1]
	v.SetX(newX)
	v.SetY(newY)
	return nil
}

// ToBlasVector returns a BLAS vector for operations.
func (v *Vector2D) ToBlasVector() blas64.Vector {
	return blas64.Vector{
		N:    2,
		Data: []float64{v.X, v.Y},
		Inc:  1,
	}
}
