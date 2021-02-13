package geometry

// Point3DReader is a read-only interface for vectors.
type Point3DReader interface {
	GetX() float64
	GetY() float64
	GetZ() float64
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
