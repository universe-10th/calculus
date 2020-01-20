package diff

import "math"

// CloseTo tells whether two values a and b are close: their distance is less than / equal a given threshold.
func CloseTo(a, b, epsilon float64) bool {
	return math.Abs(a - b) <= epsilon
}