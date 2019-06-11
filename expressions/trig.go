package diff

import (
	"fmt"
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
)

type TrigFunction struct {
	arg   Expression
}


type SinExpr struct {
	TrigFunction
}


type CosExpr struct {
	TrigFunction
}


type TanExpr struct {
	TrigFunction
}


func (sin SinExpr) String() string {
	return fmt.Sprintf("sin(%s)", sin.arg)
}


func (sin SinExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := sin.arg.Evaluate(args); err != nil {
		return ops.Sin(result), nil
	} else {
		return nil, err
	}
}


func (sin SinExpr) Derivative(wrt Variable) Expression {
	return nil
}


func (cos CosExpr) String() string {
	return fmt.Sprintf("cos(%s)", cos.arg)
}


func (cos CosExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := cos.arg.Evaluate(args); err != nil {
		return ops.Cos(result), nil
	} else {
		return nil, err
	}
}


func (cos CosExpr) Derivative(wrt Variable) Expression {
	return nil
}


func (tan TanExpr) String() string {
	return fmt.Sprintf("tan(%s)", tan.arg)
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


func (tan TanExpr) Derivative(wrt Variable) Expression {
	return nil
}


func Sin(arg Expression) SinExpr {
	return SinExpr{TrigFunction{arg}}
}


func Cos(arg Expression) CosExpr {
	return CosExpr{TrigFunction{arg}}
}


func Tan(arg Expression) TanExpr {
	return TanExpr{TrigFunction{arg}}
}
