package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func Zero(set sets.Set) sets.Number {
	switch set {
	case sets.N, sets.N0, sets.Z:
		return int32(0)
	case sets.R:
		return float64(0)
	}
	return nil
}


func One(set sets.Set) sets.Number {
	switch set {
	case sets.N, sets.N0, sets.Z:
		return int32(1)
	case sets.R:
		return float64(1)
	}
	return nil
}


func IsZero(value sets.Number) bool {
	switch c := value.(type) {
	case float64:
		return c == 0
	case int32:
		return c == 0
	default:
		panic("cannot ask for zero a non-(int32, float64) value")
	}
}


func IsOne(value sets.Number) bool {
	switch c := value.(type) {
	case float64:
		return c == 1
	case int32:
		return c == 1
	default:
		panic("cannot ask for one a non-(int32, float64) value")
	}
}