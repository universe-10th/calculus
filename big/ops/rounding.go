package ops

import (
	"github.com/universe-10th/calculus/v2/big/sets"
	"math/big"
)


type RoundType int
const (
	// If non-integer, rounds to nearest integer "to the right"
	Ceil RoundType = iota
	// If non-integer, rounds to nearest integer "to the left"
	Floor
	// If non-integer, rounds to the closest integer lower in absolute value
	Inward
	// If non-integer, rounds to the closest integer higher in absolute value
	Outward
)


func roundFloat(float *big.Float, roundType RoundType) *big.Int {
	roundedToZero, _ := float.Int(big.NewInt(0))
	if roundedToZero == nil {
		return nil
	}
	switch roundType {
	case Ceil:
		sign := float.Sign()
		if float.IsInt() || sign < 0 {
			return roundedToZero
		} else {
			return roundedToZero.Add(roundedToZero, big.NewInt(1))
		}
	case Floor:
		sign := float.Sign()
		if float.IsInt() || sign > 0 {
			return roundedToZero
		} else {
			return roundedToZero.Sub(roundedToZero, big.NewInt(1))
		}
	case Inward:
		return roundedToZero
	case Outward:
		if float.IsInt() {
			return roundedToZero
		} else {
			return roundedToZero.Add(roundedToZero, big.NewInt(int64(1 * float.Sign())))
		}
	default:
		panic("cannot round the number: invalid round type")
	}
}


// Rounds a big number according to certain modes:
// - *big.Int numbers are returned unchanged.
// - *big.Rat numbers are converted to float, and then rounded according to the rounding mode.
// - *big.Float numbers are rounded according to the rounding mode.
// The rounding mode works as follows:
// - Inward rounds toward 0.
// - Outward rounds away from 0, if the float is not already an integer.
// - Ceil rounds toward +inf, if the float is not already an integer.
// - Floor rounds toward -inf, if the float is not already an integer.
func Round(number sets.Number, roundType RoundType) *big.Int {
	switch vn := number.(type) {
	case *big.Int:
		return vn
	case *big.Rat:
		return roundFloat(big.NewFloat(0).SetRat(vn), roundType)
	case *big.Float:
		return roundFloat(vn, roundType)
	default:
		panic("cannot round a non-*big.(Int, Float, Rat) value")
	}
}


// Takes the number, rounds it INWARD, and then subtracts the result
// from the original number. Frac(big.NewFloat(1.35)) will be 0.35,
// and Frac(big.NewFloat(-1.35)) will be -0.35.
func Frac(number sets.Number) sets.Number {
	if rounded := Round(number, Inward); rounded == nil {
		return nil
	} else {
		return Sub(number, Round(number, Inward))
	}
}


// Takes the number, rounds it INWARD, and then subtracts the result
// from the original number. Both the rounded number and the result
// of the subtraction (which will be the fractional part) will be
// returned.
// Note: This function will never be used in Expression classes, at
// least while we don't support multivalued expressions.
func Split(number sets.Number) (sets.Number, sets.Number) {
	if rounded := Round(number, Inward); rounded != nil {
		return rounded, Sub(number, rounded)
	} else {
		return nil, nil
	}
}