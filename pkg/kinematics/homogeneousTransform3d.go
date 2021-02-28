package kinematics

import (
	"math"

	"github.com/tab58/v1/spatial/pkg/geometry"
)

// HomogeneousTransform3D is a 3x3 matrix that encodes a transformation for vectors in homogeneous space.
type HomogeneousTransform3D struct {
	*geometry.Matrix3D
}

// Set2DTranslation sets the matrix to a translation encoded for homogeneous coordinates.
func (t *HomogeneousTransform3D) Set2DTranslation(v geometry.Vector2DReader) error {
	x, y := v.GetX(), v.GetY()
	return t.SetElements(1, 0, x, 0, 1, y, 0, 0, 1)
}

// Set2DRotation sets the matrix to a 2D rotation encoded for homogeneous coordinates.
func (t *HomogeneousTransform3D) Set2DRotation(angle float64) error {
	c, s := math.Cos(angle), math.Sin(angle)
	return t.SetElements(c, -s, 0, s, c, 0, 0, 0, 1)
}

// Set2DTranslationRotation sets the matrix to a rotation and a translation encoded for homogeneous coordinates.
func (t *HomogeneousTransform3D) Set2DTranslationRotation(v geometry.Vector2DReader, angle float64) error {
	x, y := v.GetX(), v.GetY()
	c, s := math.Cos(angle), math.Sin(angle)
	return t.SetElements(c, -s, x, s, c, y, 0, 0, 1)
}
