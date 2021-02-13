package geometry

// MatrixReader is a read-only interface for a matrix.
type MatrixReader interface {
	Rows() uint
	Cols() uint
	ElementAt(i uint) float64
}

// MatrixWriter is a write-only interface for a matrix.
type MatrixWriter interface {
	SetElementAt(i uint, value float64) error
}

// Matrix2D is a representation of a 2x2 matrix.
type Matrix2D struct {
	elements [4]float64
}
