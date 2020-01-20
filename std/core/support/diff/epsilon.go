package diff

import (
	"math"
)

// Epsilon takes an integer precision, and creates an epsilon value with that precision.
// Precision 0 will become precision 1, and precision -X will become precision X. A sample
// precision 4 would return 0.0001, while a precision of 16: 0.0000000000000001. Use this
// function few times and keep the result the safest possible.
func Epsilon(precision int64) float64 {
	if precision == 0 {
		precision = 1
	} else if precision < 0 {
		precision = -precision
	}
	return math.Pow(10, float64(-precision))
}
