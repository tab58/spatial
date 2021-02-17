package errors

import (
	e "errors"
)

// ErrOverflow expresses that a computation has resulted in numeric overflow.
var ErrOverflow = e.New("numeric overflow")

// ErrUnderflow expresses that a computation has resulted in numeric underflow.
var ErrUnderflow = e.New("numeric underflow")

// ErrNaN expresses that a computation has resulted in an invalid results (NaN).
var ErrNaN = e.New("result is NaN")

// ErrVectorZeroLength expresses that a vector is zero length.
var ErrVectorZeroLength = e.New("length of vector is zero")

// ErrInfinity expresses that a number is infinite.
var ErrInfinity = e.New("number is infinite")

// ErrDivideByZero expresses a division by zero.
var ErrDivideByZero = e.New("division by zero")

// ErrEmptyArray expresses that a specific array is of length 0.
var ErrEmptyArray = e.New("array is of length 0")

// ErrMatrixOutOfRange expresses that the index of a matrix is out of range.
var ErrMatrixOutOfRange = e.New("matrix indices are out of range")

// ErrMatrixDims expresses that the matrix dimensions for a specific operation don't match.
var ErrMatrixDims = e.New("matrix dimensions do not match")

// ErrInvalidArgument expresses that one of the arguments supplied is unexpectedly invalid.
var ErrInvalidArgument = e.New("argument is invalid")

// ErrInvalidTol expresses that a tolerance value is invalid.
var ErrInvalidTol = e.New("invalid value for tolerance; must be nonnegative")
