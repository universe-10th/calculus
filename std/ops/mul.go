package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func mul(a, b sets.Number) sets.Number {
	switch va := a.(type) {
	case int32:
		return va * b.(int32)
	case float64:
		return va * b.(float64)
	}
	return 1.0
}


func Mul(factors ...sets.Number) sets.Number {
	if len(factors) == 0 {
		return int32(1)
	} else if len(factors) == 1 {
		return factors[0]
	}
	set := sets.BroaderAll(sets.ClosestAll(factors...)...)
	current := One(set)
	factors = sets.UpCastTo(set, factors...)
	for _, factor := range factors {
		current = mul(current, factor)
		if IsZero(current) {
			return current
		}
	}
	return current
}


func Div(dividend sets.Number, dividers ...sets.Number) sets.Number {
	if dividers == nil {
		return dividend
	}
	divider := Add(dividers...)
	set := sets.BroaderAll(sets.ClosestAll(dividend, divider)...)
	cast := sets.UpCastTo(set, dividend, divider)
	dividend = cast[0]
	divider = cast[1]
	switch vm := dividend.(type) {
	case int32:
		return vm / divider.(int32)
	case float64:
		return vm / divider.(float64)
	}
	return nil
}
