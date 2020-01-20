package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func inc(a, b sets.Number) sets.Number {
	switch va := a.(type) {
	case int32:
		return b.(int32) + va
	case float64:
		return b.(float64) + va
	default:
		return nil
	}
}


func Add(terms ...sets.Number) sets.Number {
	if len(terms) == 0 {
		return int32(0)
	} else if len(terms) == 1 {
		return terms[0]
	}
	set := sets.BroaderAll(sets.ClosestAll(terms...)...)
	current := Zero(set)
	terms = sets.UpCastTo(set, terms...)
	for _, term := range terms {
		current = inc(current, term)
	}
	return current
}


func Sub(minuend sets.Number, subtrahends ...sets.Number) sets.Number {
	if subtrahends == nil {
		return minuend
	}
	subtrahend := Add(subtrahends...)
	set := sets.BroaderAll(sets.ClosestAll(minuend, subtrahend)...)
	cast := sets.UpCastTo(set, minuend, subtrahend)
	minuend = cast[0]
	subtrahend = cast[1]
	switch vm := minuend.(type) {
	case int32:
		return vm - subtrahend.(int32)
	case float64:
		return vm - subtrahend.(float64)
	}
	return nil
}