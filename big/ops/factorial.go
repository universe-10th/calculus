package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
)


func Factorial(a sets.Number) sets.Number {
	if !sets.BelongsTo(a, sets.N0) {
		panic("only integers >= 0 will have a factorial")
	}
	v := a.(*big.Int)
	if v.IsInt64() {
		return big.NewInt(0).MulRange(1, v.Int64())
	} else {
		panic("the number is way too big to calculate a factorial")
	}
}
