package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
	"math"
)

func trigf(f func(float64) float64) func(sets.Number) sets.Number {
	return func(a sets.Number) sets.Number {
		return f(sets.UpCastTo(sets.R, a)[0].(float64))
	}
}


var Sin = trigf(math.Sin)
var Cos = trigf(math.Cos)
var Tan = trigf(math.Tan)
