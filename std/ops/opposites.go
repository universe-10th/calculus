package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func Neg(a sets.Number) sets.Number {
	switch va := a.(type) {
	case int32:
		return -va
	case float64:
		return -va
	default:
		panic("cannot negate a non-(int32, float64) value")
	}
}


func Inv(a sets.Number) sets.Number {
	switch va := a.(type) {
	case int32:
		return 1.0/float64(va)
	case float64:
		return 1.0/va
	default:
		panic("cannot invert a non-(int32, float64) value")
	}
}