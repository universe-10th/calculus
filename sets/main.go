package sets

import "math/big"

type Number interface{}
type Set uint


const (
	N Set = iota
	N0
	Z
	Q
	R
	Invalid
)


var zeroInt = big.NewInt(0)
var sets    = []Set{N, N0, Z, Q, R}


func Valid(set Set) bool {
	for _, s := range sets {
		if set == s {
			return true
		}
	}
	return false
}


func setForUInt(value uint) Set {
	set := N
	if value == 0 {
		set = N0
	}
	return set
}


func Wrap(value interface{}) (Number, Set) {
	switch c := value.(type) {
	case int64:
		return big.NewInt(c), Z
	case int32:
		return big.NewInt(int64(c)), Z
	case int16:
		return big.NewInt(int64(c)), Z
	case int8:
		return big.NewInt(int64(c)), Z
	case uint32:
		return big.NewInt(int64(c)), setForUInt(uint(c))
	case uint16:
		return big.NewInt(int64(c)), setForUInt(uint(c))
	case uint8:
		return big.NewInt(int64(c)), setForUInt(uint(c))
	case float32:
		return big.NewFloat(float64(c)), R
	case float64:
		return big.NewFloat(float64(c)), R
	default:
		panic("unsupported number type. Only int64, int32, int16, int8, uint32, uint16, uint8, float64, float43 " +
			  "are supported")
	}
}


func Clone(value interface{}) Number {
	switch c := value.(type) {
	case *big.Float:
		return big.NewFloat(0).Set(c)
	case *big.Rat:
		return big.NewRat(0, 1).Set(c)
	case *big.Int:
		return big.NewInt(0).Set(c)
	default:
		// No cloning will occur for other values
		return value
	}
}


func broader(a, b Set) Set {
	if !Valid(b) {
		return Invalid
	}
	switch a {
	case N:
		return b
	case N0:
		switch b {
		case N:
			return N0
		default:
			return b
		}
	case Z:
		switch b {
		case Q, R:
			return b
		default:
			return Z
		}
	case Q:
		switch b {
		case R:
			return R
		default:
			return Q
		}
	case R:
		return R
	}
	return Invalid
}


func BroaderAll(sets ...Set) Set {
	if sets == nil {
		return Invalid
	}
	x := sets[0]
	for _, s := range sets[1:] {
		x = broader(x, s)
	}
	return x
}


func closest(a Number) Set {
	switch va := a.(type) {
	case *big.Int:
		switch va.Cmp(zeroInt) {
		case 0:
			return N0
		case 1:
			return N
		default:
			return Z
		}
	case *big.Rat:
		return Q
	case *big.Float:
		return R
	default:
		panic("cannot get the closest set of a non-*big.(Int, Float, Rat) value")
	}
}


func ClosestAll(a ...Number) []Set {
	sets := make([]Set, len(a))
	for index, value := range a {
		sets[index] = closest(value)
	}
	return sets
}


func upCast(a Number, s Set) Number {
	switch va := a.(type) {
	case *big.Int:
		switch s {
		case Q:
			return big.NewRat(0, 1).SetInt(va)
		case R:
			return big.NewFloat(0).SetInt(va)
		}
	case *big.Rat:
		if s == R {
			return big.NewFloat(0).SetRat(va)
		}
	case *big.Float:
	default:
		panic("cannot up-cast a non-*big.(Int, Float, Rat) value")
	}
	return a
}


func UpCastTo(toSet Set, a ...Number) []Number {
	result := make([]Number, len(a))
	for index, value := range a {
		result[index] = upCast(value, toSet)
	}
	return result
}


func UpCast(a ...Number) []Number {
	if a == nil {
		return nil
	}
	toSet := BroaderAll(ClosestAll(a...)...)
	return UpCastTo(toSet, a...)
}