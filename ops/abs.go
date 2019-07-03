package ops

import (
	"github.com/universe-10th/calculus/sets"
	"math/big"
)


func Abs(value sets.Number) sets.Number {
	switch c := value.(type) {
	case *big.Float:
		return big.NewFloat(0).Abs(c)
	case *big.Rat:
		return big.NewRat(0, 1).Abs(c)
	case *big.Int:
		return big.NewInt(0).Abs(c)
	default:
		panic("cannot get absolute value of a non-*big.(Int, Float, Rat) value")
	}
}

