package geometry

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/numeric"
)

// Point3DReader is a read-only interface for vectors.
type Point3DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64

	Clone() *Point3D
	AsVector() *Vector3D
	DistanceTo(q Point3DReader) (float64, error)
	IsEqualTo(q Point3DReader, tol float64) (bool, error)
}

// Point3DWriter is a write-only interface for vectors.
type Point3DWriter interface {
	SetX(float64)
	SetY(float64)
	SetZ(float64)
}

// Origin3D is the canonical origin in the 3D plane.
var Origin3D Point3DReader = &Point3D{X: 0, Y: 0, Z: 0}

// Point3D represents a 3D point.
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// GetX returns the x-coordinate of the point.
func (p *Point3D) GetX() float64 {
	return p.X
}

// GetY returns the y-coordinate of the point.
func (p *Point3D) GetY() float64 {
	return p.Y
}

// GetZ returns the z-coordinate of the point.
func (p *Point3D) GetZ() float64 {
	return p.Z
}

// SetX sets the x-coordinate of the point.
func (p *Point3D) SetX(x float64) {
	p.X = x
}

// SetY sets the y-coordinate of the point.
func (p *Point3D) SetY(y float64) {
	p.Y = y
}

// SetZ sets the z-coordinate of the point.
func (p *Point3D) SetZ(z float64) {
	p.Z = z
}

// Clone creates a new Point3D with the same coordinate information.
func (p *Point3D) Clone() *Point3D {
	return &Point3D{
		X: p.GetX(),
		Y: p.GetY(),
		Z: p.GetZ(),
	}
}

// AsVector creates a displacement vector from the origin to this point.
func (p *Point3D) AsVector() *Vector3D {
	return &Vector3D{
		X: p.GetX(),
		Y: p.GetY(),
		Z: p.GetZ(),
	}
}

// DistanceTo calculates the distance from one point to another.
func (p *Point3D) DistanceTo(q Point3DReader) (float64, error) {
	v := q.AsVector()
	err := v.Sub(p.AsVector())
	if err != nil {
		return 0, err
	}
	return v.Length()
}

// IsEqualTo returns true if 2 points can be considered equal to within a specific tolerance, false if not.
func (p *Point3D) IsEqualTo(q Point3DReader, tol float64) (bool, error) {
	if numeric.IsInvalidTolerance(tol) {
		return false, numeric.ErrInvalidTol
	}

	px, py, pz := p.GetX(), p.GetY(), p.GetZ()
	qx, qy, qz := q.GetX(), q.GetY(), q.GetZ()

	resX := math.Abs(qx - px)
	resY := math.Abs(qy - py)
	resZ := math.Abs(qz - pz)
	isEqual := resX <= tol && resY <= tol && resZ <= tol
	return isEqual, nil
}
