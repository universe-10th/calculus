// Sets package contains utils related to numeric sets.
package sets

import "math/big"


// Number is a marker type that tells the given values will be treated as numbers.
// Hence, few types are actually allowed here.
type Number interface{}


// Set stands for numbers sets like real, integer, rational[, no complex supported yet due to the lack of *big.Complex].
type Set uint


const (
	// N stands for natural numbers.
	N Set = iota
	// N0 stands for natural numbers and zero.
	N0
	// Z stands for integer numbers.
	Z
	// Q stands for rational numbers.
	Q
	// R stands for real numbers.
	R
	// Invalid stands for an invalid set.
	Invalid
)


var zeroInt = big.NewInt(0)
var sets    = []Set{N, N0, Z, Q, R}


// Valid tests whether the given set is among N, N0, Z, Q, R.
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


// Wrap makes a *big.(Int, Rat, Float) value out of a numeric primitive type.
// If given a *big.(Int, Rat, Float) object, it is returned as-is.
func Wrap(value interface{}) (Number, Set) {
	switch c := value.(type) {
	case int:
		return big.NewInt(int64(c)), Z
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
	case *big.Int:
		return c, closest(c)
	case *big.Rat:
		return c, Q
	case *big.Float:
		return c, R
	default:
		panic("unsupported number type. Only int64, int32, int16, int8, uint32, uint16, uint8, float64, float32 " +
			  "and math/big types (they are a no-op) are supported")
	}
}


// Clone makes a copy of the given *big.(Int, Rat, Float) object.
// Other values are returned as-is.
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


// BroaderAll computes the minimal set that includes all the given sets.
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


// ClosestAll returns, for the given arguments, the minimal sets that include them each.
func ClosestAll(a ...Number) []Set {
	sets := make([]Set, len(a))
	for index, value := range a {
		sets[index] = closest(value)
	}
	return sets
}


// BelongsTo returns whether a given number is a member of a given set.
func BelongsTo(a Number, set Set) bool {
	switch va := a.(type) {
	case *big.Float:
		return set == R
	case *big.Rat:
		return set == R || set == Q
	case *big.Int:
		return set == Z || (va.Sign() >= 0 && set == N0) || (va.Sign() > 0 && set == N) || set == Q || set == R
	default:
		return false
	}
}


func UpCastOneTo(a Number, s Set) Number {
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


// UpCastTo ensures all the given numbers are converted to the appropriate set.
// Depending on the given set, the given numbers will be converted, one by one,
// to the type belonging to the set... unless they are of that type, or broader.
// E.g. a *big.Float will not be cast to *big.Int if the set is Z, but will
// remain a *big.Float.
func UpCastTo(toSet Set, a ...Number) []Number {
	result := make([]Number, len(a))
	for index, value := range a {
		result[index] = UpCastOneTo(value, toSet)
	}
	return result
}


// UpCast casts all the given numbers to their closest common set.
func UpCast(a ...Number) []Number {
	if a == nil {
		return nil
	}
	toSet := BroaderAll(ClosestAll(a...)...)
	return UpCastTo(toSet, a...)
}