package geometry

import (
	"gonum.org/v1/gonum/blas/blas64"
)

// matrixAxpy defines an addition-type operation.
type matrixAxpy struct {
	m     blas64.General
	alpha float64
}

// type MatrixElementOp func(mat blas64.General, alpha float64) matrixAxpy

// CloneBlasMatrix deep clones a BLAS matrix.
func CloneBlasMatrix(m blas64.General) blas64.General {
	mData := m.Data
	n := len(mData)
	data := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		data[i] = mData[i]
	}

	return blas64.General{
		Rows:   m.Rows,
		Cols:   m.Cols,
		Stride: m.Stride,
		Data:   data,
	}
}

// AddMatrixAxpy computes Y = alpha * X + beta * Y.
func AddMatrixAxpy(X, Y blas64.General, alpha, beta float64) (blas64.General, error) {
	aM, aN := X.Rows, X.Cols
	bM, bN := Y.Rows, Y.Cols

	xData := X.Data
	yData := Y.Data

	if aM != bM || aN != bN {
		return blas64.General{}, ErrMatrixDims
	}

	K := aM * aN
	for i := 0; i < K; i++ {
		xi := xData[i]
		yi := yData[i]
		Y.Data[i] = alpha*xi + beta*yi
	}

	return Y, nil
}

// MatrixAxpy adds an array of matrices with a scale factor: result = m1*alpha1+m2*alpha2+...
func MatrixAxpy(matrixInfos ...matrixAxpy) (blas64.General, error) {
	N := len(matrixInfos)
	if N == 0 {
		return blas64.General{}, ErrEmptyArray
	}

	firstInfo := matrixInfos[0].m
	m, n, s := firstInfo.Rows, firstInfo.Cols, firstInfo.Stride
	k := m * n
	r := make([]float64, 0, k)
	for i := 0; i < k; i++ {
		r[i] = 0
	}
	result := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: s,
		Data:   r,
	}

	for _, mInfo := range matrixInfos {
		m, alpha := mInfo.m, mInfo.alpha
		_, err := AddMatrixAxpy(m, result, alpha, 1)
		if err != nil {
			return blas64.General{}, err
		}
	}
	return result, nil
}

// BuildMatrix3DSkewSymmetric builds the skew symmetric matrix based on vector component values.
func BuildMatrix3DSkewSymmetric(v Vector3DReader) blas64.General {
	// TODO: do check on v's Length, etc.

	x, y, z := v.X(), v.Y(), v.Z()
	return blas64.General{
		Rows:   3,
		Cols:   3,
		Data:   []float64{0, -z, y, z, 0, -x, -y, x, 0},
		Stride: 3,
	}
}

// func get3DRotMatrix(axis Vector3DReader, angle float64) (blas64.General, error) {
// 	// TODO: do checks on axis Length, etc.
// 	c := math.Cos(angle)
// 	s := math.Sin(angle)
// 	c1 := 1 - c

// 	UU := blas64.General{
// 		Rows:   3,
// 		Cols:   3,
// 		Data:   []float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
// 		Stride: 3,
// 	}
// 	u := axis.ToBlasVector()
// 	blas64.Ger(c1, u, u, UU)

// 	Ux := get3DSkewSymmetricMatrix(axis, s)

// 	cI := blas64.General{
// 		Rows:   3,
// 		Cols:   3,
// 		Data:   []float64{c, 0, 0, 0, c, 0, 0, 0, c},
// 		Stride: 3,
// 	}

// 	return addMatrices(UU, Ux, cI)
// }
