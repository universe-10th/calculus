package ops

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/ALTree/bigfloat"
	"math/big"
)


func intPow(base, exponent *big.Int) *big.Int {
	if exponent.Sign() == 0 {
		return big.NewInt(1)
	}
	if base.Sign() == 0 {
		return big.NewInt(0)
	}
	one := big.NewInt(1)
	flag := big.NewInt(0)
	total := one
	factor := base
	for {
		if exponent.Cmp(one) == 0 {
			return total
		}
		if flag.And(exponent, one).Sign() > 0 {
			total.Mul(total, factor)
		}
		factor.Mul(factor, factor)
		exponent.Rsh(exponent, 1)
	}
	return total
}


func Pow(base, exponent sets.Number) sets.Number {
	set := sets.BroaderAll(sets.ClosestAll(base, exponent)...)
	cast := sets.UpCastTo(set, base, exponent)
	base = cast[0]
	exponent = cast[1]
	switch vb := base.(type) {
	case *big.Int:
		ve := exponent.(*big.Int)
		abse := big.NewInt(0).Abs(ve)
		signe := ve.Sign()
		if signe == 0 && vb.Sign() == 0 {
			panic("attempting to calculate integers 0 elevated to 0")
		} else if signe == -1 {
			return intPow(vb, abse)
		} else if signe == 1 {
			return big.NewRat(0, 1).SetFrac(big.NewInt(1), intPow(vb, abse))
		}
	case *big.Rat:
		return bigfloat.Pow(big.NewFloat(0).SetRat(vb), big.NewFloat(0).SetRat(exponent.(*big.Rat)))
	case *big.Float:
		return bigfloat.Pow(vb, exponent.(*big.Float))
	}
	return nil
}


func Root(base, exponent sets.Number) sets.Number {
	return Pow(base, Inv(exponent))
}


func Log(power, base sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, base, power)
	fbase := cast[0].(*big.Float)
	fpower := cast[1].(*big.Float)
	return big.NewFloat(0).Quo(bigfloat.Log(fpower), bigfloat.Log(fbase))
}


func Ln(power sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, power)
	return bigfloat.Log(cast[0].(*big.Float))
}


func Exp(exponent sets.Number) sets.Number {
	cast := sets.UpCastTo(sets.R, exponent)
	return bigfloat.Exp(cast[0].(*big.Float))
}