package diff

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"github.com/universe-10th-calculus/errors"
)


// TODO the power (^) operator goes here


type LnExpr struct {
	FunctionExpr
	arg   Expression
}


func (ln LnExpr) wrappedLn(power sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.LogarithmOfNegative
		}
	}()
	result = ops.Ln(power)
	return
}


func (ln LnExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := ln.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ln.wrappedLn(result)
	}
}


func (ln LnExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
}


func (ln LnExpr) Arguments() []Expression {
	return []Expression{ ln.arg }
}


type LogExpr struct {
	FunctionExpr
	power   Expression
	base    Expression
}


func (log LogExpr) wrappedLn(power, base sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.LogarithmOfNegative
		}
	}()
	result = ops.Log(power, base)
	return
}


func (log LogExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := log.power.Evaluate(args); err != nil {
		return nil, err
	} else if result2, err2 := log.base.Evaluate(args); err2 != nil {
		return nil, err
	} else {
		return log.wrappedLn(result, result2)
	}
}


func (log LogExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
}


func (log LogExpr) Arguments() []Expression {
	return []Expression{ log.base, log.power }
}


type ExpExpr struct {
	FunctionExpr
	exponent Expression
}


func (exp ExpExpr) Derivative(wrt Variable) (Expression, error) {
	return nil, nil
}


func (exp ExpExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := exp.exponent.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Exp(result), nil
	}
}


func (exp ExpExpr) Arguments() []Expression {
	return []Expression{ exp.exponent }
}


func Ln(power Expression) LnExpr {
	return LnExpr{FunctionExpr{"ln"}, power}
}


func Log(base, power Expression) LogExpr {
	return LogExpr{FunctionExpr{"log"}, power, base}
}


func Exp(exponent Expression) ExpExpr {
	return ExpExpr{FunctionExpr{"exp"}, exponent}
}