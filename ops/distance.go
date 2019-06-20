package ops

import (
	"github.com/universe-10th-calculus/sets"
	"math/big"
)


func NearBy(a, b sets.Number, epsilon *big.Float) bool {
	cast := sets.UpCast(a, b, epsilon)
	return Cmp(Abs(Sub(cast[0], cast[1])), cast[3]) <= 0
}
