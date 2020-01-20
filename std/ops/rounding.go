package ops

import (
	"github.com/universe-10th/calculus/v2/std/sets"
	"math"
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


func roundFloat(float float64, roundType RoundType) int32 {
	switch roundType {
	case Ceil:
		return int32(math.Ceil(float))
	case Floor:
		return int32(math.Floor(float))
	case Inward:
		if float > 0 {
			return int32(math.Floor(float))
		} else {
			return int32(math.Ceil(float))
		}
	case Outward:
		return int32(math.Round(float))
	default:
		panic("cannot round the number: invalid round type")
	}
}


// Rounds a number according to certain modes:
// - int32 numbers are returned unchanged.
// - float64 numbers are rounded according to the rounding mode.
// The rounding mode works as follows:
// - Inward rounds toward 0.
// - Outward rounds away from 0, if the float is not already an integer.
// - Ceil rounds toward +inf, if the float is not already an integer.
// - Floor rounds toward -inf, if the float is not already an integer.
func Round(number sets.Number, roundType RoundType) int32 {
	switch vn := number.(type) {
	case int32:
		return vn
	case float64:
		return roundFloat(vn, roundType)
	default:
		panic("cannot round a non-(int32, float64) value")
	}
}


// Takes the number, rounds it INWARD, and then subtracts the result
// from the original number. Frac(1.35) will be 0.35,
// and Frac(-1.35) will be -0.35.
func Frac(number sets.Number) sets.Number {
	return Sub(number, Round(number, Inward))
}


// Takes the number, rounds it INWARD, and then subtracts the result
// from the original number. Both the rounded number and the result
// of the subtraction (which will be the fractional part) will be
// returned.
// Note: This function will never be used in Expression classes, at
// least while we don't support multivalued expressions.
func Split(number sets.Number) (sets.Number, sets.Number) {
	rounded := Round(number, Inward)
	return rounded, Sub(number, rounded)
}