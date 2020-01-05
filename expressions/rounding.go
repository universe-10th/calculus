package expressions

import (
	"fmt"
	"math/big"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/errors"
	"github.com/universe-10th/calculus/ops"
)

type Round struct {
	arg Expression
	roundType ops.RoundType
}


func (round Round) Curry(arguments Arguments) (Expression, error) {
	if simplified, err := round.arg.Simplify(); err != nil {
		return nil, err
	} else {
		return Round{simplified, round.roundType}.Simplify()
	}
}


func (round Round) CollectVariables(variables Variables) {
	round.arg.CollectVariables(variables)
}


func (round Round) IsConstant(wrt Variable) bool {
	return round.arg.IsConstant(wrt)
}


func (round Round) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := round.arg.Evaluate(args); err != nil{
		return nil, err
	} else if result := ops.Round(result, round.roundType); result != nil {
		return result, nil
	} else {
		return nil, errors.InfiniteCannotBeRounded
	}
}


func (round Round) Derivative(wrt Variable) (Expression, error) {
	// Rounding is discontinuous on integers,
	// but continuous otherwise, and 0.
	return DefectiveOnInt{
		round.arg,
		big.NewInt(0),
	}, nil
}


func (round Round) Simplify() (Expression, error) {
	if simplified, err := round.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := (Round{num, round.roundType}).Evaluate(Arguments{}); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return Round{simplified,round.roundType}, nil
	}
}


func (round Round) String() string {
	return fmt.Sprintf("Round(%s, %v)", round.arg, round.roundType)
}


func (round Round) IsSelfContained() bool {
	return true
}


type Frac struct {
	arg Expression
}


func (frac Frac) Curry(arguments Arguments) (Expression, error) {
	if simplified, err := frac.arg.Simplify(); err != nil {
		return nil, err
	} else {
		return Frac{simplified}.Simplify()
	}
}


func (frac Frac) CollectVariables(variables Variables) {
	frac.arg.CollectVariables(variables)
}


func (frac Frac) IsConstant(wrt Variable) bool {
	return frac.arg.IsConstant(wrt)
}


func (frac Frac) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := frac.arg.Evaluate(args); err != nil{
		return nil, err
	} else {
		return ops.Frac(result), nil
	}
}


func (frac Frac) Derivative(wrt Variable) (Expression, error) {
	// Rounding is discontinuous on integers,
	// but continuous otherwise, and 1: we must
	// also apply chain rule.
	if derivative, err := frac.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(DefectiveOnInt{frac.arg, big.NewInt(1)}, derivative), nil
	}
}


func (frac Frac) Simplify() (Expression, error) {
	if simplified, err := frac.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := (Frac{num}).Evaluate(Arguments{}); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return Frac{simplified}, nil
	}
}


func (frac Frac) String() string {
	return fmt.Sprintf("Frac(%s)", frac.arg)
}


func (frac Frac) IsSelfContained() bool {
	return true
}


// This expression evaluates to a constant if the underlying (bypassed) expression
// is not an integer, and returns an error otherwise.
type DefectiveOnInt struct {
	bypassed Expression
	result sets.Number
}


func (defectiveOnInt DefectiveOnInt) Curry(arguments Arguments) (Expression, error) {
	if curried, err := defectiveOnInt.bypassed.Curry(arguments); err != nil {
		return nil, err
	} else {
		return DefectiveOnInt{curried, defectiveOnInt.result}.Simplify()
	}
}


func (defectiveOnInt DefectiveOnInt) CollectVariables(variables Variables) {
	defectiveOnInt.bypassed.CollectVariables(variables)
}


func (defectiveOnInt DefectiveOnInt) IsConstant(wrt Variable) bool {
	return defectiveOnInt.bypassed.IsConstant(wrt)
}


func (defectiveOnInt DefectiveOnInt) Evaluate(args Arguments) (sets.Number, error) {
	evaluated, _ := defectiveOnInt.bypassed.Evaluate(args)
	switch vn := evaluated.(type) {
	case *big.Int:
		return nil, errors.UndefinedOnInteger
	case *big.Rat:
		if vn.IsInt() {
			return nil, errors.UndefinedOnInteger
		} else {
			return defectiveOnInt.result, nil
		}
	case *big.Float:
		if vn.IsInt() {
			return nil, errors.UndefinedOnInteger
		} else {
			return defectiveOnInt.result, nil
		}
	default:
		return defectiveOnInt.result, nil
	}
}


func (defectiveOnInt DefectiveOnInt) Derivative(wrt Variable) (Expression, error) {
	// Since defective-on-int returns always a constant,
	// or raises an error, depending on the underlying expression,
	// its derivative will always be 0 (or undefined on int).
	return DefectiveOnInt{
		defectiveOnInt.bypassed,
		big.NewInt(0),
	}, nil
}


func (defectiveOnInt DefectiveOnInt) Simplify() (Expression, error) {
	if simplified, err := defectiveOnInt.bypassed.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := (DefectiveOnInt{num, defectiveOnInt.result}).Evaluate(Arguments{}); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return DefectiveOnInt{simplified, defectiveOnInt.result}, nil
	}
}


func (defectiveOnInt DefectiveOnInt) String() string {
	return fmt.Sprintf("DOI(%s, %s)", defectiveOnInt.bypassed, defectiveOnInt.result)
}


func (defectiveOnInt DefectiveOnInt) IsSelfContained() bool {
	return true
}
