package ops

import (
	"github.com/universe-10th/calculus/v2/sets"
	"math/big"
)


func inc(a, b sets.Number) {
	switch va := a.(type) {
	case *big.Int:
		va.Add(va, b.(*big.Int))
	case *big.Rat:
		va.Add(va, b.(*big.Rat))
	case *big.Float:
		va.Add(va, b.(*big.Float))
	}
}


func Add(terms ...sets.Number) sets.Number {
	if len(terms) == 0 {
		return nil
	} else if len(terms) == 1 {
		return terms[0]
	}
	set := sets.BroaderAll(sets.ClosestAll(terms...)...)
	current := Zero(set)
	terms = sets.UpCastTo(set, terms...)
	for _, term := range terms {
		inc(current, term)
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
	case *big.Int:
		return vm.Sub(vm, subtrahend.(*big.Int))
	case *big.Rat:
		return vm.Sub(vm, subtrahend.(*big.Rat))
	case *big.Float:
		return vm.Sub(vm, subtrahend.(*big.Float))
	}
	return nil
}