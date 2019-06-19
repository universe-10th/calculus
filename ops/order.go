package ops

import (
	"github.com/universe-10th-calculus/sets"
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
