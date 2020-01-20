// Sets package contains utils related to numeric sets.
package sets

// Number is a marker type that tells the given values will be treated as numbers.
// Hence, few types are actually allowed here.
type Number interface{}


// Set stands for numbers sets like real, integer.
type Set uint


const (
	// N stands for natural numbers.
	N Set = iota
	// N0 stands for natural numbers and zero.
	N0
	// Z stands for integer numbers.
	Z
	// R stands for real numbers.
	R
	// Invalid stands for an invalid set.
	Invalid
)


var zeroInt = 0
var sets    = []Set{N, N0, Z, R}


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


func setForInt(value int32) Set {
	if value < 0 {
		return Z
	} else if value == 0 {
		return N0
	} else {
		return N
	}
}


// Wrap makes a {int32, float64} value out of a numeric primitive type.
func Wrap(value interface{}) (Number, Set) {
	switch c := value.(type) {
	case int:
		return int32(c), Z
	case int32:
		return int32(c), Z
	case int16:
		return int32(c), Z
	case int8:
		return int32(c), Z
	case uint16:
		return int32(c), setForUInt(uint(c))
	case uint8:
		return int64(c), setForUInt(uint(c))
	case float32:
		return float64(c), R
	case float64:
		return float64(c), R
	default:
		panic("unsupported number type. Only int32, int16, int8, uint16, uint8, float64, and float32 " +
			  "are supported")
	}
}


// This function is a no-op with std types.
func Clone(value interface{}) Number {
	return value
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
		case R:
			return b
		default:
			return Z
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
	case int:
		return setForInt(int32(va))
	case uint8:
		return setForInt(int32(va))
	case uint16:
		return setForInt(int32(va))
	case int8:
		return setForInt(int32(va))
	case int16:
		return setForInt(int32(va))
	case int32:
		return setForInt(int32(va))
	case float32, float64:
		return R
	default:
		panic("cannot get the closest set of a non-(int32, float64) value")
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
	case float32, float64:
		return set == R
	case uint8:
		return set == Z || set == N0 || (va > 0 && set == N) || set == R
	case uint16:
		return set == Z || set == N0 || (va > 0 && set == N) || set == R
	case int8:
		return set == Z || (va >= 0 && set == N0) || (va > 0 && set == N) || set == R
	case int16:
		return set == Z || (va >= 0 && set == N0) || (va > 0 && set == N) || set == R
	case int32:
		return set == Z || (va >= 0 && set == N0) || (va > 0 && set == N) || set == R
	default:
		return false
	}
}


func UpCastOneTo(a Number, s Set) Number {
	switch va := a.(type) {
	case int32:
		if s == R {
			return float64(va)
		}
	case float64:
	default:
		panic("cannot up-cast a non-(int32, float64) value")
	}
	return a
}


// UpCastTo ensures all the given numbers are converted to the appropriate set.
// Depending on the given set, the given numbers will be converted, one by one,
// to the type belonging to the set... unless they are of that type, or broader.
// E.g. a float64 will not be cast to int32 if the set is Z, but will
// remain a float64.
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