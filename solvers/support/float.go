package support

import (
	"github.com/universe-10th-calculus/sets"
	"math/big"
	"errors"
)


func ForceFloat(value sets.Number) (*big.Float, error) {
	switch v := value.(type) {
	case *big.Float:
		return v, nil
	case *big.Int:
		return big.NewFloat(0).SetInt(v), nil
	case *big.Rat:
		return big.NewFloat(0).SetRat(v), nil
	default:
		return nil, errors.New("Only *big.(Float, Rat, Int) can be cast to *big.Float")
	}
}
