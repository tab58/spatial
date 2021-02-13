package geometry

import "errors"

// ErrInfinity expresses that a number is infinite.
var ErrInfinity = errors.New("number is infinite")

// ErrDivideByZero expresses a division by zero.
var ErrDivideByZero = errors.New("division by zero")

// ErrEmptyArray expresses that a specific array is of length 0.
var ErrEmptyArray = errors.New("array is of length 0")

// ErrMatrixDims expresses that the matrix dimensions for a specific operation don't match.
var ErrMatrixDims = errors.New("matrix dimensions do not match")

// ErrInvalidArgument expresses that one of the arguments supplied is unexpectedly invalid.
var ErrInvalidArgument = errors.New("argument is invalid")

// ErrInvalidTol expresses that a tolerance value is invalid.
var ErrInvalidTol = errors.New("invalid value for tolerance; must be nonnegative")
