package geometry

import "gonum.org/v1/gonum/blas/blas64"

// Vector3DReader is a read-only interface for a 3D vector.
type Vector3DReader interface {
	X() float64
	Y() float64
	Z() float64
	ToBlasVector() blas64.Vector
}

// Vector3DWriter is a write-only interface for a 3D vector.
type Vector3DWriter interface {
	SetX(float64)
	SetY(float64)
	SetZ(float64)
}

type vector3DImpl struct {
	x float64
	y float64
	z float64
}

func (v vector3DImpl) X() float64 {
	return v.x
}

func (v vector3DImpl) Y() float64 {
	return v.y
}

func (v vector3DImpl) Z() float64 {
	return v.z
}

func (v *vector3DImpl) SetX(z float64) {
	v.x = z
}

func (v *vector3DImpl) SetY(z float64) {
	v.y = z
}

func (v *vector3DImpl) SetZ(z float64) {
	v.z = z
}

func (v vector3DImpl) ToBlasVector() blas64.Vector {
	return blas64.Vector{
		N:    3,
		Data: []float64{v.x, v.y, v.z},
		Inc:  1,
	}
}
