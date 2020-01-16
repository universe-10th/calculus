package ops

import (
	"github.com/universe-10th/calculus/v2/sets"
	"math/big"
)


func mul(a, b, zero sets.Number) bool {
	switch va := a.(type) {
	case *big.Int:
		if va.Cmp(zero.(*big.Int)) == 0 {
			return true
		}
		va.Mul(va, b.(*big.Int))
	case *big.Rat:
		if va.Cmp(zero.(*big.Rat)) == 0 {
			return true
		}
		va.Mul(va, b.(*big.Rat))
	case *big.Float:
		if va.Cmp(zero.(*big.Float)) == 0 {
			return true
		}
		va.Mul(va, b.(*big.Float))
	}
	return false
}


func Mul(factors ...sets.Number) sets.Number {
	if len(factors) == 0 {
		return nil
	} else if len(factors) == 1 {
		return factors[0]
	}
	set := sets.BroaderAll(sets.ClosestAll(factors...)...)
	zero := Zero(set)
	current := One(set)
	factors = sets.UpCastTo(set, factors...)
	for _, term := range factors {
		if mul(current, term, zero) {
			return zero
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
	case *big.Int:
		return big.NewRat(0, 1).SetFrac(vm, divider.(*big.Int))
	case *big.Rat:
		return vm.Quo(vm, divider.(*big.Rat))
	case *big.Float:
		return vm.Quo(vm, divider.(*big.Float))
	}
	return nil
}
