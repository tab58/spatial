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

// GetLocalOrientation returns the rotation matrix that defines the orientation of the coordinate system from the parent coordinate system.
func (c *CoordinateSystem) GetLocalOrientation() (*geometry.Matrix3D, error) {
	b0 := c.b0

	// calculate z-axis
	b2, err := b0.Cross(c.b1)
	if err != nil {
		return nil, err
	}

	// do just in case b0 is not orthogonal to c.b1
	b1, err := b0.Cross(b2)
	if err != nil {
		return nil, err
	}

	b0x, b0y, b0z := b0.GetComponents()
	b1x, b1y, b1z := b1.GetComponents()
	b2x, b2y, b2z := b2.GetComponents()

	m := &geometry.Matrix3D{}
	m.SetElements(b0x, b0y, b0z, b1x, b1y, b1z, b2x, b2y, b2z)

	return m, nil
}

// GetGlobalOrientation returns the orientation of the coordinate system from the global system.
func (c *CoordinateSystem) GetGlobalOrientation() (*geometry.Matrix3D, error) {
	global := &geometry.Matrix3D{}
	global.Identity()

	for current := c; current != nil; current = current.parent {
		localOrient, err := current.GetLocalOrientation()
		if err != nil {
			return nil, err
		}
		err = global.Premultiply(localOrient)
		if err != nil {
			return nil, err
		}
	}

	return global, nil
}
