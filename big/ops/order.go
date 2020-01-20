package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
)


func Cmp(a, b sets.Number) int {
	cast := sets.UpCast(a, b)
	a = cast[0]
	b = cast[1]
	switch va := a.(type) {
	case *big.Int:
		return va.Cmp(b.(*big.Int))
	case *big.Rat:
		return va.Cmp(b.(*big.Rat))
	case *big.Float:
		return va.Cmp(b.(*big.Float))
	default:
		panic("cannot compare non-*big.(Int, Float, Rat) values")
	}
}


func IsNegative(value sets.Number) bool {
	switch c := value.(type) {
	case *big.Float:
		return c.Cmp(big.NewFloat(0)) == -1
	case *big.Rat:
		return c.Cmp(big.NewRat(0, 1)) == -1
	case *big.Int:
		return c.Cmp(big.NewInt(0)) == -1
	default:
		panic("cannot ask for negativeness a non-*big.(Int, Float, Rat) value")
	}
}


func IsPositive(value sets.Number) bool {
	switch c := value.(type) {
	case *big.Float:
		return c.Cmp(big.NewFloat(0)) == 1
	case *big.Rat:
		return c.Cmp(big.NewRat(0, 1)) == 1
	case *big.Int:
		return c.Cmp(big.NewInt(0)) == 1
	default:
		panic("cannot ask for negativeness a non-*big.(Int, Float, Rat) value")
	}
}


func OpposeSigns(a, b sets.Number) bool {
	cast := sets.UpCast(a, b)
	a = cast[0]
	b = cast[1]

	switch va := a.(type) {
	case *big.Int:
		return va.Sign() == -b.(*big.Int).Sign()
	case *big.Rat:
		return va.Sign() == -b.(*big.Rat).Sign()
	case *big.Float:
		return va.Sign() == -b.(*big.Float).Sign()
	default:
		panic("cannot compare non-*big.(Int, Float, Rat) values")
	}
}