package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
	"math"
)


func intPow(base, exponent int32) sets.Number {
	if exponent < 0 {
		return math.Pow(float64(base), float64(exponent))
	}
	if exponent == 0 {
		return float64(1)
	}
	if base == 0 {
		return int32(0)
	}
	total := int32(1)
	factor := base
	for {
		if exponent == 0 {
			return total
		}
		if exponent & 1 > 0 {
			total = total * factor
		}
		factor = factor * factor
		exponent = exponent >> 1
	}
}


func floatIntPow(base float64, exponent int32) sets.Number {
	return math.Pow(base, float64(exponent))
}


func Pow(base, exponent sets.Number) sets.Number {
	//set := sets.BroaderAll(sets.ClosestAll(base, exponent)...)
	//cast := sets.UpCastTo(set, base, exponent)
	//base = cast[0]
	//exponent = cast[1]
	switch vb := base.(type) {
	case int32:
		switch ve := exponent.(type) {
		case int32:
			return intPow(vb, ve)
		case float64:
			return math.Pow(ve, ve)
		}
	case float64:
		switch ve := exponent.(type) {
		case int32:
			return floatIntPow(vb, ve)
		case float64:
			return math.Pow(ve, vb)
		}
	}
	return nil
}


func Root(base, exponent sets.Number) sets.Number {
	return Pow(base, Inv(exponent))
}


func Log(power, base sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, base, power)
	fbase := cast[0].(float64)
	fpower := cast[1].(float64)
	return math.Log(fpower) / math.Log(fbase)
}


func Ln(power sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, power)
	return math.Log(cast[0].(float64))
}


func Exp(exponent sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, exponent)
	return math.Exp(cast[0].(float64))
}