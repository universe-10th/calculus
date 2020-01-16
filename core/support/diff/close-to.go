package diff

import (
	"math/big"
)


// CloseTo tells whether two values a and b are close: their distance is less than / equal a given threshold.
// Since this method is intended as a support method for solvers, and due to the way big.Float works, two
// extra arguments must be given: diff, dist. They serve as intermediate steps, actually: You'll have few
// diff/dist variables in your execution frame, and you'll not face race conditions most of the times.
func CloseTo(a, b, diff, dist, epsilon *big.Float) bool {
	return dist.Abs(diff.Sub(a, b)).Cmp(epsilon) <= 0
}