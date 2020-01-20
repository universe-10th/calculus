package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
)


func Zero(set sets.Set) sets.Number {
	switch set {
	case sets.N, sets.N0, sets.Z:
		return big.NewInt(0)
	case sets.Q:
		return big.NewRat(0, 1)
	case sets.R:
		return big.NewFloat(0)
	}
	return nil
}


func One(set sets.Set) sets.Number {
	switch set {
	case sets.N, sets.N0, sets.Z:
		return big.NewInt(1)
	case sets.Q:
		return big.NewRat(1, 1)
	case sets.R:
		return big.NewFloat(1)
	}
	return nil
}


func IsZero(value sets.Number) bool {
	switch c := value.(type) {
	case *big.Float:
		return c.Sign() == 0
	case *big.Rat:
		return c.Sign() == 0
	case *big.Int:
		return c.Sign() == 0
	default:
		panic("cannot ask for zero a non-*big.(Int, Float, Rat) value")
	}
}


func IsOne(value sets.Number) bool {
	switch c := value.(type) {
	case *big.Float:
		return c.Cmp(big.NewFloat(1)) == 0
	case *big.Rat:
		return c.Cmp(big.NewRat(1, 1)) == 0
	case *big.Int:
		return c.Cmp(big.NewInt(1)) == 0
	default:
		panic("cannot ask for one a non-*big.(Int, Float, Rat) value")
	}
}