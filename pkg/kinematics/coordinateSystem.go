package kinematics

import (
	"github.com/tab58/v1/spatial/pkg/geometry"
)

// CoordinateSystem defines a coordinate system for referencing vectors and points.
type CoordinateSystem struct {
	origin geometry.Point3DReader
	b0     geometry.Vector3DReader
	b1     geometry.Vector3DReader
	parent *CoordinateSystem
}

// Origin returns the origin of the coordinate system.
func (c *CoordinateSystem) Origin() geometry.Point3DReader {
	return c.origin
}

// B0 returns the "first" basis vector for the coordinate system expressed in the parent coordinate system.
func (c *CoordinateSystem) B0() geometry.Vector3DReader {
	return c.b0
}

// B1 returns the "second" basis vector for the coordinate system expressed in the parent coordinate system.
func (c *CoordinateSystem) B1() geometry.Vector3DReader {
	return c.b1
}

// Rotate rotates (body-fixed) the coordinate system about an axis and angle defined in the parent coordinate system.
func (c *CoordinateSystem) Rotate(axis geometry.Vector3DReader, angle float64) error {
	transform := Transform3D{}
	err := transform.Set3DRotation(axis, angle)
	if err != nil {
		return err
	}
}
