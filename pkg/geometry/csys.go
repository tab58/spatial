package geometry

// CoordinateSystem defines a coordinate system for referencing vectors and points.
type CoordinateSystem struct {
	origin Point3DReader
	b0     Vector3DReader
	b1     Vector3DReader
	parent *CoordinateSystem
}

// Origin returns the origin of the coordinate system.
func (c *CoordinateSystem) Origin() Point3DReader {
	return c.origin
}

// B0 returns the "first" basis vector for the coordinate system expressed in the parent coordinate system.
func (c *CoordinateSystem) B0() Vector3DReader {
	return c.b0
}

// B1 returns the "second" basis vector for the coordinate system expressed in the parent coordinate system.
func (c *CoordinateSystem) B1() Vector3DReader {
	return c.b1
}

// Rotate rotates the coordinate system about an axis and angle defined in the parent coordinate system.
func (c *CoordinateSystem) Rotate(axis Vector3DReader, angle float64) {

}
