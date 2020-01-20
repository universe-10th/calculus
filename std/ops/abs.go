package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
	"math"
)


func Abs(value sets.Number) sets.Number {
	switch c := value.(type) {
	case int32:
		if c >= 0 {
			return c
		} else {
			return -c
		}
	case float64:
		return math.Abs(c)
	default:
		panic("cannot get absolute value of a non-(int32, float64) value")
	}
}

