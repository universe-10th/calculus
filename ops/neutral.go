package ops

import (
	"github.com/universe-10th-calculus/sets"
	"math/big"
)


func Zero(set sets.Set) sets.Number {
	switch set {
	case sets.N0, sets.Z:
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
	case sets.N0, sets.Z:
		return big.NewInt(1)
	case sets.Q:
		return big.NewRat(1, 1)
	case sets.R:
		return big.NewFloat(1)
	}
	return nil
}
