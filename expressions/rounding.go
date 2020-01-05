package expressions

import (
	"fmt"
	"math/big"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/errors"
)

type Round struct {
	arg Expression
}


type Frac struct {
	arg Expression
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
