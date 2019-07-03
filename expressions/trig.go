package expressions

import (
	"github.com/universe-10th/calculus/sets"
	"github.com/universe-10th/calculus/ops"
	"github.com/universe-10th/calculus/errors"
)


// TrigFunctionExpr is a structure providing some helpers over simple FunctionExpr structs.
// As an example, only one argument is used for these functions.
type TrigFunctionExpr struct {
	FunctionExpr
	arg   Expression
}


// String returns the standard string representation of these functions.
// They have one argument and their names are sin, cos, tan.
func (trigFunctionExpr TrigFunctionExpr) String() string {
	return FunctionDisplay(trigFunctionExpr)
}


// Arguments builds a list consisting on the single argument.
func (trig TrigFunctionExpr) Arguments() []Expression {
	return []Expression{ trig.arg }
}


// CollectVariables digs into the single argument expression.
func (trig TrigFunctionExpr) CollectVariables(variables Variables) {
	trig.arg.CollectVariables(variables)
}


// IsConstant evaluates whether the single argument expression is constant with respect to the given variable.
func (trig TrigFunctionExpr) IsConstant(wrt Variable) bool {
	return trig.arg.IsConstant(wrt)
}


// SinExpr stands for a sine expression.
type SinExpr struct {
	TrigFunctionExpr
}


// CosExpr stands for a cosine expression.
type CosExpr struct {
	TrigFunctionExpr
}


// TanExpr stands for a tangent expression.
type TanExpr struct {
	TrigFunctionExpr
}


// Simplify attempts reducing a sine expression to a constant.
// It first simplifies the argument and, if it turns to be constant, returns a constant expression with its sine.
func (sin SinExpr) Simplify() (Expression, error) {
	if simplified, err := sin.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Sin(num.number)}, nil
	} else {
		return Sin(simplified), nil
	}
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (sin SinExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := sin.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Sin(curried).Simplify()
	}
}


// Evaluate computes the sine expression by first computing the inner expression, and then applying the sine.
func (sin SinExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := sin.arg.Evaluate(args); err == nil {
		return ops.Sin(result), nil
	} else {
		return nil, err
	}
}


// Derivative uses the sine rule and also applies the chain rule.
func (sin SinExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := sin.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Cos(sin.arg), derivative).Simplify()
	}
}


// Simplify attempts reducing a cosine expression to a constant.
// It first simplifies the argument and, if it turns to be constant, returns a constant expression with its sine.
func (cos CosExpr) Simplify() (Expression, error) {
	if simplified, err := cos.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Cos(num.number)}, nil
	} else {
		return Cos(simplified), nil
	}
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (cos CosExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := cos.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Sin(curried).Simplify()
	}
}


// Evaluate computes the cosine over the evaluated value of the inner argument.
func (cos CosExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := cos.arg.Evaluate(args); err == nil {
		return ops.Cos(result), nil
	} else {
		return nil, err
	}
}


// Derivative applies the cosine rule and also the chain rule.
func (cos CosExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := cos.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Negated(Sin(cos.arg)), derivative).Simplify()
	}
}


func (tan TanExpr) wrappedTan(input sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.TangentOfVertical
		}
	}()
	result = ops.Tan(input)
	return
}


// Simplify attempts reducing a tangent expression to a constant.
// It first simplifies the argument and, if it turns to be constant, returns a constant expression with its sine.
func (tan TanExpr) Simplify() (Expression, error) {
	if simplified, err := tan.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		if result, err := tan.wrappedTan(num.number); err != nil {
			return nil, err
		} else {
			return Constant{result}, nil
		}
	} else {
		return Tan(simplified), nil
	}
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (tan TanExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := tan.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Sin(curried).Simplify()
	}
}


// Evaluate computes the tangent over the evaluated value of the inner expression.
func (tan TanExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := tan.arg.Evaluate(args); err == nil {
		return tan.wrappedTan(result)
	} else {
		return nil, err
	}
}


// Derivative uses the tangent rule (actually, the derivative of sin(x)/cos(x)) and also applies the chain rule.
func (tan TanExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := tan.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Pow(Sin(tan.arg), Num(-2)), derivative).Simplify()
	}
}


// Sin construct a sine expression.
func Sin(arg Expression) Expression {
	return SinExpr{TrigFunctionExpr{FunctionExpr{"sin"},arg}}
}


// Cos constructs a cosine expression.
func Cos(arg Expression) Expression {
	return CosExpr{TrigFunctionExpr{FunctionExpr{"cos"},arg}}
}


// Tan constructs a tangent expression.
func Tan(arg Expression) Expression {
	return TanExpr{TrigFunctionExpr{FunctionExpr{"tan"},arg}}
}
