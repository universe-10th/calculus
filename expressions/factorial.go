package diff

import (
	"math/big"
	"github.com/universe-10th-calculus/errors"
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"fmt"
)

type FactorialExpr struct {
	arg Expression
}


func (factorial FactorialExpr) Derivative(wrt Variable) (Expression, error) {
	// This expression is not derivable!!! must not include the
	// variable in the wrt parameter.
	variables := Variables{}
	factorial.arg.CollectVariables(variables)
	if _, ok := variables[wrt]; ok {
		return nil, errors.NotDerivableExpression
	} else {
		return Num(big.NewFloat(0)), nil
	}
}


func (factorial FactorialExpr) CollectVariables(variables Variables) {
	factorial.arg.CollectVariables(variables)
}


func (factorial FactorialExpr) wrappedFactorial(input sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.InvalidFactorialArgument
		}
	}()
	return ops.Factorial(input), nil
}


func (factorial FactorialExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := factorial.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return factorial.wrappedFactorial(result)
	}
}


func (factorial FactorialExpr) String() string {
	if _, ok := factorial.arg.(SelfContained); ok {
		return fmt.Sprintf("%s!", factorial.arg)
	} else {
		return fmt.Sprintf("(%s)!", factorial.arg)
	}
}


func (factorial FactorialExpr) IsSelfContained() bool {
	return true
}