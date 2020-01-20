package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
	"math"
)

var float64limit = big.NewFloat(math.MaxFloat64)

func trigf(f func(float64) float64) func(sets.Number) sets.Number {
	return func(a sets.Number) sets.Number {
		ra := sets.UpCastTo(sets.R, a)[0].(*big.Float)
		if big.NewFloat(0).Abs(ra).Cmp(float64limit) > 0 {
			panic("the float value exceeds 64 bits - trigonometric functions use that precision")
		} else {
			value, _ := ra.Float64()
			return big.NewFloat(f(value))
		}
	}
}


var Sin = trigf(math.Sin)
var Cos = trigf(math.Cos)
var Tan = trigf(math.Tan)
