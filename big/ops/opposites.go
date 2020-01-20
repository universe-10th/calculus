package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
)


func Neg(a sets.Number) sets.Number {
	switch va := a.(type) {
	case *big.Int:
		return big.NewInt(0).Neg(va)
	case *big.Rat:
		return big.NewRat(0, 1).Neg(va)
	case *big.Float:
		return big.NewFloat(0).Neg(va)
	default:
		panic("cannot negate a non-*big.(Int, Float, Rat) value")
	}
}


func Inv(a sets.Number) sets.Number {
	switch va := a.(type) {
	case *big.Int:
		return big.NewRat(0, 1).SetFrac(oneInt, va)
	case *big.Rat:
		return big.NewRat(0, 1).Inv(va)
	case *big.Float:
		return big.NewFloat(0).Quo(oneFloat, va)
	default:
		panic("cannot invert a non-*big.(Int, Float, Rat) value")
	}
}