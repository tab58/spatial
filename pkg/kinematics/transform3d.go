package kinematics

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/geometry"
)

// Transform3D is a 3x3 matrix that encodes a transformation.
type Transform3D struct {
	*geometry.Matrix3D
}

func rotation3DFromAxisAngle(axis geometry.Vector3DReader, angle float64) [9]float64 {
	elements := [9]float64{}
	u := axis.Clone()
	u.Normalize()
	ux, uy, uz := u.GetX(), u.GetY(), u.GetZ()

	c := math.Cos(angle)
	s := math.Sin(angle)
	c1 := 1.0 - c

	elements[0] = c + ux*ux*c1
	elements[1] = ux*uy*c1 - uz*s
	elements[2] = ux*uz*c1 + uy*s

	elements[3] = ux*uy*c1 + uz*s
	elements[4] = c + uy*uy*c1
	elements[5] = uy*uz*c1 - ux*s

	elements[6] = ux*uz*c1 - uy*s
	elements[7] = uy*uz*c1 + ux*s
	elements[8] = c + uz*uz*c1
	return elements
}

// Set3DRotation sets the matrix to a 3D rotation about the specified axis and angle.
func (m *Transform3D) Set3DRotation(axis geometry.Vector3DReader, angle float64) error {
	u := axis.Clone()
	u.Normalize()
	ux, uy, uz := u.GetX(), u.GetY(), u.GetZ()

	c := math.Cos(angle)
	s := math.Sin(angle)
	c1 := 1.0 - c

	a00 := c + ux*ux*c1
	a01 := ux*uy*c1 - uz*s
	a02 := ux*uz*c1 + uy*s

	a10 := ux*uy*c1 + uz*s
	a11 := c + uy*uy*c1
	a12 := uy*uz*c1 - ux*s

	a20 := ux*uz*c1 - uy*s
	a21 := uy*uz*c1 + ux*s
	a22 := c + uz*uz*c1
	m.Matrix3D.SetElements(a00, a01, a02, a10, a11, a12, a20, a21, a22)
	return nil
}

// Set3DXRotation sets the matrix to a 3D rotation about x-axis with the specified angle.
func (m *Transform3D) Set3DXRotation(angle float64) error {
	c := math.Cos(angle)
	s := math.Sin(angle)
	m.Matrix3D.SetElements(1, 0, 0, 0, c, -s, 0, s, c)
	return nil
}

// Set3DYRotation sets the matrix to a 3D rotation about y-axis with the specified angle.
func (m *Transform3D) Set3DYRotation(angle float64) error {
	c := math.Cos(angle)
	s := math.Sin(angle)
	m.Matrix3D.SetElements(c, 0, s, 0, 1, 0, -s, 0, c)
	return nil
}

// Set3DZRotation sets the matrix to a 3D rotation about z-axis with the specified angle.
func (m *Transform3D) Set3DZRotation(angle float64) error {
	c := math.Cos(angle)
	s := math.Sin(angle)
	m.Matrix3D.SetElements(c, -s, 0, s, c, 0, 0, 0, 1)
	return nil
}

// Set3DScaling sets the matrix to a scaling matrix.
func (m *Transform3D) Set3DScaling(v geometry.Vector3DReader) error {
	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	return m.Matrix3D.SetElements(x, 0, 0, 0, y, 0, 0, 0, z)
}

// Set3DMirror sets the matrix to encode a mirror operation for a vector about a line direction defined by the given vector n.
func (m *Transform3D) Set3DMirror(n geometry.Vector3DReader) error {
	n1, n2, n3 := n.GetX(), n.GetY(), n.GetZ()

	n1sq := n1 * n1
	n1n2 := n1 * n2
	n1n3 := n1 * n3
	n2sq := n2 * n2
	n2n3 := n2 * n3
	n3sq := n3 * n3

	return m.Matrix3D.SetElements(1-2*n1sq, -2*n1n2, -2*n1n3, -2*n1n2, 1-2*n2sq, -2*n2n3, -2*n1n3, -2*n2n3, 1-2*n3sq)
}
