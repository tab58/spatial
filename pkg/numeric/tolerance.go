package numeric

// IsInvalidTolerance returns true if the tolerance value is invalid, false if valid.
func IsInvalidTolerance(tol float64) bool {
	return tol < 0
}
