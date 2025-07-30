package util

import "math"

const (
	// Epsilon normal precision floating point
	Epsilon = 0.00000001
	// LowEpsilon low precision floating point
	LowEpsilon = 0.01
)

// AC floating point precision function
var AC Accuracy = func() float64 { return LowEpsilon }

// Accuracy precision type
type Accuracy func() float64

// Equal checks if two floating point numbers are equal
func (ac Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < LowEpsilon
}

// Greater checks if a is greater than b
func (ac Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > ac()
}

// Smaller checks if a is smaller than b
func (ac Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > ac()
}

// GreaterOrEqual checks if a is greater than or equal to b
func (ac Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < ac()
}

// SmallerOrEqual checks if a is smaller than or equal to b
func (ac Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < ac()
}
