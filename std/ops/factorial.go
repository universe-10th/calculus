package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
)


func Factorial(a sets.Number) sets.Number {
	if !sets.BelongsTo(a, sets.N0) {
		panic("only integers >= 0 will have a factorial")
	}
	v := a.(int32)
	if v <= 12 {
		total := int32(1)
		index := int32(1)
		for index = 1; index <= v; index++  {
			total *= index
		}
		return total
	} else {
		panic("the number is way too big to calculate a factorial")
	}
}
