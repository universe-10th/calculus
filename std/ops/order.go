package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func IsNegative(value sets.Number) bool {
	switch c := value.(type) {
	case float64:
		return c < 0
	case int32:
		return c < 0
	default:
		panic("cannot ask for negativity a non-(int32, float64) value")
	}
}


func IsPositive(value sets.Number) bool {
	switch c := value.(type) {
	case float64:
		return c > 0
	case int32:
		return c > 0
	default:
		panic("cannot ask for positivity a non-(int32, float64) value")
	}
}


func OpposeSigns(a, b sets.Number) bool {
	cast := sets.UpCast(a, b)
	a = cast[0]
	b = cast[1]

	switch va := a.(type) {
	case int32:
		return va * b.(int32) < 0
	case float64:
		return va * b.(float64) < 0
	default:
		panic("cannot compare non-(int32, float64) sign values")
	}
}