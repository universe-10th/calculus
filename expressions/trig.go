package diff

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
)

type TrigFunctionExpr struct {
	FunctionExpr
	arg   Expression
}


func (trig TrigFunctionExpr) Arguments() []Expression {
	return []Expression{ trig.arg }
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


func (sin SinExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := sin.arg.Evaluate(args); err != nil {
		return ops.Sin(result), nil
	} else {
		return nil, err
	}
}


func (sin SinExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
}


func (cos CosExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := cos.arg.Evaluate(args); err != nil {
		return ops.Cos(result), nil
	} else {
		return nil, err
	}
}


func (cos CosExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
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


func (tan TanExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := tan.arg.Evaluate(args); err != nil {
		return tan.wrappedTan(result)
	} else {
		return nil, err
	}
}


func (tan TanExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
}


func Sin(arg Expression) SinExpr {
	return SinExpr{TrigFunctionExpr{FunctionExpr{"sin"},arg}}
}


func Cos(arg Expression) CosExpr {
	return CosExpr{TrigFunctionExpr{FunctionExpr{"cos"},arg}}
}


func Tan(arg Expression) TanExpr {
	return TanExpr{TrigFunctionExpr{FunctionExpr{"tan"},arg}}
}
