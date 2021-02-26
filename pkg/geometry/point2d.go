package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/errors"
	"github.com/tab58/v1/spatial/pkg/numeric"
)

// Point2DReader is a write-only interface for vectors.
type Point2DReader interface {
	GetX() float64
	GetY() float64

	Clone() *Point2D
	AsVector() *Vector2D
	DistanceTo(q Point2DReader) (float64, error)
	IsEqualTo(q Point2DReader, tol float64) (bool, error)
}

// Point2DWriter is a read-only interface for vectors.
type Point2DWriter interface {
	SetX() float64
	SetY() float64

	Copy(q Point2DReader)
	Scale(f float64) error
	Add(v Vector2DReader) error
	Sub(v Vector2DReader) error
}

// Origin2D is the canonical origin in the 2D plane.
var Origin2D Point2DReader = &Point2D{X: 0, Y: 0}

// Point2D represents a 2D point.
type Point2D struct {
	X float64
	Y float64
}

// GetX returns the x-coordinate of the point.
func (p *Point2D) GetX() float64 {
	return p.X
}

// GetY returns the y-coordinate of the point.
func (p *Point2D) GetY() float64 {
	return p.Y
}

// SetX sets the x-coordinate of the point.
func (p *Point2D) SetX(x float64) {
	p.X = x
}

// SetY sets the x-coordinate of the point.
func (p *Point2D) SetY(y float64) {
	p.Y = y
}

// Copy copies the coordinate info from another Point2D.
func (p *Point2D) Copy(q Point2DReader) {
	p.SetX(q.GetX())
	p.SetY(q.GetY())
}

// Scale scales the displacement vector from the origin by the given factor.
func (p *Point2D) Scale(f float64) error {
	if math.IsNaN(f) {
		return errors.ErrInvalidArgument
	}

	newX := p.GetX() * f
	newY := p.GetY() * f
	if math.IsInf(newX, 0) || math.IsInf(newY, 0) {
		return errors.ErrOverflow
	}

	p.SetX(newX)
	p.SetY(newY)
	return nil
}

// Add adds the given displacement vector to this point.
func (p *Point2D) Add(v Vector2DReader) error {
	x, y := p.GetX(), p.GetY()
	vx, vy := v.GetX(), v.GetY()

	newX := x + vx
	newY := y + vy
	if math.IsInf(newX, 0) || math.IsInf(newY, 0) {
		return errors.ErrOverflow
	}

	p.SetX(newX)
	p.SetY(newY)
	return nil
}

// Sub subtracts the given displacement vector to this point.
func (p *Point2D) Sub(v Vector2DReader) error {
	x, y := p.GetX(), p.GetY()
	vx, vy := v.GetX(), v.GetY()

	newX := x - vx
	newY := y - vy
	if math.IsInf(newX, 0) || math.IsInf(newY, 0) {
		return errors.ErrOverflow
	}

	p.SetX(newX)
	p.SetY(newY)
	return nil
}

// Clone creates a new Point2D with the same coordinate data.
func (p *Point2D) Clone() *Point2D {
	return &Point2D{
		X: p.GetX(),
		Y: p.GetY(),
	}
}

// AsVector gets the displacement vector for the point.
func (p *Point2D) AsVector() *Vector2D {
	return &Vector2D{
		X: p.GetX(),
		Y: p.GetY(),
	}
}

// DistanceTo gets the length of the displacement vector from a point to the given point.
func (p *Point2D) DistanceTo(q Point2DReader) (float64, error) {
	qx, qy := q.GetX(), q.GetY()
	px, py := p.GetX(), p.GetY()

	newX := qx - px
	newY := qy - py
	if math.IsInf(newX, 0) || math.IsInf(newY, 0) {
		return 0, errors.ErrOverflow
	}

	len := numeric.Nrm2(newX, newY)
	if math.IsInf(len, 0) {
		return 0, errors.ErrOverflow
	}
	return len, nil
}

// IsEqualTo returns true if 2 points can be considered equal to within a specific tolerance, false if not.
func (p *Point2D) IsEqualTo(q Point2DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, errors.ErrInvalidTol
	}
	px, py := p.GetX(), p.GetY()
	qx, qy := q.GetX(), q.GetY()

	x := math.Abs(px - qx)
	y := math.Abs(py - qy)
	isEqual := x <= tol && y <= tol
	return isEqual, nil
}
