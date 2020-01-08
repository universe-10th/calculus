package expressions

import (
	"math/big"
	"github.com/universe-10th/calculus/errors"
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/ops"
	"fmt"
)


// FactorialExpr is a factorial expression. It only makes sense with numbers in N0 set.
// Factorial takes a single inner expression which MUST evaluate into N0.
type FactorialExpr struct {
	arg Expression
}


// Derivative returns a 0 constant expression, or an error, when attempting the derivative.
// This means: a 0 constant expression will be returned if the inner expression is constant
// with respect to the variable, or an error if it is not: factorial is not a derivable
// expression since it is not defined on R.
func (factorial FactorialExpr) Derivative(wrt Variable) (Expression, error) {
	// This expression is not derivable!!! must not include the
	// variable in the wrt parameter.
	variables := Variables{}
	factorial.arg.CollectVariables(variables)
	if !factorial.arg.IsConstant(wrt) {
		return nil, errors.ErrNotDerivableExpression
	} else {
		return Num(big.NewFloat(0)), nil
	}
}


// CollectVariables digs into the factorial's inner terms.
func (factorial FactorialExpr) CollectVariables(variables Variables) {
	factorial.arg.CollectVariables(variables)
}


// IsConstant will return whether the inner expression is constant with respect to the given variable.
func (factorial FactorialExpr) IsConstant(wrt Variable) bool {
	return factorial.arg.IsConstant(wrt)
}


// Simplify simplifies the inner expression and, if the simplified inner expression is constant,
// generates a new constant expression with the factorial of the returned constant's value.
// It will be an error if the inner constant does not evaluate into N0.
func (factorial FactorialExpr) Simplify() (Expression, error) {
	if simplified, err := factorial.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := factorial.wrappedFactorial(num.number); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return Factorial(simplified), nil
	}
}


func (factorial FactorialExpr) wrappedFactorial(input sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.ErrInvalidFactorialArgument
		}
	}()
	return ops.Factorial(input), nil
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (factorial FactorialExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := factorial.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Factorial(curried).Simplify()
	}
}


// Evaluate computes the factorial over the evaluated inner argument's value.
// It will be an error if the inner value does not evaluate into N0.
func (factorial FactorialExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := factorial.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return factorial.wrappedFactorial(result)
	}
}


// String represents the factorial as x! or (x)! appropriately.
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


// Factorial constructs a new factorial node with the given inner term.
func Factorial(value Expression) Expression {
	return FactorialExpr{value}
}