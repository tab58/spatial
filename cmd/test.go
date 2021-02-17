package main

import (
	"fmt"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

func main() {
	A := blas64.General{
		Rows:   2,
		Cols:   2,
		Data:   []float64{-2, 1, 0, 4},
		Stride: 2,
	}

	B := blas64.General{
		Rows:   2,
		Cols:   2,
		Data:   []float64{6, 5, -7, 1},
		Stride: 2,
	}

	C := blas64.General{
		Rows:   2,
		Cols:   2,
		Data:   []float64{0, 0, 0, 0},
		Stride: 2,
	}

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, A, B, 0, C)

	fmt.Printf("C: %#v\n", C)
	fmt.Printf("B: %#v\n", B)

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, A, B, 0, B)

	fmt.Printf("new B: %#v\n", B)
}
