package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
)

type TrigFunctionExpr struct {
	FunctionExpr
	arg   Expression
}


func (trigFunctionExpr TrigFunctionExpr) String() string {
	return FunctionDisplay(trigFunctionExpr)
}


func (trig TrigFunctionExpr) Arguments() []Expression {
	return []Expression{ trig.arg }
}


func (trig TrigFunctionExpr) CollectVariables(variables Variables) {
	trig.arg.CollectVariables(variables)
}


func (trig TrigFunctionExpr) IsConstant(wrt Variable) bool {
	return trig.arg.IsConstant(wrt)
}


type SinExpr struct {
	TrigFunctionExpr
}


type CosExpr struct {
	TrigFunctionExpr
}


type TanExpr struct {
	TrigFunctionExpr
}


func (sin SinExpr) Simplify() (Expression, error) {
	if simplified, err := sin.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Sin(num.number)}, nil
	} else {
		return Sin(simplified), nil
	}
}


func (sin SinExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := sin.arg.Evaluate(args); err == nil {
		return ops.Sin(result), nil
	} else {
		return nil, err
	}
}


func (sin SinExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := sin.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Cos(sin.arg), derivative).Simplify()
	}
}


func (cos CosExpr) Simplify() (Expression, error) {
	if simplified, err := cos.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Cos(num.number)}, nil
	} else {
		return Cos(simplified), nil
	}
}


func (cos CosExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := cos.arg.Evaluate(args); err == nil {
		return ops.Cos(result), nil
	} else {
		return nil, err
	}
}


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


func (tan TanExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := tan.arg.Evaluate(args); err == nil {
		return tan.wrappedTan(result)
	} else {
		return nil, err
	}
}


func (tan TanExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := tan.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Mul(Pow(Sin(tan.arg), Num(-2)), derivative).Simplify()
	}
}


func Sin(arg Expression) Expression {
	return SinExpr{TrigFunctionExpr{FunctionExpr{"sin"},arg}}
}


func Cos(arg Expression) Expression {
	return CosExpr{TrigFunctionExpr{FunctionExpr{"cos"},arg}}
}


func Tan(arg Expression) Expression {
	return TanExpr{TrigFunctionExpr{FunctionExpr{"tan"},arg}}
}
