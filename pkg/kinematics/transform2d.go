package kinematics

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/geometry"
)

// Transform2D is a 2x2 matrix that encodes a transformation.
type Transform2D struct {
	*geometry.Matrix2D
}

// Set2DMirror sets the matrix to encode a mirror operation for a vector about a line direction defined by the given vector n.
func (m *Transform2D) Set2DMirror(n geometry.Vector2DReader) error {
	n1, n2 := n.GetX(), n.GetY()

	n1sq := n1 * n1
	n1n2 := n1 * n2
	n2sq := n2 * n2

	return m.Matrix2D.SetElements(1-2*n1sq, -2*n1n2, -2*n1n2, 1-2*n2sq)
}

// Set2DRotation sets the matrix to encode a rotation about the plane normal with the given angle.
func (m *Transform2D) Set2DRotation(angle float64) error {
	c := math.Cos(angle)
	s := math.Sin(angle)
	m.Matrix2D.SetElements(c, -s, s, c)
	return nil
}

// Set2DScaling sets the matrix to a scaling matrix.
func (m *Transform2D) Set2DScaling(v geometry.Vector2DReader) error {
	x, y := v.GetX(), v.GetY()
	return m.Matrix2D.SetElements(x, 0, 0, y)
}
