package kinematics

import "github.com/tab58/v1/spatial/pkg/geometry"

// HomogeneousTransform4D is a 4x4 matrix that encodes a transformation.
type HomogeneousTransform4D struct {
	*geometry.Matrix4D
}

// Set3DTranslation sets the matrix to a translation encoded for homogeneous coordinates.
func (m *HomogeneousTransform4D) Set3DTranslation(v geometry.Vector3DReader) error {
	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	return m.Matrix4D.Adjoint().SetElements(1, 0, 0, x, 0, 1, 0, y, 0, 0, 1, z, 0, 0, 0, 1)
}

// Set3DRotation sets the matrix to a rotation encoded for homogeneous coordinates.
func (m *HomogeneousTransform4D) Set3DRotation(axis geometry.Vector3DReader, angle float64) error {
	e := rotation3DFromAxisAngle(axis, angle)
	return m.Matrix4D.SetElements(e[0], e[1], e[2], 0, e[3], e[4], e[5], 0, e[6], e[7], e[8], 0, 0, 0, 0, 1)
}

// Set3DTranslationRotation sets the matrix to a rotation and translation encoded for homogeneous coordinates.
func (m *HomogeneousTransform4D) Set3DTranslationRotation(axis geometry.Vector3DReader, angle float64, v geometry.Vector3DReader) error {
	x, y, z := v.GetX(), v.GetY(), v.GetZ()
	e := rotation3DFromAxisAngle(axis, angle)
	return m.Matrix4D.SetElements(e[0], e[1], e[2], x, e[3], e[4], e[5], y, e[6], e[7], e[8], z, 0, 0, 0, 1)
}
